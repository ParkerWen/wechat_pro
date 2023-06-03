package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/silenceper/wechat/util"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	sendButtonInteractionTemplateCardURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

type SendButtonInteractionTemplateCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendButtonInteractionTemplateCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendButtonInteractionTemplateCardLogic {
	return &SendButtonInteractionTemplateCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendButtonInteractionTemplateCardLogic) SendButtonInteractionTemplateCard(req *types.SendButtonInteractionTemplateCardReq) (*types.SendButtonInteractionTemplateCardResp, error) {
	// 获取企业微信 Access Token
	stableAccessToken := credential.NewWorkAccessToken(l.svcCtx.Config.Midjourney.CorpId, l.svcCtx.Config.Midjourney.CorpSecret, cacheKeyPrefix, cache.NewRedis(l.ctx, &cache.RedisOpts{
		Host:     l.svcCtx.Config.Midjourney.Cache.Host,
		Password: l.svcCtx.Config.Midjourney.Cache.Pass,
	}))
	accessToken, err := stableAccessToken.GetAccessToken()
	if err != nil {
		return nil, err
	}
	// 请求参数转 JSON 格式
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	// 发起http请求
	response, err := util.HTTPPost(fmt.Sprintf(sendButtonInteractionTemplateCardURL, accessToken), string(jsonData))
	if err != nil {
		return nil, err
	}

	var sendButtonInteractionTemplateCardResp *types.SendButtonInteractionTemplateCardResp

	err = json.Unmarshal(response, &sendButtonInteractionTemplateCardResp)
	if err != nil {
		return nil, err
	}
	return sendButtonInteractionTemplateCardResp, nil
}
