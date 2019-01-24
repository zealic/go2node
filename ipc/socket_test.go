package ipc

import (
	"testing"
)

func TestSocketpair(t *testing.T) {
	fds, err := Socketpair()

	if err != nil {
		t.Fatal(err)
	}

	if len(fds) != 2 {
		t.Fatal("Invalid count of fds")
	}

	if fds[0].Fd() <= 0 || fds[1].Fd() <= 0 {
		t.Fatal("Invalid fds")
	}
}
