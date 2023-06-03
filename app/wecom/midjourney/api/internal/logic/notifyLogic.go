package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/svc"
	"github.com/ParkerWen/wechat_pro/app/wecom/midjourney/api/internal/types"
	"github.com/silenceper/wechat/util"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	sendURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

type NotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type ImageField struct {
	MediaID string `json:"media_id"`
}

type SendImageReq struct {
	ToUser                 string     `json:"touser"`
	ToParty                string     `json:"toparty"`
	ToTag                  string     `json:"totag"`
	MsgType                string     `json:"msgtype"`
	AgentID                string     `json:"agentid"`
	Image                  ImageField `json:"image"`
	Safe                   int        `json:"safe"`
	EnableIDTrans          int        `json:"enable_id_trans"`
	EnableDuplicateCheck   int        `json:"enable_duplicate_check"`
	DuplicateCheckInterval int        `json:"duplicate_check_interval"`
}

type SendImageResp struct {
	ErrCode        int64  `json:"errcode"`
	ErrMsg         string `json:"errmsg"`
	InvalidUser    string `json:"invaliduser"`    // 不合法的userid，不区分大小写，统一转为小写
	InvalidParty   string `json:"invalidparty"`   // 不合法的partyid
	InvalidTag     string `json:"invalidtag"`     // 不合法的标签id
	UnlicensedUser string `json:"unlicenseduser"` // 没有基础接口许可(包含已过期)的userid
	MsgID          string `json:"msgid"`          // 消息id
	ResponseCode   string `json:"response_code"`
}

func NewNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyLogic {
	return &NotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotifyLogic) Notify(req *types.NotifyReq) (*types.NotifyResp, error) {
	if req.Status != "SUCCESS" {
		return nil, errors.New("图片还未处理完")
	}
	// 获取企业微信 Access Token
	stableAccessToken := credential.NewWorkAccessToken(l.svcCtx.Config.Midjourney.CorpId, l.svcCtx.Config.Midjourney.CorpSecret, cacheKeyPrefix, cache.NewRedis(l.ctx, &cache.RedisOpts{
		Host:     l.svcCtx.Config.Midjourney.Cache.Host,
		Password: l.svcCtx.Config.Midjourney.Cache.Pass,
	}))
	accessToken, err := stableAccessToken.GetAccessToken()
	if err != nil {
		return nil, err
	}
	imgUrl := req.ImageUrl
	url, err := url.Parse(imgUrl)
	if err != nil {
		return nil, err
	}
	url.Scheme = "http"
	url.Host = "img.itcity.cc"
	imgUrl = url.String()
	// 上传临时素材
	var uploadTempFileResp *types.UploadTempFileResp
	upload := NewUploadTempFileLogic(l.ctx, l.svcCtx)
	uploadTempFileResp, err = upload.UploadTempFile(&types.UploadTempFileReq{
		AccessToken: accessToken,
		Media:       imgUrl,
	})
	if err != nil {
		return nil, err
	}
	mediaID := uploadTempFileResp.MediaId
	// 构造请求参数
	var sendImageReq SendImageReq
	sendImageReq.ToUser = req.State
	sendImageReq.ToParty = ""
	sendImageReq.ToTag = ""
	sendImageReq.MsgType = "image"
	sendImageReq.AgentID = l.svcCtx.Config.Midjourney.AgentId
	sendImageReq.Image.MediaID = mediaID
	sendImageReq.Safe = 0
	sendImageReq.EnableDuplicateCheck = 0
	sendImageReq.DuplicateCheckInterval = 1800
	// 请求参数转 JSON 格式
	jsonData, err := json.Marshal(sendImageReq)
	if err != nil {
		return nil, err
	}
	// 发起http请求
	response, err := util.HTTPPost(fmt.Sprintf(sendURL, accessToken), string(jsonData))
	if err != nil {
		return nil, err
	}

	var sendImageResp SendImageResp

	err = json.Unmarshal(response, &sendImageResp)
	if err != nil {
		return nil, err
	}

	return &types.NotifyResp{
		ErrCode: int(sendImageResp.ErrCode),
		ErrMsg:  sendImageResp.ErrMsg,
	}, nil
}
