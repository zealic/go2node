package ipc

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendRecv(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	fds, err := Socketpair()
	assert.NoError(err)

	err = Send(fds[0], &Message{Data: []byte("123")})
	assert.NoError(err)

	msg, err := Recv(fds[1])
	assert.NoError(err)
	require.Equal("123", string(msg.Data))
}

func TestSendRecvWithFiles(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	fds, err := Socketpair()
	assert.NoError(err)

	err = Send(fds[0], &Message{
		Data:  []byte("bilibili"),
		Files: []*os.File{os.Stdin, os.Stdout},
	})
	assert.NoError(err)

	msg, err := Recv(fds[1])
	assert.NoError(err)
	require.Equal("bilibili", string(msg.Data))
	require.True(len(msg.Files) == 2)
}
