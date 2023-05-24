package svc

import (
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/api/internal/config"
	"github.com/ParkerWen/wechat_pro/app/midjourney/user/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	UserModel model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
