package ipc

import (
	"os/exec"
	"strings"
	"testing"
)

func TestExec_Reader(t *testing.T) {
	cmd := exec.Command("node", "process_test.js", "reader")
	channel, err := Exec(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
	}()

	msg := <-channel.Reader
	const expectedContent = "{\"hello\":\"123\"}\n"
	if strings.Compare(string(msg.Data), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Data))
	}
}
