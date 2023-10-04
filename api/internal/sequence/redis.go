package sequence

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"

	"strconv"
)

type RedisSequence struct {
	conn *redis.Client
}

func NewRedisSequence(Host string) Sequence {

	redisConn := redis.NewClient(&redis.Options{Addr: Host})

	return &RedisSequence{conn: redisConn}
}

var once sync.Once

func (s *RedisSequence) Next() (uint64, error) {
	once.Do(func() {
		_, err := s.conn.Set("stub", "1", 0).Result()
		fmt.Printf(err.Error())
	})

	n, err := s.conn.Get("stub").Result()
	if err != nil {
		logx.Errorf("redis get sequence failed,err:", err.Error())
		return 0, err
	}
	s.conn.Incr("stub").Result()
	uintn, _ := strconv.Atoi(n)
	return uint64(uintn), nil
}
