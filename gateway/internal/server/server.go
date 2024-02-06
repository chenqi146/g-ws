package server

import (
	"context"
	"fmt"
	"g-ws/common"
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

func NewGatewayWsServer(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayWsServer {
	return &GatewayWsServer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		upgrader: &websocket.Upgrader{
			// ws握手过程中允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *GatewayWsServer) ServeWs(w http.ResponseWriter, r *http.Request) error {

	defer func() {
		if err := recover(); err != nil {
			s.Logger.Infof("ws server panic: %v\n", err)
		}
	}()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		_, _ = w.Write([]byte("websocket upgrader err: " + err.Error()))
		return err
	}

	// todo 鉴权?

	// 分配id 发送mq
	clientId, err := s.initClientId(conn)
	if err != nil {
		_, _ = w.Write([]byte("websocket init client id err: " + err.Error()))
		return err
	}
	// todo 连接成功发送通知

	defer s.handlerClientDisconnect(clientId)

	// 心跳
	go s.handlerClientHeartbeat(clientId)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logx.Errorf("read message err: %s", err.Error())
			return err
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
	s.Logger.Info("[ws]: 收到消息: " + string(message))

	client := value.(connection.WebsocketClient)

	if messageType != websocket.TextMessage {
		return
	}

	if string(message) == "ping" {

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

}

func (s *GatewayWsServer) handlerClientHeartbeat(clientId string) {
	heartbeatPeriod := s.svcCtx.Config.HeartbeatPeriod

	for {
		select {
		case <-time.After(5 * time.Second):
			value, ok := connection.Clients.Load(clientId)
			if !ok {
				break
			}
			client := value.(connection.WebsocketClient)

			fmt.Printf("%s: 当前: %d, 客户端: %d, diff: %d\n", clientId, carbon.Now().Timestamp(), client.LastHeartbeatTime, int64(heartbeatPeriod))
			if carbon.Now().Timestamp()-client.LastHeartbeatTime > int64(heartbeatPeriod) {
				s.Logger.Infof("[ws]: client: %s heartbeat timeout", clientId)
				s.handlerClientDisconnect(clientId)
				break
			}
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
		logx.Infof("[ws]: 发生clientId(%s)异常: %v", clientId, err)
		s.handlerClientDisconnect(clientId)
		return "", err
	}
	value, ok := connection.Clients.Load(clientId)
	if ok {
		s.Logger.Info("dwqd", value)
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
		//  解绑用户
		if userValue, ok := connection.Users.Load(client.UserId); ok {
			user := userValue.(*connection.WebsocketUser)
			user.Clients = common.DeleteSlice(user.Clients, clientId)
			connection.Users.Store(client.UserId, user)
		}
	}

	if len(client.Groups) > 0 {
		// 解绑组
		for _, groupId := range client.Groups {
			if groupValue, ok := connection.Groups.Load(groupId); ok {
				group := groupValue.(connection.WebsocketGroup)
				group.Clients = common.DeleteSlice(group.Clients, clientId)
				connection.Groups.Store(groupId, group)
			}
		}
	}

	// todo 发生通知  断开连接
}
