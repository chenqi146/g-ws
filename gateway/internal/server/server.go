package server

import (
	"context"
	"g-ws/gateway/internal/svc"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
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
	
	s.upgrader.Upgrade(w, r, nil)
	
	// 有人连接  发送mq
	
}
