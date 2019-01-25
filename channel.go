package go2node

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/zealic/go2node/ipc"
)

// NodeChannel node ipc channel
type NodeChannel interface {
	// read a node message
	Read() (*NodeMessage, error)
	// write a node message
	Write(*NodeMessage) error
}

type nodeChannel struct {
	reader     chan *NodeMessage
	writer     chan *NodeMessage
	ipcChannel ipc.Channel
	queue      []*NodeMessage
}

type internalNodeMessage struct {
	Cmd  string          `json:"cmd"`
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg,omitempty"`
}

const nodeChannelFD = "NODE_CHANNEL_FD"
const nodeChannelDelim = '\n'

// ExecNode execute new nodejs child process with Node ipc channel
func ExecNode(cmd *exec.Cmd) (NodeChannel, error) {
	ipcChannel, e := ipc.Exec(cmd, nodeChannelFD)
	if e != nil {
		return nil, e
	}

	return newNodeChannel(ipcChannel)
}

func newNodeChannel(ipc ipc.Channel) (NodeChannel, error) {
	// Handle message
	readChan := make(chan *NodeMessage, 1)
	writeChan := make(chan *NodeMessage, 1)
	channel := &nodeChannel{
		reader:     readChan,
		writer:     writeChan,
		ipcChannel: ipc,
		queue:      []*NodeMessage{}}

	return channel, nil
}

func (c *nodeChannel) Read() (*NodeMessage, error) {
	for {
		ipcMsg, e := c.ipcChannel.ReadMessage(nodeChannelDelim)
		if e != nil {
			return nil, e
		}

		// Handle internal message
		intMsg := new(internalNodeMessage)
		e = json.Unmarshal(ipcMsg.Data, intMsg)
		if e != nil {
			return nil, e
		}

		msg, err := c.handleInternalMsg(ipcMsg, intMsg)

		// Confirm ACK and NACK, read next message
		if msg == nil && err == nil {
			continue
		}
		return msg, err
	}
}

func (c *nodeChannel) handleInternalMsg(
	ipcMsg *ipc.Message,
	intMsg *internalNodeMessage) (*NodeMessage, error) {
	var err error
	switch intMsg.Cmd {
	case "NODE_HANDLE":
		err = c.ipcChannel.WriteMessage(&ipc.Message{
			Data: []byte(`{"cmd":"NODE_HANDLE_ACK"}`),
		}, '\n')
		if err != nil {
			return nil, err
		}
		return &NodeMessage{
			Message: intMsg.Msg,
			Handle:  ipcMsg.Files[0],
		}, nil
	case "NODE_HANDLE_NACK":
		queue := c.queue
		c.queue = []*NodeMessage{}
		for _, m := range queue {
			err := c.Write(&NodeMessage{
				Message: m.Message,
				Handle:  m.Handle,
				nack:    true,
			})
			if err != nil {
				return nil, err
			}
		}
	case "NODE_HANDLE_ACK":
		c.queue = []*NodeMessage{}
	default:
		return normNodeMessage(ipcMsg), nil
	}
	// ACK and NACK
	return nil, nil
}

func (c *nodeChannel) Write(msg *NodeMessage) error {
	var ipcMsg *ipc.Message
	if msg.Handle == nil { // Normal message
		ipcMsg = &ipc.Message{
			Data:  []byte(msg.Message),
			Files: []*os.File{},
		}
	} else {
		// Default use naked message
		// NACK message will be naked too
		ipcMsg = &ipc.Message{
			Data:  []byte(msg.Message),
			Files: []*os.File{msg.Handle},
		}
		// Send raw message
		if !msg.nack {
			c.queue = append(c.queue, msg)
			intMsg := &internalNodeMessage{
				Cmd:  "NODE_HANDLE",
				Type: "net.Native",
				Msg:  json.RawMessage(msg.Message),
			}

			data, e := json.Marshal(intMsg)
			if e != nil {
				return e
			}
			ipcMsg.Data = data
		}
	}

	return c.ipcChannel.WriteMessage(ipcMsg, nodeChannelDelim)
}
