package go2node

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/zealic/go2node/ipc"
)

// NodeMessage node ipc message
type NodeMessage struct {
	Message []byte
	Handle  *os.File
	nack    bool
}

// Unmarshal unmarshal json encoded message
func (m *NodeMessage) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(m.Message), v)
}

func normNodeMessage(msg *ipc.Message) *NodeMessage {
	var handle *os.File
	if len(msg.Files) > 0 {
		handle = msg.Files[0]
	}

	data := strings.TrimSuffix(string(msg.Data), "\n")
	return &NodeMessage{
		Message: []byte(data),
		Handle:  handle,
	}
}
