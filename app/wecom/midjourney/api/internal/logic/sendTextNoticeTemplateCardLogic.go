package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendTextNoticeTemplateCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendTextNoticeTemplateCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendTextNoticeTemplateCardLogic {
	return &SendTextNoticeTemplateCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendTextNoticeTemplateCardLogic) SendTextNoticeTemplateCard(req *types.SendTextNoticeTemplateCardReq) (resp *types.SendTextNoticeTemplateCardResp, err error) {
	// todo: add your logic here and delete this line

	return
}
