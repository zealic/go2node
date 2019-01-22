package go2node

import (
	"os"
	"strconv"
	"syscall"
)

// Socketpair create a socketpair
func Socketpair() ([]*os.File, error) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	return []*os.File{
		os.NewFile(uintptr(fds[0]), "@/go2node/left"+strconv.Itoa(fds[0])),
		os.NewFile(uintptr(fds[1]), "@/go2node/right"+strconv.Itoa(fds[1])),
	}, err
}
