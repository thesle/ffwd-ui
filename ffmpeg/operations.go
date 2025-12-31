package ffmpeg

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func BuildTrimStartCommand(input, output string, seconds float64) []string {
	return []string{
		"-ss", fmt.Sprintf("%.2f", seconds),
		"-i", input,
		"-c", "copy",
		output,
	}
}

func BuildTrimToLengthCommand(input, output string, duration float64) []string {
	return []string{
		"-i", input,
		"-t", fmt.Sprintf("%.2f", duration),
		"-c", "copy",
		output,
	}
}

func BuildExtractAudioCommand(input, output, format string) []string {
	format = strings.ToLower(format)

	var args []string
	args = append(args, "-i", input, "-vn")

	switch format {
	case "mp3":
		args = append(args, "-acodec", "libmp3lame", "-q:a", "2")
	case "aac":
		args = append(args, "-acodec", "aac", "-b:a", "192k")
	case "wav":
		args = append(args, "-acodec", "pcm_s16le")
	case "flac":
		args = append(args, "-acodec", "flac")
	default:
		args = append(args, "-acodec", "libmp3lame", "-q:a", "2")
	}

	args = append(args, output)
	return args
}

func BuildConvertFormatCommand(input, output string) []string {
	return []string{
		"-i", input,
		"-c", "copy",
		output,
	}
}

func BuildChangeResolutionCommand(input, output string, width, height int, hwAccel string) []string {
	var scale string
	if width > 0 && height > 0 {
		scale = fmt.Sprintf("%d:%d", width, height)
	} else if width > 0 {
		scale = fmt.Sprintf("%d:-1", width)
	} else if height > 0 {
		scale = fmt.Sprintf("-1:%d", height)
	} else {
		scale = "1280:720"
	}

	args := []string{"-i", input}

	if hwAccel != "" && hwAccel != "none" {
		args = append(args, "-c:v", hwAccel)
	}

	args = append(args, "-vf", fmt.Sprintf("scale=%s", scale), "-c:a", "copy", output)
	return args
}

func BuildAdjustVolumeCommand(input, output string, volumePercent int) []string {
	volumeMultiplier := float64(volumePercent) / 100.0
	return []string{
		"-i", input,
		"-af", fmt.Sprintf("volume=%.2f", volumeMultiplier),
		"-c:v", "copy",
		output,
	}
}

func BuildTrimRangeCommand(input, output string, startSeconds, endSeconds float64) []string {
	return []string{
		"-ss", fmt.Sprintf("%.2f", startSeconds),
		"-to", fmt.Sprintf("%.2f", endSeconds),
		"-i", input,
		"-c", "copy",
		output,
	}
}

func BuildCropVideoCommand(input, output string, width, height, x, y int) []string {
	return []string{
		"-i", input,
		"-vf", fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y),
		"-c:a", "copy",
		output,
	}
}

func BuildAdjustBitrateCommand(input, output string, videoBitrate, audioBitrate, hwAccel string, twoPass bool) []string {
	args := []string{"-i", input}

	if videoBitrate != "" {
		if hwAccel != "" && hwAccel != "none" {
			args = append(args, "-c:v", hwAccel)
		}
		args = append(args, "-b:v", videoBitrate)

		if twoPass {
			args = append(args, "-pass", "1", "-f", "null")
		}
	} else {
		args = append(args, "-c:v", "copy")
	}

	if audioBitrate != "" {
		args = append(args, "-b:a", audioBitrate)
	} else {
		args = append(args, "-c:a", "copy")
	}

	if !twoPass {
		args = append(args, output)
	}
	return args
}

func DetectHardwareEncoder() string {
	// Try to detect available hardware encoders
	encoders := []string{"h264_nvenc", "h264_qsv", "h264_videotoolbox", "h264_vaapi"}

	for _, encoder := range encoders {
		cmd := exec.Command("ffmpeg", "-hide_banner", "-encoders")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		if strings.Contains(string(output), encoder) {
			return encoder
		}
	}

	return "none"
}

func BuildCommandString(args []string) string {
	quotedArgs := make([]string, len(args))
	for i, arg := range args {
		quotedArgs[i] = quoteArg(arg)
	}
	return "ffmpeg " + strings.Join(quotedArgs, " ")
}

func quoteArg(arg string) string {
	// Check if argument needs quoting
	needsQuote := strings.ContainsAny(arg, " \t\n'\"()[]{}$&|;<>~`#*?")

	if !needsQuote {
		return arg
	}

	if runtime.GOOS == "windows" {
		// Windows: use double quotes and escape internal quotes
		escaped := strings.ReplaceAll(arg, "\"", "\\\"")
		return "\"" + escaped + "\""
	}

	// Linux/macOS: use single quotes (safest for shell)
	// Replace any single quotes with '\'' sequence
	escaped := strings.ReplaceAll(arg, "'", "'\\''")
	return "'" + escaped + "'"
}
