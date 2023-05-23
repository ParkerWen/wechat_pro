package main

import (
	"context"
	"flag"
	"os"

	"github.com/ParkerWen/wechat_pro/app/mqueue/scheduler/internal/config"
	"github.com/ParkerWen/wechat_pro/app/mqueue/scheduler/internal/logic"
	"github.com/ParkerWen/wechat_pro/app/mqueue/scheduler/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	logx.DisableStat()
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	svcContext := svc.NewServiceContext(c)

	mqueueScheduler := logic.NewCornScheduler(ctx, svcContext)
	mqueueScheduler.Register()

	if err := svcContext.Scheduler.Run(); err != nil {
		logx.Errorf("!!!MqueueSchedulerErr!!!  run err:%+v", err)
		os.Exit(1)
	}

}
