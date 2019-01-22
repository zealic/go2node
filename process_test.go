package go2node

import (
	"os"
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

	msg := <-channel.Reader
	const expectedContent = "{\"hello\":\"123\"}\n"
	if strings.Compare(string(msg.Data), expectedContent) != 0 {
		t.Fatal("Message not matched: ", string(msg.Data))
	}
}

func TestExec_Writer(t *testing.T) {
	cmd := exec.Command("node", "process_test.js", "writer")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	channel, err := Exec(cmd)
	if err != nil {
		t.Fatal(err)
	}

	msg := &Message{
		Data:  []byte(`{"cmd": "NODE_HANDLE", "type": "net.Socket", "msg": {"name": 10098}}` + "\n"),
		Files: []*os.File{os.Stdout},
	}
	channel.Writer <- msg

	// Verify
	ackMsg := <-channel.Reader

	if string(ackMsg.Data) != (`{"cmd":"NODE_HANDLE_ACK"}` + "\n") {
		t.Fatal("Ack msg not matched: ", string(ackMsg.Data))
	}

	msg = <-channel.Reader
	if string(msg.Data) != (`{"cmd":"NODE_HANDLE","type":"net.Socket","msg":"10087"}` + "\n") {
		t.Fatal("Message not matched: ", string(msg.Data))
	}
	cmd.Wait()
}
