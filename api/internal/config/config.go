package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	ShortUrlDB struct {
		DSN string
	}
	SequenceDB struct {
		DSN string
	}
	Redis struct {
		Host string
	}

	ShortBlackList []string

	ShortDoamin string

	CacheRedis cache.CacheConf
}
