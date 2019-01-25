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

	buff := make([]byte, 1024)
	n, _, err := Recv(fds[1], buff, 4)
	assert.NoError(err)
	require.Equal(3, n)
	require.Equal("123", string(buff[:3]))
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

	buff := make([]byte, 1024)
	n, files, err := Recv(fds[1], buff, 4)
	assert.NoError(err)
	require.Equal(8, n)
	require.Equal("bilibili", string(buff[:8]))
	require.True(len(files) == 2)
}
