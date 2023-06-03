package logic

import (
	"context"
	"encoding/xml"
	"errors"

	"github.com/ParkerWen/wechat_pro/app/wecom/callback/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/callback/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallbackPostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type Msg struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	MsgId        string `xml:"MsgId"`
	AgentID      string `xml:"AgentID"`
}

func NewCallbackPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackPostLogic {
	return &CallbackPostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackPostLogic) CallbackPost(req *types.CallbackReq, body []byte) (*types.CallbackResp, error) {
	echoStr, cryptErr := l.svcCtx.WxBiz.DecryptMsg(req.MsgSignature, req.TimeStamp, req.Nonce, body)
	if cryptErr != nil {
		return &types.CallbackResp{}, errors.New(cryptErr.ErrMsg)
	}
	v := Msg{}
	err := xml.Unmarshal(echoStr, &v)
	if err != nil {
		panic(err)
	}
	l.Info(string(v.Content))
	return &types.CallbackResp{ErrCode: 0, ErrMsg: "ok", Data: echoStr}, nil
}
