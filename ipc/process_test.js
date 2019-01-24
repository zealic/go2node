const METHOD = process.argv.slice(2)[0];

function reader() {
  process.send({
    hello: "123"
  });
}

let methods = {
  reader: reader
};

methods[METHOD]();
