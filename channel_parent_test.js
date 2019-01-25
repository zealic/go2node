const cp = require('child_process');

const METHOD = process.argv.slice(2)[0];

function spawnTest() {
  let gtest = cp.spawn('go',
    ['test', '-tags', 'channel_children'], {
      stdio: [0, 1, 2, 'ipc']
    });

  gtest.on('close', (code) => {
    process.exit(code);
  });
  return gtest;
}

function channel_children() {
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
  channel_children: channel_children
};

methods[METHOD]();