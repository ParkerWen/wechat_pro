package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/model"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpscaleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpscaleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpscaleLogic {
	return &UpscaleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpscaleLogic) Upscale(req *types.UpscaleReq) (*types.UpscaleResp, error) {
	var action = "UPSCALE"
	// 判断 task_id 是否已经执行过 UPSCALE，如果执行过则不再执行 UPSCALE 操作
	whereBuilder := l.svcCtx.TaskModel.RowBuilder()
	list, err := l.svcCtx.TaskModel.FindByParentTaskIdAndAction(l.ctx, req.TaskId, action, whereBuilder)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(errors.New("数据库繁忙,请稍后再试"), "Failed to get pending task err : %v ", err)
	}
	if len(list) > 0 {
		return nil, errors.New("该任务已经执行过 UPSCALE 操作，不可重复执行")
	}
	m := map[string]interface{}{
		"action": action,
		"taskId": req.TaskId,
		"index":  req.Index,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	r, _ := http.NewRequest("POST", "http://38.95.233.164:8088/mj/trigger/submit", contentReader)
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(r)
	var res ImagineRes
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	task := new(model.Task)
	task.Action = action
	task.Index = req.Index
	task.Description = res.Description
	task.TaskId = res.Result
	task.ParentTaskId = req.TaskId
	task.Status = "PENDING"
	task.State = "valid"
	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	_, err = l.svcCtx.TaskModel.InsertByImagine(l.ctx, task)
	if err != nil {
		return nil, errors.Wrapf(err, "Task Database Exception task : %+v , err: %v", task, err)
	}
	return &types.UpscaleResp{
		Code: http.StatusOK,
		Msg:  "Success",
	}, nil
}
