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
	m := map[string]interface{}{
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
		return nil, err
	}
	task := new(model.Task)
	task.Action = "IMAGINE"
	task.Prompt = req.Prompt
	task.Index = 0
	task.Description = res.Description
	task.TaskId = res.Result
	task.Status = "PENDING"
	task.State = "valid"
	task.CreatedAt = time.Now().Unix()
	task.UpdatedAt = time.Now().Unix()

	_, err = l.svcCtx.TaskModel.InsertByImagine(l.ctx, task)
	if err != nil {
		return nil, errors.Wrapf(err, "Task Database Exception task : %+v , err: %v", task, err)
	}
	return &types.ImagineResp{
		Code: http.StatusOK,
		Msg:  "Success",
	}, nil
}
