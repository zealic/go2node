// +build channel_children

package go2node

import (
	"strings"
	"testing"

	"github.com/zealic/go2node/ipc"
)

func TestRunAsNodeChilren(t *testing.T) {
	channel, err := RunAsNodeChilren()
	if err != nil {
		panic(err)
	}

	// Parent
	channel.Writer <- &NodeMessage{Message: `{"name":"parent"}`}
	msg := <-channel.Reader
	if strings.Compare(string(msg.Message), `{"say":"We are one!"}`) != 0 {
		t.Fatal("Message not matched: ", string(msg.Message))
	}

	// ParentWithHandle
	sp, _ := ipc.Socketpair()
	channel.Writer <- &NodeMessage{
		Message: `{"name":"parentWithHandle"}`,
		Handle:  sp[0],
	}
	msg = <-channel.Reader
	if strings.Compare(string(msg.Message), `{"say":"For the Lich King!"}`) != 0 {
		t.Fatal("Message not matched: ", string(msg.Message))
	}
	if msg.Handle == nil {
		t.Fatal("Reply handle is required.")
	}
}
