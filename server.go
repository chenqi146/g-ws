package main

import (
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type WebsocketClient struct {
	Id                string
	Conn              *websocket.Conn
	LastHeartbeatTime int64
	
	Groups []string
	UserId string
}

type WebsocketUser struct {
	Id      string
	Clients []string
}

type WebsocketGroup struct {
	Id      string
	Clients []string
}

var (
	Clients sync.Map
	Users   sync.Map
	Groups  sync.Map
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
	
	clientId := initClientId(conn)
	
	if clientId == "" {
		return
	}
	log.Printf("[ws]: 连接websocket成功, clientId: %s\n", clientId)
	
	go handlerClientHeartbeat(clientId)
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			handlerClientDisconnect(clientId)
		} else {
			handlerClientMessage(conn, clientId, messageType, message)
		}
		
	}
}

func handlerClientMessage(conn *websocket.Conn, clientId string, messageType int, message []byte) {
	
	// ping  回复pong
	
	value, ok := Clients.Load(clientId)
	if !ok {
		handlerClientDisconnect(clientId)
		return
	}
	
	client := value.(WebsocketClient)
	
	if messageType != websocket.TextMessage {
		return
	}
	
	if string(message) != "ping" {
		return
	}
	
	if err := conn.WriteMessage(messageType, []byte("pong")); err != nil {
		handlerClientDisconnect(clientId)
		return
	}
	
	Clients.Store(clientId, WebsocketClient{
		Id:                clientId,
		Conn:              conn,
		LastHeartbeatTime: carbon.Now().Timestamp(),
		Groups:            client.Groups,
		UserId:            client.UserId,
	})
}

func handlerClientDisconnect(clientId string) {
	value, loaded := Clients.LoadAndDelete(clientId)
	if !loaded {
		return
	}
	client := value.(*WebsocketClient)
	
	if client.UserId != "" {
		// todo 解绑用户
	}
	
	if len(client.Groups) > 0 {
		// todo 解绑组
	}
}

func handlerClientHeartbeat(clientId string) {

}

func initClientId(conn *websocket.Conn) string {
	clientId := uuid.New().String()
	
	client := WebsocketClient{
		Id:                clientId,
		Conn:              conn,
		LastHeartbeatTime: carbon.Now().Timestamp(),
	}
	Clients.Store(clientId, client)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(clientId)); err != nil {
		log.Println("[ws]: 发生clientId异常", err)
		handlerClientDisconnect(clientId)
		return ""
	}
	
	return clientId
}

func main() {
	http.HandleFunc("/ws", ws)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
