package handler

import (
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/logic"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func sendButtonInteractionTemplateCardHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendButtonInteractionTemplateCardReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSendButtonInteractionTemplateCardLogic(r.Context(), svcCtx)
		resp, err := l.SendButtonInteractionTemplateCard(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
