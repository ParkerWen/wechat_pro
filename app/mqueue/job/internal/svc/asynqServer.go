package svc

import (
	"fmt"

	"github.com/ParkerWen/wechat_pro/app/mqueue/job/internal/config"
	"github.com/hibiken/asynq"
)

func newAsynqServer(c config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Host,
			Password: c.Redis.Pass,
		},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>>  err : %+v \n", err)
				return true
			},
			Concurrency: 20,
		},
	)
}
