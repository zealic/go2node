// +build channel_child

package go2node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zealic/go2node/ipc"
)

func TestRunAsNodeChild(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	channel, err := RunAsNodeChild()
	assert.NoError(err)

	// Parent
	channel.Write(&NodeMessage{Message: []byte(`{"name":"parent"}`)})
	msg, err := channel.Read()
	assert.NoError(err)
	require.Equal(`{"say":"We are one!"}`, string(msg.Message))

	// ParentWithHandle
	sp, err := ipc.Socketpair()
	assert.NoError(err)
	err = channel.Write(&NodeMessage{
		Message: []byte(`{"name":"parentWithHandle"}`),
		Handle:  sp[0],
	})
	assert.NoError(err)

	msg, err = channel.Read()
	assert.NoError(err)
	require.Equal(`{"say":"For the Lich King!"}`, string(msg.Message))
	assert.NotNil(msg.Handle)
}
