package inits

import (
	"context"
	"fmt"
	"lianxi/srv/dasic/config"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var Rdb *redis.Client

func RedisInit() {
	RedisConfig := config.Gen.Redis
	Add := fmt.Sprintf("%s:%d",
		RedisConfig.Host,
		RedisConfig.Port,
	)
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Add,
		Password: RedisConfig.Password, // no password set
		DB:       RedisConfig.Database, // use default DB
	})
	err := Rdb.Ping(Ctx).Err()
	if err != nil {
		return
	}
	fmt.Println("redis连接成功")
}
