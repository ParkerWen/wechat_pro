package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	CorpId         string
	AppToken       string
	EncodingAesKey string

	Midjourney struct {
		CorpId     string
		CorpSecret string
		AgentId    string
		Cache      struct {
			Host string
			Pass string
		}
		NotifyHook string
	}
}
