//go:build windows

package system

import (
	"ffwd-ui/models"
	"syscall"
)

var (
	getLogicalDrives = kernel32.NewProc("GetLogicalDrives")
)

func GetAllMountPoints() ([]models.MountPoint, error) {
	ret, _, _ := getLogicalDrives.Call()
	if ret == 0 {
		return nil, syscall.GetLastError()
	}

	var mounts []models.MountPoint
	drives := uint32(ret)

	for i := 0; i < 26; i++ {
		if drives&(1<<uint(i)) != 0 {
			driveLetter := string(rune('A'+i)) + ":\\"

			diskSpace, err := GetDiskSpace(driveLetter)
			if err != nil {
				continue
			}

			mounts = append(mounts, *diskSpace)
		}
	}

	return mounts, nil
}
