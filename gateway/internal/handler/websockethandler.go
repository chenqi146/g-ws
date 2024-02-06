package handler

import (
	"g-ws/gateway/internal/server"
	"g-ws/gateway/internal/svc"
	"g-ws/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		wsServer := server.NewGatewayWsServer(r.Context(), svcCtx)
		err := wsServer.ServeWs(w, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, "")
		}
	}
}
