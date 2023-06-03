package logic

import (
	"context"
	"net/http"
	"net/url"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"
	"github.com/ParkerWen/wechat_pro/common/xerr"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

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

func (l *ListLogic) List(req *types.ListReq) (*types.ListResp, error) {
	var resp []types.Task
	list, err := l.svcCtx.TaskModel.FindAll(l.ctx, l.svcCtx.TaskModel.RowBuilder(), "id DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "TaskList db err req : %+v , err : %v", req, err)
	}

	if len(list) > 0 {
		for _, task := range list {
			var typeTask types.Task
			url, err := url.Parse(task.ImageUrl)
			if err != nil {
				return nil, err
			}
			url.Scheme = "http"
			url.Host = "img.itcity.cc"
			task.ImageUrl = url.String()
			_ = copier.Copy(&typeTask, task)
			resp = append(resp, typeTask)
		}
	}

	return &types.ListResp{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: resp,
	}, nil
}
