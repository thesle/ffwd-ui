//go:build linux

package system

import (
	"bufio"
	"ffwd-ui/models"
	"os"
	"strings"
)

func GetAllMountPoints() ([]models.MountPoint, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var mounts []models.MountPoint
	seen := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		mountPoint := fields[1]

		if seen[mountPoint] {
			continue
		}

		if strings.HasPrefix(mountPoint, "/dev") ||
			strings.HasPrefix(mountPoint, "/proc") ||
			strings.HasPrefix(mountPoint, "/sys") ||
			strings.HasPrefix(mountPoint, "/run") ||
			mountPoint == "/boot" ||
			mountPoint == "/boot/efi" {
			continue
		}

		diskSpace, err := GetDiskSpace(mountPoint)
		if err != nil {
			continue
		}

		mounts = append(mounts, *diskSpace)
		seen[mountPoint] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mounts, nil
}
