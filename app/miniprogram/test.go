package main

import (
	"fmt"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

func main() {
	wc := wechat.NewWechat()
	memory := cache.NewMemory()
	cfg := &miniConfig.Config{
		AppID:     "wx279a9281699dc592",
		AppSecret: "bf2e5f9e2e2d6bf198dfab4d6073b910",
		Cache:     memory,
	}
	miniprogram := wc.GetMiniProgram(cfg)
	res, err := miniprogram.GetAnalysis().GetAnalysisDailyRetain("20230519", "20230519")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
