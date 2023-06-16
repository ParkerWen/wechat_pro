package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/ParkerWen/wechat_pro/common/ctxdata"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpscaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type UpscaleRes struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

func NewUpscaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpscaleLogic {
	return &UpscaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpscaleLogic) Upscale(req *types.UpscaleReq) (*types.UpscaleResp, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	var action = "UPSCALE"
	// 判断 task_id 是否已经执行过 UPSCALE，如果执行过则不再执行 UPSCALE 操作
	whereBuilder := l.svcCtx.TaskModel.RowBuilder().Where("`index` = ?", req.Index)
	list, err := l.svcCtx.TaskModel.FindByParentTaskIdAndAction(l.ctx, req.TaskId, action, whereBuilder)
	if err != nil && err != model.ErrNotFound {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  errors.Wrapf(errors.New("数据库繁忙,请稍后再试"), "Failed to get pending task err : %v ", err).Error(),
		}, nil
	}
	if len(list) > 0 {
		url, err := url.Parse(list[0].ImageUrl)
		if err != nil {
			return nil, err
		}
		url.Scheme = "http"
		url.Host = "img.itcity.cc"
		data := make(map[string]any)
		data["id"] = list[0].Id
		data["task_id"] = list[0].TaskId
		data["parent_task_id"] = list[0].ParentTaskId
		data["image_url"] = url.String()
		return &types.UpscaleResp{
			Code: http.StatusOK,
			Msg:  "Success",
			Data: data,
		}, nil
	}
	m := map[string]interface{}{
		"action": action,
		"taskId": req.TaskId,
		"index":  req.Index,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	r, _ := http.NewRequest("POST", "http://38.95.233.164:8080/mj/submit/change", contentReader)
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(r)
	var res UpscaleRes
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}

	// 接口创建任务失败
	if len(res.Result) <= 0 {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  res.Description,
		}, nil
	}

	task := new(model.Task)
	task.UserId = userId
	task.Action = action
	task.Index = req.Index
	task.Description = res.Description
	task.TaskId = res.Result
	task.ParentTaskId = req.TaskId
	task.Status = "PENDING"
	task.State = "valid"
	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	insertResult, err := l.svcCtx.TaskModel.InsertByImagine(l.ctx, task)
	if err != nil {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  errors.Wrapf(err, "Task Database Exception task : %+v , err: %v", task, err).Error(),
		}, nil
	}
	data := make(map[string]any)
	id, err := insertResult.LastInsertId()
	if err != nil {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}
	data["id"] = id
	return &types.UpscaleResp{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	}, nil
}
