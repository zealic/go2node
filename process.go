package go2node

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const nodeChannelID = "NODE_CHANNEL_FD"

// IpcChannel ipc channel
type IpcChannel struct {
	Reader <-chan *Message
	Writer chan<- *Message
}

// Exec execute new nodejs child process with ipc channel
func Exec(cmd *exec.Cmd) (*IpcChannel, error) {
	fds, err := Socketpair()
	if err != nil {
		return nil, err
	}
	localSock := fds[0]
	remoteSocket := fds[1]

	cmd.ExtraFiles = append(cmd.ExtraFiles, remoteSocket)
	cmd.Env = []string{
		fmt.Sprintf("%s=%d", nodeChannelID, 2+len(cmd.ExtraFiles)),
	}

	// Handle message
	readChan := make(chan *Message, 1)
	go readIpcMessage(localSock, readChan)
	writeChan := make(chan *Message, 1)
	go writeIpcMessage(localSock, writeChan)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &IpcChannel{readChan, writeChan}, nil
}

func readIpcMessage(fd *os.File, msgChan chan *Message) {
	for {
		message, err := Recv(fd)
		if err != nil {
			log.Panic(err)
			break
		}

		msgChan <- message
	}
}

func writeIpcMessage(fd *os.File, msgChan chan *Message) {
	for {
		msg := <-msgChan
		err := Send(fd, msg)
		if err != nil {
			log.Panic(err)
			break
		}
	}
}
