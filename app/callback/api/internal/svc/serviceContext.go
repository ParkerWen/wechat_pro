package svc

import (
	"github.com/ParkerWen/wechat_pro/app/callback/api/internal/config"
	"github.com/go-laoji/wxbizmsgcrypt"
)

type ServiceContext struct {
	Config config.Config
	WxBiz  *wxbizmsgcrypt.WXBizMsgCrypt
}

func NewServiceContext(c config.Config) *ServiceContext {
	wxbiz := wxbizmsgcrypt.NewWXBizMsgCrypt(c.AppToken, c.EncodingAesKey, c.CorpId, wxbizmsgcrypt.XmlType)
	return &ServiceContext{
		Config: c,
		WxBiz:  wxbiz,
	}
}
