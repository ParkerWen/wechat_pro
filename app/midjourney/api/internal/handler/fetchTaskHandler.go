package handler

import (
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/logic"
	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func fetchTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FetchTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewFetchTaskLogic(r.Context(), svcCtx)
		resp, err := l.FetchTask(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
