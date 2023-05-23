package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	CorpId         string
	AppToken       string
	EncodingAesKey string
}
