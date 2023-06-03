package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/types"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/ParkerWen/wechat_pro/common/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImagineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type ImagineRes struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

func NewImagineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImagineLogic {
	return &ImagineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImagineLogic) Imagine(req *types.ImagineReq) (*types.ImagineResp, error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	var action = "IMAGINE"
	m := map[string]interface{}{
		"action": action,
		"prompt": req.Prompt,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	r, _ := http.NewRequest("POST", "http://38.95.233.164:8088/mj/trigger/submit", contentReader)
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(r)
	var res ImagineRes
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return &types.ImagineResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}

	// 接口创建任务失败
	if len(res.Result) <= 0 {
		return &types.ImagineResp{
			Code: http.StatusBadRequest,
			Msg:  res.Description,
		}, nil
	}

	task := new(model.Task)
	task.UserId = userId
	task.Action = action
	task.Prompt = req.Prompt
	task.Index = 0
	task.Description = res.Description
	task.TaskId = res.Result
	task.Status = "PENDING"
	task.State = "valid"
	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	insertResult, err := l.svcCtx.TaskModel.InsertByImagine(l.ctx, task)
	if err != nil {
		return &types.ImagineResp{
			Code: http.StatusBadRequest,
			Msg:  fmt.Sprintf("Task Database Exception task : %+v , err: %v", task, err),
		}, nil
	}
	data := make(map[string]any)
	id, err := insertResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	data["id"] = id
	return &types.ImagineResp{
		Code: http.StatusOK,
		Msg:  "Success",
		Data: data,
	}, nil
}
