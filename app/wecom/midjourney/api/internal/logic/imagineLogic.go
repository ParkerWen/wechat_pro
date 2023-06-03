package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"

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
	var action = "IMAGINE"
	m := map[string]interface{}{
		"action":     action,
		"prompt":     req.Prompt,
		"state":      req.State,
		"notifyHook": req.NotifyHook,
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
	if len(res.Result) <= 0 {
		return nil, errors.New("任务创建失败")
	}
	return &types.ImagineResp{
		Code: http.StatusOK,
		Msg:  "Success",
	}, nil
}
