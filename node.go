package go2node

import (
	"encoding/json"
	"os"
	"os/exec"
)

// NodeChannel node ipc channel
type NodeChannel struct {
	Reader <-chan *NodeMessage
	Writer chan<- *NodeMessage
}

// NodeMessage node ipc message
type NodeMessage struct {
	Message string
	Handle  *os.File
}

type rawNodeMessage struct {
	Cmd  string          `json:"cmd"`
	Type string          `json:"type"`
	Msg  json.RawMessage `json:"msg,omitempty"`
}

// ExecNode execute new nodejs child process with ipc channel
func ExecNode(cmd *exec.Cmd) (*NodeChannel, error) {
	ipcChannel, e := Exec(cmd)
	if e != nil {
		return nil, e
	}
	return newNodeChannel(ipcChannel)
}

func newNodeChannel(ipc *IpcChannel) (*NodeChannel, error) {
	// Handle message
	readChan := make(chan *NodeMessage, 1)
	go readNodeMessage(ipc, readChan)
	writeChan := make(chan *NodeMessage, 1)
	go writeNodeMessage(ipc, writeChan)

	return &NodeChannel{readChan, writeChan}, nil
}

func readNodeMessage(ipc *IpcChannel, msgChan chan *NodeMessage) {
	for msg := range ipc.Reader {
		rawMessage := new(rawNodeMessage)
		e := json.Unmarshal(msg.Data, rawMessage)
		if e != nil {
			panic(e)
		}

		switch rawMessage.Cmd {
		case "NODE_HANDLE":
			ipc.Writer <- &Message{
				Data: []byte(`{"cmd":"NODE_HANDLE_ACK"}` + "\n"),
			}
			msgChan <- &NodeMessage{
				Message: string(rawMessage.Msg),
				Handle:  msg.Files[0],
			}
		case "NODE_HANDLE_ACK":
		default:
			msgChan <- &NodeMessage{
				Message: string(msg.Data),
			}
		}
	}
}

func writeNodeMessage(ipc *IpcChannel, msgChan chan *NodeMessage) {
	for {
		msg := <-msgChan
		var ipcMsg *Message
		if msg.Handle == nil {
			ipcMsg = &Message{
				Data:  []byte(msg.Message),
				Files: []*os.File{},
			}
		} else {
			rawMsg := &rawNodeMessage{
				Cmd:  "NODE_HANDLE",
				Type: "net.Socket",
				Msg:  json.RawMessage(msg.Message),
			}

			bin, e := json.Marshal(rawMsg)
			if e != nil {
				panic(e)
			}

			ipcMsg = &Message{
				Data:  bin,
				Files: []*os.File{msg.Handle},
			}
		}

		ipcMsg.Data = append(ipcMsg.Data, '\n')
		ipc.Writer <- ipcMsg
	}
}
