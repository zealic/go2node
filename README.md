[![Build Status][travis-image]][travis-url]
[![GoDoc][godoc-image]][godoc-url]
[![GitHub release][release-image]][release-url]

# go2node

Package go2node provides a simple API to inter-process communicating for go and node.

go2node will **CREATE** or **RUN AS** Node child process with Node IPC protocol.


## Requirements

* Golang ≥ 1.11.x
* Node.js ≥ 8.x


## Quick start

### Golang to Node

`main.go`:

```golang
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
```

`child.js`:

```js
process.on('message', function (msg, handle) {
    console.log(msg);
    process.exit(0);
  });
process.send({hello: 'golang'});
```

Run `go run main.go` to test.

**Output**:

```
{"hello":"golang"}
{ hello: 'node' }
```

### Node to Golang

`main.js`:

```node
const child_process = require('child_process');

let child = child_process.spawn('go', ['run', 'gochild.go'], {
  stdio: [0, 1, 2, 'ipc']
});

child.on('close', (code) => {
  process.exit(code);
});

child.send({hello: "child"});
child.on('message', function(msg, handle) {
  if(msg.hello === "parent") {
    console.log(msg);
    process.exit(0);
  }
  process.exit(1);
});
```

`gochild.go`:

```golang
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
```

Run `node ./main.js` to test.

**Output**:

```
{"hello":"child"}
{ hello: 'parent' }
```


## Reference

* https://github.com/nodejs/node/blob/master/lib/child_process.js
* https://github.com/nodejs/node/blob/master/lib/internal/child_process.js
* https://github.com/nodejs/node/blob/master/src/stream_wrap.cc
* https://github.com/libuv/libuv/blob/master/src/unix/stream.c
* https://medium.com/js-imaginea/clustering-inter-process-communication-ipc-in-node-js-748f981214e9

[travis-image]:  https://travis-ci.org/zealic/go2node.svg
[travis-url]:    https://travis-ci.org/zealic/go2node
[godoc-image]:   https://godoc.org/github.com/zealic/go2node?status.svg
[godoc-url]:     https://godoc.org/github.com/zealic/go2node
[release-image]: https://img.shields.io/github/release/zealic/go2node.svg
[release-url]:   https://github.com/zealic/go2node/releases/latest
