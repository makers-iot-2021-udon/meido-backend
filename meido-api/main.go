package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const CLIENT_NUM = "CLIENT_NUM"
const ACCEPT_USER = "ACCEPT_USER"
const DENIED_USER = "DENIED_USER"

type ByteBroadCast struct {
	Message []byte
	Type    int
	Conn    *websocket.Conn
}

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan ByteBroadCast)
var StatusBroadCast = make(chan CurrentStatusMessage)
var MultiBroadCast = make(chan ByteBroadCast)

// var ch = make(chan bool)

func broadcastMessageToClients() {
	for {
		select {
		case message := <-StatusBroadCast:
			log.Println("Dispatched")
			// クライアントの数だけループ
			for client := range Clients {
				//　書き込む
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("error occurred while writing message to client: %v", err)
					client.Close()
					delete(Clients, client)
				}
			}

		case p := <-Broadcast:
			if err := p.Conn.WriteMessage(p.Type, p.Message); err != nil {
				log.Println(err)
				//	return
			}

		case p := <-MultiBroadCast:
			for client := range Clients {
				if err := client.WriteMessage(p.Type, p.Message); err != nil {
					log.Println(err)
					//return
				}
			}
		}
	}
}

func sendStatusRoutines() {
	for {
		StatusBroadCast <- currentStatus()
		time.Sleep(10 * time.Second)
	}
}

func reader(conn *websocket.Conn) {
	for {
		// Broadcast <- currentStatus()
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("recv:", string(p))
		flag := false
		p, flag = handler(p)
		log.Println("res:", string(p))

		//ここは同期的な処理だから特定のタイミングでAPI側からのメッセージを送信することも可能か（システムステータスなど）

		if flag {
			// //全体メッセージ
			// for client := range Clients {
			// 	if err := client.WriteMessage(messageType, p); err != nil {
			// 		log.Println(err)
			// 		//return
			// 	}
			// }
			MultiBroadCast <- ByteBroadCast{
				Type:    messageType,
				Message: p,
			}
		} else {
			// if err := conn.WriteMessage(messageType, p); err != nil {
			// 	log.Println(err)
			// 	//	return

			// }
			MultiBroadCast <- ByteBroadCast{
				Type:    messageType,
				Message: p,
				Conn:    conn,
			}
		}
		// //接続中の全ユーザーにパラメーターに現在のパラメーターを書き込む
		// statusMessage := currentStatus()

		// for client := range Clients {
		// 	if err := client.WriteMessage(messageType, statusMessage); err != nil {
		// 		log.Println(err)
		// 	}
		// }
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please connect via WebSocket")

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer ws.Close()

	Clients[ws] = true

	//クライアントカウント
	err = addValue(CLIENT_NUM)

	if err != nil {
		log.Println(err)
		ws.WriteMessage(1, []byte("Failed to addValue.close connection"))
		return
	}

	err = ws.WriteMessage(1, []byte("Hi Client!"))

	if err != nil {
		log.Println(err)
	}

	//メッセージの読み込みと書き込み
	reader(ws)

	log.Println("Client Disconnected")

	currentNum, err := declValue(CLIENT_NUM)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Successfully decrement CLIENT_NUM.\ncurrent num is :%d\n", currentNum)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
	// http.HandleFunc("/meido",meidoEndPoint)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Hello world")
	setupRoutes()
	go broadcastMessageToClients()
	go sendStatusRoutines()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
