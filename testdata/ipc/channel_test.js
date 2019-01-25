const METHOD = process.argv.slice(2)[0];

function read() {
  process.send({
    hello: "123"
  });
}

let methods = {
  read: read
};

methods[METHOD]();
