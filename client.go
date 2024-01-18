package main

import (
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"time"
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
	
	for {
		select {
		case <-time.After(time.Second * time.Duration(1)):
			err := conn.WriteMessage(websocket.TextMessage, []byte("Hello World"))
			if err != nil {
				log.Fatal(err)
				return
			}
		case <-interrupt:
			log.Println("接收到中断信号")
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Fatal(err)
				return
			}
			select {
			case <-done:
				log.Println("接收到客户端关闭")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("超时")
			}
		}
	}
}
