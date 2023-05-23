package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.ListReq) (resp *types.ListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
