package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/callback/api/internal/logic"
	"github.com/ParkerWen/wechat_pro/app/wecom/callback/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/callback/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func callbackPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCallbackPostLogic(r.Context(), svcCtx)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpx.Error(w, err)
		}
		resp, err := l.CallbackPost(&req, body)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}