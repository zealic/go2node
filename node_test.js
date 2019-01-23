const METHOD = process.argv.slice(2)[0];

function reader() {
  process.send({
    black: "heart"
  });
}

function writer() {
  process.on("message", function (msg, handle) {
    process.send({value:"6553588"}, handle);
  });
}

let methods = {
  reader: reader,
  writer: writer
};

methods[METHOD]();
