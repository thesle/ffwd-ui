//go:build darwin

package system

import (
	"ffwd-ui/models"
	"os/exec"
	"strings"
)

func GetAllMountPoints() ([]models.MountPoint, error) {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return []models.MountPoint{}, nil
	}

	var mounts []models.MountPoint
	lines := strings.Split(string(output), "\n")

	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		mountPoint := fields[len(fields)-1]

		if strings.HasPrefix(mountPoint, "/dev") ||
			strings.HasPrefix(mountPoint, "/System") ||
			strings.HasPrefix(mountPoint, "/private") {
			continue
		}

		diskSpace, err := GetDiskSpace(mountPoint)
		if err != nil {
			continue
		}

		mounts = append(mounts, *diskSpace)
	}

	return mounts, nil
}
