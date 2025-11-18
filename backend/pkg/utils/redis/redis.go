package redis

import (
	"webGL-720yun/config"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 该项目中Redis的实现确实是一种相对简单的方式，通过全局变量模式实现
// 虽然这种实现简单直观，但在扩展性、测试性和维护性方面都不如SVC服务上下文模式灵活，特别是对于更复杂的企业级应用
var client *redis.Client

func Init() {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host + ":" + strconv.Itoa(config.Conf.Redis.Port),
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
}

// Set 存储键值对
func Set(key, value string, expiration time.Duration) {
	client.Set(key, value, expiration)
}

// Get 获取键值对
func Get(key string) (string, error) {
	return client.Get(key).Result()
}

// Del 删除键值对
func Del(key string) bool {
	res := client.Del(key)
	return res.Val() > 0
}

// GetIdleTime 获取键的空闲时间（没有被访问的时间）
func GetIdleTime(key string) (time.Duration, error) {
	duration, err := client.ObjectIdleTime(key).Result()
	if err != nil {
		return 0, err
	}
	return duration, nil
}
