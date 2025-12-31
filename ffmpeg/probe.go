package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"ffwd-ui/models"
)

type FFProbeOutput struct {
	Format  FFProbeFormat   `json:"format"`
	Streams []FFProbeStream `json:"streams"`
}

type FFProbeFormat struct {
	Filename   string `json:"filename"`
	Size       string `json:"size"`
	Duration   string `json:"duration"`
	FormatName string `json:"format_name"`
}

type FFProbeStream struct {
	CodecType string `json:"codec_type"`
	CodecName string `json:"codec_name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func ProbeFile(path string) (*models.FileInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %w", err)
	}

	var probe FFProbeOutput
	if err := json.Unmarshal(output, &probe); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	fileInfo := &models.FileInfo{
		Path:   path,
		Format: probe.Format.FormatName,
	}

	if size, err := strconv.ParseInt(probe.Format.Size, 10, 64); err == nil {
		fileInfo.Size = size
	}

	if duration, err := strconv.ParseFloat(probe.Format.Duration, 64); err == nil {
		fileInfo.Duration = duration
	}

	stat, err := os.Stat(path)
	if err == nil {
		fileInfo.Size = stat.Size()
	}

	for _, stream := range probe.Streams {
		if stream.CodecType == "video" {
			fileInfo.Codec = stream.CodecName
			fileInfo.Width = stream.Width
			fileInfo.Height = stream.Height
			break
		}
	}

	return fileInfo, nil
}
