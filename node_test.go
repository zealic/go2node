package go2node

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestExecNode_Reader(t *testing.T) {
	cmd := exec.Command("node", "node_test.js", "reader")
	channel, err := ExecNode(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
	}()

	msg := <-channel.Reader
	const expectedContent = `{"black":"heart"}` + "\n"
	if strings.Compare(string(msg.Message), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Message))
	}
}

func TestExecNode_Writer(t *testing.T) {
	cmd := exec.Command("node", "node_test.js", "writer")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	channel, err := ExecNode(cmd)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cmd.Process.Kill()
	}()

	msg := &NodeMessage{
		Message: `65535`,
		Handle:  os.Stdout,
	}
	channel.Writer <- msg

	msg = <-channel.Reader
	if string(msg.Message) != `"65535"` {
		t.Fatal("Message not matched: ", msg.Message)
	}
	if msg.Handle.Fd() == 0 {
		t.Fatal("Handle is empty")
	}
	cmd.Wait()
}
