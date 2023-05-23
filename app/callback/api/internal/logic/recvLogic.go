package logic

import (
	"context"
	"encoding/xml"
	"errors"

	"github.com/ParkerWen/wechat_pro/app/callback/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/callback/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecvLogic struct {
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

func NewRecvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecvLogic {
	return &RecvLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// access token "nk0Bg3hPR0J4wdbxdk0_NgymVQyaLeH8WSal2lys-x4r3osZX6TUdwHNh_JUuHNLNyU_IyxicZsPuvTTCoKrTb-VUaENnbbQYIHukV9pyiiA_4Zl-8t9HBkX_mSJhuRKRc9qoVWtj6culAlMc1gju1eSP5L156v90jMrb3vwgIwqEvcZpaqOZ31Gtfaj0HM44FH0aN4h5dP7TW4IkF0QEw"
func (l *RecvLogic) Recv(req *types.CallbackReq, body []byte) (*types.CallbackResp, error) {
	echoStr, cryptErr := l.svcCtx.WxBiz.DecryptMsg(req.MsgSignature, req.TimeStamp, req.Nonce, body)
	if cryptErr != nil {
		return &types.CallbackResp{}, errors.New(cryptErr.ErrMsg)
	}
	// l.Info(string(echoStr))
	v := Msg{}
	err := xml.Unmarshal(echoStr, &v)
	if err != nil {
		panic(err)
	}
	l.Info(string(v.Content))
	return &types.CallbackResp{ErrCode: 0, ErrMsg: "ok", Data: echoStr}, nil
}
