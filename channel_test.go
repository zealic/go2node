package go2node

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/zealic/go2node/ipc"
)

const testFile = "channel_test.js"

func execNodeFile(handler string) (*os.Process, *NodeChannel) {
	cmd := exec.Command("node", testFile, handler)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	channel, err := ExecNode(cmd)
	if err != nil {
		panic(err)
	}

	return cmd.Process, channel
}

func TestExecNode_Reader(t *testing.T) {
	proc, channel := execNodeFile("reader")
	defer func() {
		proc.Kill()
	}()

	msg := <-channel.Reader
	const expectedContent = `{"black":"heart"}`
	if strings.Compare(string(msg.Message), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Message))
	}
}

func TestExecNode_Writer(t *testing.T) {
	proc, channel := execNodeFile("writer")
	defer func() {
		proc.Kill()
	}()

	sp, _ := ipc.Socketpair()
	msg := &NodeMessage{
		Message: `65535`,
		Handle:  sp[0],
	}
	channel.Writer <- msg

	msg = <-channel.Reader
	if string(msg.Message) != `{"value":"6553588"}` {
		t.Fatal("Message not matched: ", msg.Message)
	}
	if msg.Handle.Fd() == 0 {
		t.Fatal("Handle is empty")
	}
}
