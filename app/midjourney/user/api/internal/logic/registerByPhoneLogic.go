package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterByPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterByPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterByPhoneLogic {
	return &RegisterByPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterByPhoneLogic) RegisterByPhone(req *types.RegisterByPhoneReq) (resp *types.RegisterByPhoneResp, err error) {
	// todo: add your logic here and delete this line

	return
}
