process.on('message', function (msg, handle) {
  console.log(msg);
  process.exit(0);
});
process.send({hello: 'golang'});
