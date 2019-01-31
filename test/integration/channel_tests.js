const cp = require('child_process');

const METHOD = process.argv.slice(2)[0];

function spawnTest() {
  let gtest = cp.spawn('go',
    ['test', '-tags', 'integration'], {
      stdio: [1, 2, 3, 'ipc'] // ipc is required
    });

  gtest.on('close', (code) => {
    process.exit(code);
  });
  return gtest;
}

function channel_child() {
  let child = spawnTest();

  child.on("message", function (msg, handle) {
    if (msg.name === 'parent') {
      child.send({
        say: "We are one!"
      });
    } else if (msg.name === 'parentWithHandle') {
      child.send({
        say: "For the Lich King!"
      }, handle);
    }
  });
}

let methods = {
  channel_child: channel_child
};

methods[METHOD]();