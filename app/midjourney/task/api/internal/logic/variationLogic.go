package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type VariationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type VariationRes struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

func NewVariationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VariationLogic {
	return &VariationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VariationLogic) Variation(req *types.VariationReq) (*types.VariationResp, error) {
	var action = "VARIATION"
	// 判断 task_id 所属的任务是否 UPSCALE，如果是则拒绝 VARIATION 操作
	ptask, err := l.svcCtx.TaskModel.FindOneByTaskId(l.ctx, req.TaskId)
	if err != nil {
		return &types.VariationResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}
	if ptask.Action == "UPSCALE" {
		return &types.VariationResp{
			Code: http.StatusBadRequest,
			Msg:  "放大后的图画不能做变换操作",
		}, nil
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
	var res VariationRes
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return &types.VariationResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}

	// 接口创建任务失败
	if len(res.Result) <= 0 {
		return &types.VariationResp{
			Code: http.StatusBadRequest,
			Msg:  res.Description,
		}, nil
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
		return &types.VariationResp{
			Code: http.StatusBadRequest,
			Msg:  errors.Wrapf(err, "Task Database Exception task : %+v , err: %v", task, err).Error(),
		}, nil
	}
	return &types.VariationResp{
		Code: http.StatusOK,
		Msg:  "Success",
	}, nil
}
