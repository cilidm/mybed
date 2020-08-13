package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"mybedv2/conf"
	"sync"
)

var (
	Client *redis.Client
	once   sync.Once
)

func init() {
	once.Do(func() {
		Client = redis.NewClient(&redis.Options{
			Addr:     conf.Setting.RedisAddr,
			Password: conf.Setting.RedisPWD, // no password set
			DB:       conf.Setting.RedisDB,  // use default DB
		})
		pong, err := Client.Ping().Result()
		fmt.Println(pong, err)
	})
}
