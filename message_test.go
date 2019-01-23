package go2node

import (
	"os"
	"testing"
)

func TestSendRecv(t *testing.T) {
	fds, err := Socketpair()
	if err != nil {
		t.Fatal(err)
	}

	err = Send(fds[0], &Message{Data: []byte("123")})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Recv(fds[1])
	if err != nil {
		t.Fatal(err)
	}
	if string(msg.Data) != "123" {
		t.Error("Message content not matched")
	}
}

func TestSendRecvWithFiles(t *testing.T) {
	fds, err := Socketpair()
	if err != nil {
		t.Fatal(err)
	}

	err = Send(fds[0], &Message{
		Data:  []byte("bilibili"),
		Files: []*os.File{os.Stdin, os.Stdout},
	})
	if err != nil {
		t.Fatal(err)
	}

	msg, err := Recv(fds[1])
	if err != nil {
		t.Fatal(err)
	}
	if string(msg.Data) != "bilibili" {
		t.Error("Message content not matched")
	}
	if len(msg.Files) != 2 {
		t.Error("Files count not matched")
	}
}
