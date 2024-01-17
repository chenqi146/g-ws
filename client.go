package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
)

var done chan interface{}
var interrupt chan os.Signal

func receive(conn *websocket.Conn) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
			return
		}
		log.Printf("接收到：%s", message)
	}
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)
	
	signal.Notify(interrupt, os.Interrupt)
	
	wsUrl := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	go receive(conn)
	
}
