const METHOD = process.argv.slice(2)[0];

function reader() {
  process.send({
    hello: "123"
  });
}

function writer() {
  process.on("message", function (msg, handle) {
    process.send(`10087`, handle);
  });
}

let methods = {
  reader: reader,
  writer: writer
};

methods[METHOD]();
