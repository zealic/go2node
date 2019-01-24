package ipc

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const nodeChannelFD = "NODE_CHANNEL_FD"

const testFile = "channel_test.js"

func execNodeFile(handler string) (*os.Process, *Channel) {
	cmd := exec.Command("node", testFile, handler)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	channel, err := Exec(cmd, nodeChannelFD)
	if err != nil {
		panic(err)
	}

	return cmd.Process, channel
}

func TestExec_Reader(t *testing.T) {
	proc, channel := execNodeFile("reader")
	defer func() {
		proc.Kill()
	}()

	msg := <-channel.Reader
	const expectedContent = "{\"hello\":\"123\"}\n"
	if strings.Compare(string(msg.Data), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Data))
	}
}
