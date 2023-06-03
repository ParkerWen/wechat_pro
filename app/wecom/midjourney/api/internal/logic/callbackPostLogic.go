package logic

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/ParkerWen/wechat_pro/common/tool"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	replyMsgXml = `<xml>
					<ToUserName><![CDATA[%s]]></ToUserName>
					<FromUserName><![CDATA[%s]]></FromUserName> 
					<CreateTime>%d</CreateTime>
					<MsgType><![CDATA[text]]></MsgType>
					<Content><![CDATA[%s]]></Content>
				</xml>`
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

type Reply struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   string `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
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
	l.Info(string(v.MsgType))
	if v.MsgType == "text" {
		var msg = "任务已提交，请耐心等待"
		imagine := NewImagineLogic(l.ctx, l.svcCtx)
		_, err = imagine.Imagine(&types.ImagineReq{
			Prompt:     v.Content,
			State:      v.FromUserName,
			NotifyHook: l.svcCtx.Config.Midjourney.NotifyHook,
		})
		if err != nil {
			msg = "创建任务失败，请从新提交任务"
		}
		t := time.Now().Unix()
		data, cryptErr := l.svcCtx.WxBiz.EncryptMsg(fmt.Sprintf(replyMsgXml, v.FromUserName, v.ToUserName, t, msg), fmt.Sprintf("%d", t), tool.Krand(8, tool.KC_RAND_KIND_NUM))
		if cryptErr != nil {
			return &types.CallbackResp{}, errors.New(cryptErr.ErrMsg)
		}
		return &types.CallbackResp{
			ErrCode: 0,
			ErrMsg:  "ok",
			Data:    data,
		}, nil
	}
	return &types.CallbackResp{ErrCode: 0, ErrMsg: "ok", Data: nil}, nil
}
