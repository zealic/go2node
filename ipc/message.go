package ipc

import (
	"os"
	"strconv"
	"syscall"
)

// Message message
type Message struct {
	Data  []byte
	Files []*os.File
}

// Recv receives data and file descriptors from a Unix domain socket.
//
// Num specifies the expected number of file descriptors in one message.
// Internal files' names to be assigned are specified via optional filenames
// argument.
//
// You need to close all files in the returned slice. The slice can be
// non-empty even if this function returns an error.
//
// Use net.FileConn() if you're receiving a network connection.
func Recv(sock *os.File) (*Message, error) {
	const buffSize int = 1024 * 64
	const maxFdCap = 64

	sockFd := int(sock.Fd())

	// recvmsg
	var mbuf []byte
	data := make([]byte, buffSize)
	mbuf = make([]byte, syscall.CmsgSpace(maxFdCap*4))
	n, oobn, _, _, err := syscall.Recvmsg(sockFd, data, mbuf, 0)
	if err != nil {
		return nil, err
	}
	if n < buffSize {
		data = data[:n]
	}

	files, err := parseCmsg(oobn, mbuf)
	if err != nil {
		return nil, err
	}

	return &Message{data, files}, nil
}

func parseCmsg(oobn int, mbuf []byte) ([]*os.File, error) {
	if oobn == 0 {
		return []*os.File{}, nil
	}

	// parse control msgs
	var msgs []syscall.SocketControlMessage
	msgs, err := syscall.ParseSocketControlMessage(mbuf[:oobn])
	if err != nil {
		return nil, err
	}

	// fds to files
	files := make([]*os.File, 0, len(msgs))
	for i := 0; i < len(msgs) && err == nil; i++ {
		var fds []int
		fds, err = syscall.ParseUnixRights(&msgs[i])

		for _, fd := range fds {
			files = append(files, os.NewFile(uintptr(fd), "@/go2node/fd/"+strconv.Itoa(fd)))
		}
	}

	return files, nil
}

// Send sends file descriptors to Unix domain socket.
//
// Please note that the number of descriptors in one message is limited
// and is rather small.
// Use conn.File() to get a file if you want to put a network connection.
func Send(sock *os.File, message *Message) error {
	sockFd := int(sock.Fd())
	return syscall.Sendmsg(sockFd, message.Data, makeRights(message), nil, 0)
}

func makeRights(msg *Message) []byte {
	if msg.Files == nil || len(msg.Files) == 0 {
		return nil
	}

	files := msg.Files
	fds := make([]int, len(files))
	for i := range files {
		fds[i] = int(files[i].Fd())
	}

	return syscall.UnixRights(fds...)
}
