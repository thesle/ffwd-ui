//go:build windows

package system

import (
	"ffwd-ui/models"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	getDiskFreeSpace = kernel32.NewProc("GetDiskFreeSpaceExW")
)

func GetDiskSpace(path string) (*models.MountPoint, error) {
	var freeBytesAvailable, totalBytes, totalFreeBytes int64

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}

	ret, _, _ := getDiskFreeSpace.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if ret == 0 {
		return nil, syscall.GetLastError()
	}

	return &models.MountPoint{
		Path:      path,
		Total:     uint64(totalBytes),
		Available: uint64(freeBytesAvailable),
		Used:      uint64(totalBytes - totalFreeBytes),
	}, nil
}
