package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"

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
	var action = "UPSCALE"
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
	var res UpscaleRes
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return &types.UpscaleResp{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		}, nil
	}
	if len(res.Result) <= 0 {
		return nil, err
	}
	return &types.UpscaleResp{
		Code: http.StatusOK,
		Msg:  "Success",
	}, nil
}
