package svc

import (
	"github.com/zeromicro/go-zero/core/bloom"
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortner/api/internal/config"
	"shortner/api/internal/sequence"
	"shortner/model"
)

type ServiceContext struct {
	Config                 config.Config
	ShortUrlModel          model.ShortUrlMapModel
	Sequence               sequence.Sequence   //发号器接口
	ShortBlackMap          map[string]struct{} //黑名单列表
	ShortDoamin            string              //短链接域名
	UserVisitResourceTotal metric.CounterVec   //prometheus 计数
	//bloom filter
	Filter *bloom.Filter
}

func NewServiceContext(c config.Config) *ServiceContext {

	total := metric.NewCounterVec(&metric.CounterVecOpts{ //Counter指标
		Namespace: "user",
		Subsystem: "visit_resource",
		Name:      "total",
		Help:      "user visit shortUrl 127.0.0.1:8888/x count",
		Labels:    []string{"longUrl", "shortUrl"},
	})

	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	shortBlackMap := make(map[string]struct{})

	for _, v := range c.ShortBlackList {
		shortBlackMap[v] = struct{}{}
	}

	cache := c.CacheRedis
	//初始化redisBitSet
	store := redis.New(c.CacheRedis[0].Host, func(r *redis.Redis) {
		r.Type = redis.NodeType
	})
	bitSet := bloom.New(store, "bloom_filter", 1<<20)

	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewShortUrlMapModel(conn, cache),
		Sequence:      sequence.NewMysql(c.SequenceDB.DSN),
		//Sequence:      sequence.NewRedisSequence(c.Redis.Host),
		ShortBlackMap:          shortBlackMap,
		ShortDoamin:            c.ShortDoamin,
		Filter:                 bitSet,
		UserVisitResourceTotal: total,
	}
}
