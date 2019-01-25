package go2node

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zealic/go2node/ipc"
)

const testFile = "testdata/channel_test.js"

func execNodeFile(handler string) (*os.Process, NodeChannel) {
	cmd := exec.Command("node", testFile, handler)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	channel, err := ExecNode(cmd)
	if err != nil {
		panic(err)
	}

	return cmd.Process, channel
}

func TestExecNode_Read(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("read")
	defer func() {
		proc.Kill()
	}()

	msg, err := channel.Read()
	assert.NoError(err)
	require.Equal(`{"black":"heart"}`, string(msg.Message))
}

func TestExecNode_Read_BigMsg(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("read_bigmsg")
	defer func() {
		proc.Kill()
	}()

	msg, err := channel.Read()
	assert.NoError(err)
	require.True(len(msg.Message) > 1000000, string(msg.Message))
	require.NotNil(msg.Handle)
}

func TestExecNode_Write(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("write")
	defer func() {
		proc.Kill()
	}()

	sp, err := ipc.Socketpair()
	require.NoError(err)
	msg := &NodeMessage{
		Message: []byte(`65535`),
		Handle:  sp[0],
	}
	err = channel.Write(msg)
	require.NoError(err)

	msg, err = channel.Read()
	require.NoError(err)
	require.Equal(`{"value":"6553588"}`, string(msg.Message))
	assert.NotNil(msg.Handle)
}

func TestExecNode_Write_BigMsg(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	proc, channel := execNodeFile("write_bigmsg")
	defer func() {
		proc.Kill()
	}()

	sp, err := ipc.Socketpair()
	require.NoError(err)
	msg := &NodeMessage{
		Message: []byte(`{"value":"` + strings.Repeat("X", 1000000) + `"}`),
		Handle:  sp[0],
	}
	err = channel.Write(msg)
	require.NoError(err)

	msg, err = channel.Read()
	require.NoError(err)
	require.Equal(`{"value":"wtf"}`, string(msg.Message))
	assert.NotNil(msg.Handle)
}
