package ipc

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nodeChannelFD = "NODE_CHANNEL_FD"

const testFile = "../testdata/ipc/channel_test.js"

func execNodeFile(handler string) (*os.Process, Channel) {
	cmd := exec.Command("node", testFile, handler)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	channel, err := Exec(cmd, nodeChannelFD)
	if err != nil {
		panic(err)
	}

	return cmd.Process, channel
}

func TestExec_Read(t *testing.T) {
	_, require := assert.New(t), require.New(t)

	proc, channel := execNodeFile("read")
	defer func() {
		proc.Kill()
	}()

	msg, err := channel.ReadMessage('\n')
	require.NoError(err)
	require.Equal(`{"hello":"123"}`+"\n", string(msg.Data))
}
