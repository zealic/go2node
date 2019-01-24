package go2node

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/zealic/go2node/ipc"
)

// NodeChannel node ipc channel
type NodeChannel struct {
	Reader <-chan *NodeMessage
	Writer chan<- *NodeMessage
	queue  []*NodeMessage
}

// NodeMessage node ipc message
type NodeMessage struct {
	Message string
	Handle  *os.File
	nack    bool
}

type rawNodeMessage struct {
	Cmd  string          `json:"cmd"`
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg,omitempty"`
}

const nodeChannelFD = "NODE_CHANNEL_FD"

// ExecNode execute new nodejs child process with ipc channel
func ExecNode(cmd *exec.Cmd) (*NodeChannel, error) {
	ipcChannel, e := ipc.Exec(cmd, nodeChannelFD)

	if e != nil {
		return nil, e
	}
	return newNodeChannel(ipcChannel)
}

func newNodeChannel(ipc *ipc.IpcChannel) (*NodeChannel, error) {
	// Handle message
	readChan := make(chan *NodeMessage, 1)
	writeChan := make(chan *NodeMessage, 1)
	channel := &NodeChannel{readChan, writeChan, []*NodeMessage{}}
	go channel.read(ipc, readChan, writeChan)
	go channel.write(ipc, writeChan)

	return channel, nil
}

func (c *NodeChannel) read(ipcChannel *ipc.IpcChannel,
	readChan chan *NodeMessage,
	writeChan chan *NodeMessage) {
	for msg := range ipcChannel.Reader {
		rawMessage := new(rawNodeMessage)
		e := json.Unmarshal(msg.Data, rawMessage)
		if e != nil {
			//readChan <- normNodeMessage(msg)
			panic(e)
		}

		switch rawMessage.Cmd {
		case "NODE_HANDLE":
			ipcChannel.Writer <- &ipc.Message{
				Data: []byte(`{"cmd":"NODE_HANDLE_ACK"}` + "\n"),
			}
			readChan <- &NodeMessage{
				Message: string(rawMessage.Msg),
				Handle:  msg.Files[0],
			}
		case "NODE_HANDLE_NACK":
			queue := c.queue
			c.queue = []*NodeMessage{}
			for _, m := range queue {
				writeChan <- &NodeMessage{
					Message: m.Message,
					Handle:  m.Handle,
					nack:    true,
				}
			}
		case "NODE_HANDLE_ACK":
			c.queue = []*NodeMessage{}
		default:
			readChan <- normNodeMessage(msg)
		}
	}
}

func normNodeMessage(msg *ipc.Message) *NodeMessage {
	var handle *os.File
	if len(msg.Files) > 0 {
		handle = msg.Files[0]
	}
	data := strings.TrimRight(string(msg.Data), "\n")
	return &NodeMessage{
		Message: data,
		Handle:  handle,
	}
}

func (c *NodeChannel) write(ipcChannel *ipc.IpcChannel, msgChan chan *NodeMessage) {
	for {
		msg := <-msgChan
		var ipcMsg *ipc.Message
		if msg.Handle == nil { // Normal message
			ipcMsg = &ipc.Message{
				Data:  []byte(msg.Message),
				Files: []*os.File{},
			}
		} else {
			// Default use naked message
			// NACK message will beo naked too
			ipcMsg = &ipc.Message{
				Data:  []byte(msg.Message),
				Files: []*os.File{msg.Handle},
			}
			// Send raw message
			if !msg.nack {
				c.queue = append(c.queue, msg)
				rawMsg := &rawNodeMessage{
					Cmd:  "NODE_HANDLE",
					Type: "net.Native",
					Msg:  json.RawMessage(msg.Message),
				}

				data, e := json.Marshal(rawMsg)
				if e != nil {
					panic(e)
				}
				ipcMsg.Data = data
			}
		}

		ipcMsg.Data = append(ipcMsg.Data, '\n')
		ipcChannel.Writer <- ipcMsg
	}
}
