package ipc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSocketpair(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	fds, err := Socketpair()
	assert.NoError(err)

	require.Equal(2, len(fds))
	assert.True(fds[0].Fd() > 0)
	assert.True(fds[1].Fd() > 0)
}
