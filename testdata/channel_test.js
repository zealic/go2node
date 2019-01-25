const net = require('net');

const METHOD = process.argv.slice(2)[0];

function read() {
  process.send({black: "heart"});
}

function read_bigmsg() {
  var handle = net.connect(53, '1.2.4.8', ()=> {
    process.send({black: 'c'.repeat(1000000)}, handle);
  });
}

function write() {
  process.on("message", function (msg, handle) {
    process.send({value:"6553588"}, handle);
  });
}

function write_bigmsg() {
  process.on("message", function (msg, handle) {
    if(msg.value.length >= 1000000) {
      process.send({value:"wtf"}, handle);
      return;
    }
    process.exit(1);
  });
}

let methods = {
  read: read,
  read_bigmsg: read_bigmsg,
  write: write,
  write_bigmsg: write_bigmsg
};

methods[METHOD]();
