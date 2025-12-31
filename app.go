package main

import (
	"context"
	"fmt"
	"path/filepath"

	"ffwd-ui/ffmpeg"
	"ffwd-ui/models"
	"ffwd-ui/system"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx      context.Context
	executor *ffmpeg.Executor
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.executor = ffmpeg.NewExecutor(ctx)

	a.executor.SetProgressCallback(func(percent float64, message string) {
		runtime.EventsEmit(ctx, "ffmpeg:progress", models.ProgressUpdate{
			Percent: percent,
			Message: message,
		})
	})

	a.executor.SetCompleteCallback(func() {
		runtime.EventsEmit(ctx, "ffmpeg:complete", true)
	})

	a.executor.SetErrorCallback(func(err error) {
		runtime.EventsEmit(ctx, "ffmpeg:error", err.Error())
	})
}

func (a *App) SelectInputFile() (string, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Input File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Media Files (Audio & Video)",
				Pattern:     "*.mp4;*.avi;*.mkv;*.mov;*.wmv;*.flv;*.webm;*.m4v;*.mp3;*.wav;*.aac;*.flac;*.m4a;*.ogg;*.wma;*.opus",
			},
			{
				DisplayName: "Video Files",
				Pattern:     "*.mp4;*.avi;*.mkv;*.mov;*.wmv;*.flv;*.webm;*.m4v",
			},
			{
				DisplayName: "Audio Files",
				Pattern:     "*.mp3;*.wav;*.aac;*.flac;*.m4a;*.ogg;*.wma;*.opus",
			},
			{
				DisplayName: "All Files",
				Pattern:     "*.*",
			},
		},
	})
	return file, err
}

func (a *App) SelectOutputFile(defaultName string) (string, error) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Select Output File",
		DefaultFilename: defaultName,
	})
	return file, err
}

func (a *App) GetFileInfo(path string) (*models.FileInfo, error) {
	return ffmpeg.ProbeFile(path)
}

func (a *App) ExtractThumbnail(inputPath string) (string, error) {
	return ffmpeg.ExtractThumbnail(inputPath)
}

func (a *App) TrimStart(input, output string, seconds float64) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildTrimStartCommand(input, output, seconds)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) TrimToLength(input, output string, duration float64) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildTrimToLengthCommand(input, output, duration)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) ExtractAudio(input, output, format string) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildExtractAudioCommand(input, output, format)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) CancelOperation() error {
	return a.executor.Cancel()
}

func (a *App) ConvertFormat(input, output string) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildConvertFormatCommand(input, output)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) ChangeResolution(input, output string, width, height int, hwAccel string) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildChangeResolutionCommand(input, output, width, height, hwAccel)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) AdjustVolume(input, output string, volumePercent int) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildAdjustVolumeCommand(input, output, volumePercent)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) TrimRange(input, output string, startSeconds, endSeconds float64) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildTrimRangeCommand(input, output, startSeconds, endSeconds)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) CropVideo(input, output string, width, height, x, y int) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildCropVideoCommand(input, output, width, height, x, y)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) AdjustBitrate(input, output, videoBitrate, audioBitrate, hwAccel string, twoPass bool) error {
	if a.executor.IsRunning() {
		return fmt.Errorf("operation already running")
	}

	args := ffmpeg.BuildAdjustBitrateCommand(input, output, videoBitrate, audioBitrate, hwAccel, twoPass)

	fileInfo, err := ffmpeg.ProbeFile(input)
	if err != nil {
		return err
	}

	return a.executor.Execute(args, fileInfo.Duration)
}

func (a *App) DetectHardwareEncoder() string {
	return ffmpeg.DetectHardwareEncoder()
}

func (a *App) GetDiskSpace() ([]models.MountPoint, error) {
	return system.GetAllMountPoints()
}

func (a *App) PreviewCommand(operation string, input, output string, params map[string]interface{}) (string, error) {
	var args []string

	switch operation {
	case "trim_start":
		if seconds, ok := params["seconds"].(float64); ok {
			args = ffmpeg.BuildTrimStartCommand(input, output, seconds)
		}
	case "trim_length":
		if duration, ok := params["duration"].(float64); ok {
			args = ffmpeg.BuildTrimToLengthCommand(input, output, duration)
		}
	case "extract_audio":
		if format, ok := params["format"].(string); ok {
			args = ffmpeg.BuildExtractAudioCommand(input, output, format)
		}
	case "convert_format":
		args = ffmpeg.BuildConvertFormatCommand(input, output)
	case "change_resolution":
		width := int(params["width"].(float64))
		height := int(params["height"].(float64))
		hwAccel := ""
		if hw, ok := params["hw_accel"].(string); ok {
			hwAccel = hw
		}
		args = ffmpeg.BuildChangeResolutionCommand(input, output, width, height, hwAccel)
	case "adjust_volume":
		volumePercent := int(params["volume_percent"].(float64))
		args = ffmpeg.BuildAdjustVolumeCommand(input, output, volumePercent)
	case "trim_range":
		startSeconds := params["start_seconds"].(float64)
		endSeconds := params["end_seconds"].(float64)
		args = ffmpeg.BuildTrimRangeCommand(input, output, startSeconds, endSeconds)
	case "crop_video":
		width := int(params["width"].(float64))
		height := int(params["height"].(float64))
		x := int(params["x"].(float64))
		y := int(params["y"].(float64))
		args = ffmpeg.BuildCropVideoCommand(input, output, width, height, x, y)
	case "adjust_bitrate":
		videoBitrate := params["video_bitrate"].(string)
		audioBitrate := params["audio_bitrate"].(string)
		hwAccel := ""
		twoPass := false
		if hw, ok := params["hw_accel"].(string); ok {
			hwAccel = hw
		}
		if tp, ok := params["two_pass"].(bool); ok {
			twoPass = tp
		}
		args = ffmpeg.BuildAdjustBitrateCommand(input, output, videoBitrate, audioBitrate, hwAccel, twoPass)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}

	return ffmpeg.BuildCommandString(args), nil
}

func (a *App) GetDefaultOutputName(inputPath, operation string) string {
	ext := filepath.Ext(inputPath)
	base := inputPath[:len(inputPath)-len(ext)]

	switch operation {
	case "trim_start":
		return base + "_trimmed" + ext
	case "trim_length":
		return base + "_cut" + ext
	case "trim_range":
		return base + "_trimmed" + ext
	case "extract_audio":
		return base + "_audio.mp3"
	case "convert_format":
		return base + "_converted" + ext
	case "change_resolution":
		return base + "_resized" + ext
	case "adjust_volume":
		return base + "_volume" + ext
	case "crop_video":
		return base + "_cropped" + ext
	case "adjust_bitrate":
		return base + "_bitrate" + ext
	default:
		return base + "_output" + ext
	}
}
