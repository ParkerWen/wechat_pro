package logic

import (
	"context"
	"net/url"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchLogic {
	return &FetchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchLogic) Fetch(req *types.FetchReq) (*types.FetchResp, error) {
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
		url, err := url.Parse(task.ImageUrl)
		if err != nil {
			return nil, err
		}
		url.Scheme = "http"
		url.Host = "img.itcity.cc"
		task.ImageUrl = url.String()
		_ = copier.Copy(&tpTask, task)
	}

	return &types.FetchResp{
		TaskInfo: tpTask,
	}, nil
}
