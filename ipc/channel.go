package ipc

import (
	"fmt"
	"os"
	"os/exec"
)

// Channel ipc channel
type Channel interface {
	// read a ipc message
	Read() (*Message, error)
	// write a ipc message
	Write(*Message) error
}

type channel struct {
	fd *os.File
}

// Exec execute new child process with ipc channel, ipc fd will pass by fdEnvVarName
func Exec(cmd *exec.Cmd, fdEnvVarName string) (Channel, error) {
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
func FromFD(fd *os.File) Channel {
	return makeChannel(fd)
}

func makeChannel(fd *os.File) Channel {
	return &channel{fd}
}

func (c *channel) Read() (*Message, error) {
	return Recv(c.fd)
}

func (c *channel) Write(m *Message) error {
	return Send(c.fd, m)
}
