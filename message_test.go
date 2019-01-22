package go2node

import (
	"testing"
)

func TestSend(t *testing.T) {
	fds, err := Socketpair()
	if err != nil {
		t.Fatal(err)
	}

	err = Send(fds[0], &Message{Data: []byte("123")})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecv(t *testing.T) {
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
