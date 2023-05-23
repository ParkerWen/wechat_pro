package svc

import (
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/api/internal/config"
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	TaskModel model.TaskModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		TaskModel: model.NewTaskModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
