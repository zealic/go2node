package main

import (
	"fmt"

	"github.com/zealic/go2node"
)

func main() {
	channel, err := go2node.RunAsNodeChild()
	if err != nil {
		panic(err)
	}

	// Golang will output: {"hello":"child"}
	msg, err := channel.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg.Message))

	// Node will output: {"hello":'parent'}
	err = channel.Write(&go2node.NodeMessage{
		Message: []byte(`{"hello":"parent"}`),
	})
	if err != nil {
		panic(err)
	}
}
