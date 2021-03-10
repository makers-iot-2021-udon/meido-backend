const WebSocket = require("ws");

const app = require('express')();
const server = require('http').Server(app);
const io = require('socket.io')(server);

const PORT = process.env.PORT || 5000;

app.get('/', (req, res) => {
    res.send('<h1>Hello world</h1>');
  });

  io.on("connection", (socket) => {
    console.log(socket.handshake.query); // prints { x: "42", EIO: "4", transport: "polling" }
    socket.on('message',function(msg){
        console.log(msg)
        io.emit('message',msg)
    })  
});

// io.listen(5001)
server.listen(PORT, function(){
    console.log('server listening. Port:' + PORT);
});

const ws = new WebSocket("wss://meido-app.cf/backend/ws", {
  perMessageDeflate: false,
});

ws.on("open", function open() {
  ws.send("something");
});

ws.on("message", function incoming(data) {
  // console.log(data);
});
