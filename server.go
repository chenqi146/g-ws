package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
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
	bytes, _ := io.ReadAll(r.Body)
	fmt.Println(string(bytes))
	fmt.Println(r.RequestURI)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
			break
		}
		log.Println(message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Fatal(err)
		}
		
	}
}

func main() {
	http.HandleFunc("/ws", ws)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
