package ipc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Channel ipc channel
type Channel struct {
	Reader <-chan *Message
	Writer chan<- *Message
}

// Exec execute new nodejs child process with ipc channel
func Exec(cmd *exec.Cmd, fdEnvVarName string) (*Channel, error) {
	fds, err := Socketpair()
	if err != nil {
		return nil, err
	}
	localSock := fds[0]
	remoteSocket := fds[1]

	cmd.ExtraFiles = append(cmd.ExtraFiles, remoteSocket)
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", fdEnvVarName, 2+len(cmd.ExtraFiles)))

	// Handle message
	channel := makeChannel(localSock)
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// FromFD setup channel from parent passed fd
func FromFD(fd *os.File) *Channel {
	return makeChannel(fd)
}

func makeChannel(fd *os.File) *Channel {
	// Handle message
	readChan := make(chan *Message, 1)
	go readIpcMessage(fd, readChan)
	writeChan := make(chan *Message, 1)
	go writeIpcMessage(fd, writeChan)
	return &Channel{readChan, writeChan}
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
