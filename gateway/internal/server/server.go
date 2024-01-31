package server

import (
	"context"
	"g-ws/gateway/internal/connection"
	"g-ws/gateway/internal/svc"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"
)

type GatewayWsServer struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	upgrader *websocket.Upgrader
}

func (s *GatewayWsServer) ServeWs(w http.ResponseWriter, r *http.Request) {
	
	defer func() {
		if err := recover(); err != nil {
			s.Logger.Infof("ws server panic: %v\n", err)
		}
	}()
	
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		_, _ = w.Write([]byte("websocket upgrader err: " + err.Error()))
		return
	}
	
	// todo 鉴权?
	
	// 分配id 发送mq
	clientId, err := s.initClientId(conn)
	if err != nil {
		_, _ = w.Write([]byte("websocket init client id err: " + err.Error()))
		return
	}
	defer s.handlerClientDisconnect(clientId)
	
	// 心跳
	go s.handlerClientHeartbeat(clientId)
	
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logx.Errorf("read message err: %s", err.Error())
			break
		} else {
			s.handlerClientMessage(conn, clientId, messageType, message)
		}
		
	}
}

func (s *GatewayWsServer) handlerClientMessage(conn *websocket.Conn, clientId string, messageType int, message []byte) {
	
	// ping  回复pong
	
	value, ok := connection.Clients.Load(clientId)
	if !ok {
		s.handlerClientDisconnect(clientId)
		return
	}
	
	client := value.(connection.WebsocketClient)
	
	if messageType != websocket.TextMessage {
		return
	}
	
	if string(message) != "ping" {
		return
	}
	
	if err := conn.WriteMessage(messageType, []byte("pong")); err != nil {
		s.handlerClientDisconnect(clientId)
		return
	}
	
	connection.Clients.Store(clientId, connection.WebsocketClient{
		Id:                clientId,
		Conn:              conn,
		LastHeartbeatTime: carbon.Now().Timestamp(),
		Groups:            client.Groups,
		UserId:            client.UserId,
	})
}

func (s *GatewayWsServer) handlerClientHeartbeat(clientId string) {
	// todo
	value, loaded := connection.Clients.LoadAndDelete(clientId)
	if !loaded {
		return
	}
	
	heartbeatPeriod := s.svcCtx.Config.HeartbeatPeriod
	defer func() {
		s.handlerClientDisconnect(clientId)
	}()
	
	client := value.(connection.WebsocketClient)
	
	for {
		select {
		case <-time.After(time.Duration(heartbeatPeriod) * time.Second):
			if err := client.Conn.WriteMessage(websocket.TextMessage, []byte("pong")); err != nil {
				s.Logger.Error("websocket heartbeat send err: " + err.Error())
				break
			}
			client.LastHeartbeatTime = carbon.Now().Timestamp()
			connection.Clients.Store(clientId, client)
		}
	}
}

func (s *GatewayWsServer) initClientId(conn *websocket.Conn) (string, error) {
	clientId := uuid.New().String()
	
	client := connection.WebsocketClient{
		Id:                clientId,
		Conn:              conn,
		LastHeartbeatTime: carbon.Now().Timestamp(),
	}
	connection.Clients.Store(clientId, client)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(clientId)); err != nil {
		logx.Info("[ws]: 发生clientId异常", err)
		s.handlerClientDisconnect(clientId)
		return "", err
	}
	
	return clientId, nil
}

func (s *GatewayWsServer) handlerClientDisconnect(clientId string) {
	value, loaded := connection.Clients.LoadAndDelete(clientId)
	if !loaded {
		return
	}
	client := value.(*connection.WebsocketClient)
	
	if client.UserId != "" {
		// todo 解绑用户
	}
	
	if len(client.Groups) > 0 {
		// todo 解绑组
	}
	
	// todo 发生通知  断开连接
}
