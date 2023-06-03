package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/logic"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func notifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyReq
		rBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		json.Unmarshal(rBytes, &req)
		l := logic.NewNotifyLogic(r.Context(), svcCtx)
		resp, err := l.Notify(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
