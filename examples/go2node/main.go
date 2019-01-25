package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zealic/go2node"
)

func main() {
	cmd := exec.Command("node", "child.js")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	channel, err := go2node.ExecNode(cmd)
	if err != nil {
		panic(err)
	}
	defer cmd.Process.Kill()

	// Node will output: {hello: "node"}
	channel.Write(&go2node.NodeMessage{
		Message: []byte(`{"hello": "node"}`),
	})

	// Golang will output: {"hello":"golang"}
	msg, err := channel.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg.Message))

	// Wait node child process exit
	cmd.Process.Wait()
}
