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

var CLIENT_NUM = "CLIENT_NUM"
var ACCEPT_USER = "ACCEPT_USER"
var DENIED_USER = "DENIED_USER"
var Clients = make(map[*websocket.Conn]bool)
var MultiBroadcast = make(chan []byte)
var Broadcast = make(chan []byte)

func broadcastMessageToClients() {
	for {
		// メッセージ受け取り
		message := <-Broadcast
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
	}
}
func reader(conn *websocket.Conn) {
	for {
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
			//全体メッセージ
			for client := range Clients {
				if err := client.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					//return
				}
			}
			//Broadcast <- p
		} else {
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				//	return
			}
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please connect via WebSocket")

}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	//go broadcastMessageToClients()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// defer ws.Close()

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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
