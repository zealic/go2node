package go2node

import (
	"os"
	"strings"

	"github.com/zealic/go2node/ipc"
)

// NodeMessage node ipc message
type NodeMessage struct {
	Message string
	Handle  *os.File
	nack    bool
}

func normNodeMessage(msg *ipc.Message) *NodeMessage {
	var handle *os.File
	if len(msg.Files) > 0 {
		handle = msg.Files[0]
	}

	data := strings.TrimSuffix(string(msg.Data), "\n")
	return &NodeMessage{
		Message: data,
		Handle:  handle,
	}
}
