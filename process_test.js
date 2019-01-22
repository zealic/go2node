const METHOD = process.argv.slice(2)[0];

function reader() {
  setTimeout(function () {
    process.send({
      hello: "123"
    });
    process.exit();
  }, 500);
}

function writer() {
  let exitFlag;
  process.on("message", function (msg, handle) {
    process.send(`10087`, handle);
    exitFlag = true;
  });

  let counter = 0;
  setInterval(() => {
    counter++;
    if (exitFlag || counter > 100) process.exit();
  }, 10);
}

let methods = {
  reader: reader,
  writer: writer
};

methods[METHOD]();
