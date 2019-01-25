// +build channel_children

package go2node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zealic/go2node/ipc"
)

func TestRunAsNodeChilren(t *testing.T) {
	assert, require := assert.New(t), require.New(t)

	channel, err := RunAsNodeChilren()
	assert.NoError(err)

	// Parent
	channel.Writer <- &NodeMessage{Message: []byte(`{"name":"parent"}`)}
	msg := <-channel.Reader
	require.Equal(`{"say":"We are one!"}`, string(msg.Message))

	// ParentWithHandle
	sp, err := ipc.Socketpair()
	assert.NoError(err)
	channel.Writer <- &NodeMessage{
		Message: []byte(`{"name":"parentWithHandle"}`),
		Handle:  sp[0],
	}
	msg = <-channel.Reader
	require.Equal(`{"say":"For the Lich King!"}`, string(msg.Message))
	assert.NotNil(msg.Handle)
}
