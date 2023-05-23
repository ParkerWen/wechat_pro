package logic

import (
	"context"
	"errors"

	"github.com/ParkerWen/wechat_pro/app/callback/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/callback/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) Callback(req *types.CallbackReq) (*types.CallbackResp, error) {
	echoStr, cryptErr := l.svcCtx.WxBiz.VerifyURL(req.MsgSignature, req.TimeStamp, req.Nonce, req.EchoStr)
	if cryptErr != nil {
		return &types.CallbackResp{}, errors.New(cryptErr.ErrMsg)
	}
	return &types.CallbackResp{ErrCode: 0, ErrMsg: "ok", Data: echoStr}, nil
}
