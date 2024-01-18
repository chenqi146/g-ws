package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	// ws握手过程中允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws(w http.ResponseWriter, r *http.Request) {
	
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ws server panic: %v\n", err)
		}
	}()
	
	// 升级ws
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	
	initClientId(conn)
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		log.Println(string(message))
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Fatal(err)
		}
		
	}
}

func initClientId(conn *websocket.Conn) string {
	return ""
}

func main() {
	http.HandleFunc("/ws", ws)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
