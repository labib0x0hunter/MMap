package MMap

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	PROT_READ  = syscall.PROT_READ
	PROT_WRITE = syscall.PROT_WRITE
	MAP_SHARED = syscall.MAP_SHARED
)

func Mmap(filename string, length int, prot int, flags int) (data []byte, err error) {
	fd, err := syscall.Open(filename, syscall.O_RDWR, 0)
	if err != nil {
		return
	}
	defer syscall.Close(fd)

	if length == -1 {
		var info syscall.Stat_t
		if err = syscall.Fstat(fd, &info); err != nil {
			return
		}
		length = int(info.Size)
	}

	if length == 0 {
		err = fmt.Errorf("cannot map zero size")
		return
	}

	data, err = syscall.Mmap(
		fd,
		0,
		length,
		prot,
		flags,
	)
	return
}

func Msync(data []byte) (err error) {
	_, _, errno := syscall.Syscall(
		syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(syscall.MS_SYNC),
	)

	if errno != 0 {
		return errno
	}
	return
}

func Munmap(data []byte) (err error) {
	return syscall.Munmap(data)
}