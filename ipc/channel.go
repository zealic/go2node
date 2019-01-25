package ipc

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// Channel ipc channel
type Channel interface {
	// read a ipc message
	ReadMessage(delim byte) (*Message, error)
	// write a ipc message
	WriteMessage(m *Message, delim byte) error
}

type channel struct {
	fd     *os.File
	files  []*os.File
	reader *bufio.Reader
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
	c := &channel{fd: fd}
	c.reader = bufio.NewReader(c)
	return c
}

func (c *channel) Read(buff []byte) (int, error) {
	n, files, err := Recv(c.fd, buff, 4)
	for _, f := range files {
		if f != nil {
			c.files = append(c.files, files...)
		}
	}
	return n, err
}

func (c *channel) ReadMessage(delim byte) (*Message, error) {
	data, err := c.reader.ReadBytes(delim)
	if err != nil {
		return nil, err
	}
	files := c.files
	c.files = make([]*os.File, 0)
	return &Message{
		Data:  data,
		Files: files,
	}, nil
}

func (c *channel) WriteMessage(m *Message, delim byte) error {
	m.Data = append(m.Data, delim)
	return Send(c.fd, m)
}
