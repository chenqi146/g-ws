package connection

import (
	"github.com/gorilla/websocket"
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
