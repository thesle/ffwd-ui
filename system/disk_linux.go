//go:build linux || darwin

package system

import (
	"ffwd-ui/models"
	"syscall"
)

func GetDiskSpace(path string) (*models.MountPoint, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - (stat.Bfree * uint64(stat.Bsize))

	return &models.MountPoint{
		Path:      path,
		Total:     total,
		Available: available,
		Used:      used,
	}, nil
}
