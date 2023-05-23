package svc

import (
	"github.com/ParkerWen/wechat_pro/app/midjourney/task/model"
	"github.com/ParkerWen/wechat_pro/app/mqueue/job/internal/config"
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server
	TaskModel   model.TaskModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		AsynqServer: newAsynqServer(c),
		TaskModel:   model.NewTaskModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
