const child_process = require('child_process');

let child = child_process.spawn('go', ['run', 'gochild.go'], {
  stdio: [0, 1, 2, 'ipc']
});

child.on('close', (code) => {
  process.exit(code);
});

child.send({hello: "child"});
child.on('message', function(msg, handle) {
  if(msg.hello === "parent") {
    console.log(msg);
    process.exit(0);
  }
  process.exit(1);
});
