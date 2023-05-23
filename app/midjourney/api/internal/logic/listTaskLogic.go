package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTaskLogic {
	return &ListTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTaskLogic) ListTask(req *types.ListTaskReq) (resp *types.ListTaskResp, err error) {
	// todo: add your logic here and delete this line

	return
}
