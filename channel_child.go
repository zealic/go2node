package go2node

import (
	"errors"
	"os"
	"strconv"

	"github.com/zealic/go2node/ipc"
)

// RunAsNodeChild setup current process as node child process
func RunAsNodeChild() (NodeChannel, error) {
	fdStr := os.Getenv(nodeChannelFD)
	if len(fdStr) == 0 {
		return nil, errors.New(nodeChannelFD + " is required.")
	}
	channelFD, err := strconv.Atoi(fdStr)
	if err != nil {
		return nil, err
	}
	fd := os.NewFile(uintptr(channelFD), "/go2node/"+nodeChannelFD)
	ipcChannel := ipc.FromFD(fd)
	return newNodeChannel(ipcChannel)
}
