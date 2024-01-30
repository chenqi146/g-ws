package handler

import (
	"net/http"

	"g-ws/gateway/internal/logic"
	"g-ws/gateway/internal/svc"
	"g-ws/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GatewayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGatewayLogic(r.Context(), svcCtx)
		resp, err := l.Gateway(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
