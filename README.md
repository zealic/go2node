[![Build Status][travis-image]][travis-url]
[![GoDoc][godoc-image]][godoc-url]

# go2node

Package go2node provides a simple API to inter-process communicating for go and node.

go2node will *CREATE* or *RUN AS* Node child process with Node IPC channel protocol.


## Reference

* https://github.com/nodejs/node/blob/master/lib/child_process.js
* https://github.com/nodejs/node/blob/master/lib/internal/child_process.js
* https://github.com/nodejs/node/blob/master/src/stream_wrap.cc
* https://github.com/libuv/libuv/blob/master/src/unix/stream.c
* https://medium.com/js-imaginea/clustering-inter-process-communication-ipc-in-node-js-748f981214e9

[travis-image]: https://travis-ci.org/zealic/go2node.svg
[travis-url]:   https://travis-ci.org/zealic/go2node
[godoc-image]:  https://godoc.org/github.com/zealic/go2node?status.svg
[godoc-url]:    https://godoc.org/github.com/zealic/go2node
