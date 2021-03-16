const WebSocket = require("ws");

const app = require("express")();
const server = require("http").Server(app);
const io = require("socket.io")(server);

const PORT = process.env.PORT || 5000;

app.get("/", (req, res) => {
  res.send("<h1>Hello world</h1>");
});

io.on("connection", (socket) => {
  console.log(socket.handshake.query); // prints { x: "42", EIO: "4", transport: "polling" }
  socket.on("message", function (msg) {
    console.log(msg);
    ws.send(`{"action":"LOVE_MESSAGE","message":"${msg}"}`);
    // io.emit("message", msg);
  });
});

// io.listen(5001)
server.listen(PORT, function () {
  console.log("server listening. Port:" + PORT);
});

const ws = new WebSocket("ws://localhost:8080/ws", {
  perMessageDeflate: false,
});

ws.on("open", function open() {
  ws.send("something");
});
const acceptMessage = "SUCCESS"
const deniedMessage = "REJECTED"
ws.on("message", function incoming(data) {
  // console.log(data);
  try {
    const jsonObj = JSON.parse(data);
    switch (jsonObj.action) {
      case "LOVE_MESSAGE": {
        console.log(jsonObj);
        let msg = "";
        jsonObj.messages.forEach((str, index) => {
          if (index == 0) msg += str;
          else msg += "\n" + str;
        });
        console.log("msg:", msg);
        io.emit("message", msg);
        if (jsonObj.cert_message === acceptMessage) {
          io.emit("like", "like")
        } else if(jsonObj.cert_message === deniedMessage) {
          io.emit("dislike", "dislike")
        }
        // } else if (jsonObj.type === "SPECIAL") {
        //   io.emit("special", "special")
        // }
      }
      break;
    case "NOTIFY_CURRENT_STATUS": {
      console.log(jsonObj)
      let connectCount = jsonObj.connect_count
      let msg = "hoge"
      io.emit("conn", String(connectCount))

    }
    break;
    default:
      break;
    }
  } catch (e) {
    // Error handling
    console.log(e); // SyntaxError: Unexpected token o in JSON at position 1
    //console.log('ここには来ます');
  }
});