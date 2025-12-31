package ffmpeg

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ExtractThumbnail(inputPath string) (string, error) {
	// Create temp directory for thumbnail
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, fmt.Sprintf("ffwd_thumb_%d.jpg", os.Getpid()))

	// Extract frame at 1 second - use -ss after -i for better accuracy with some codecs
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale=320:-1",
		"-q:v", "2",
		"-y",
		outputPath,
	)

	// Capture stderr for better error messages
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	// Check if file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("thumbnail file not created at %s", outputPath)
	}

	// Read the file and convert to base64 data URL for reliable display in Wails
	imageData, err := os.ReadFile(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to read thumbnail: %w", err)
	}

	// Clean up temp file
	os.Remove(outputPath)

	// Return as data URL
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return "data:image/jpeg;base64," + base64Image, nil
}
