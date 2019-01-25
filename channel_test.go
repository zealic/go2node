package go2node

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	_, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("reader")
	defer func() {
		proc.Kill()
	}()

	msg := <-channel.Reader
	require.Equal(`{"black":"heart"}`, string(msg.Message))
}

func TestExecNode_Writer(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("writer")
	defer func() {
		proc.Kill()
	}()

	sp, err := ipc.Socketpair()
	assert.NoError(err)
	msg := &NodeMessage{
		Message: []byte(`65535`),
		Handle:  sp[0],
	}
	channel.Writer <- msg

	msg = <-channel.Reader
	require.Equal(`{"value":"6553588"}`, string(msg.Message))
	assert.NotNil(msg.Handle)
}
