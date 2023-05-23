package logic

import (
	"context"

	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/types"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchTaskLogic {
	return &FetchTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchTaskLogic) FetchTask(req *types.FetchTaskReq) (*types.FetchTaskResp, error) {
	var tpTask types.Task
	if req.ID > 0 {
		task, err := l.svcCtx.TaskModel.FindOne(l.ctx, req.ID)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&tpTask, task)
	}
	if len(req.TaskId) > 0 {
		task, err := l.svcCtx.TaskModel.FindOneByTaskId(l.ctx, req.TaskId)
		if err != nil {
			return nil, err
		}
		_ = copier.Copy(&tpTask, task)
	}

	return &types.FetchTaskResp{
		TaskInfo: tpTask,
	}, nil
}
