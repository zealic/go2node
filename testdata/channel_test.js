const METHOD = process.argv.slice(2)[0];

function read() {
  process.send({black: "heart"});
}

function write() {
  process.on("message", function (msg, handle) {
    process.send({value:"6553588"}, handle);
  });
}

let methods = {
  read: read,
  write: write
};

methods[METHOD]();
