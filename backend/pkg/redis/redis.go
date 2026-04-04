package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"webGL-720yun/config"

	"github.com/go-redis/redis"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService(config *config.RedisConfig) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	ctx := context.Background()

	if err := client.Ping().Err(); err != nil {
		fmt.Printf("警告: Redis连接失败: %v\n", err)
		fmt.Println("程序将继续运行，但缓存功能将不可用")
		return &RedisService{
			client: nil,
			ctx:    ctx,
		}
	}

	return &RedisService{
		client: client,
		ctx:    ctx,
	}
}

const (
	KeyPrefixUserSession    = "user:session:"
	KeyPrefixUserOnline     = "user:online:"
	KeyPrefixTokenBlacklist = "token:blacklist:"
	KeyPrefixUserInfo       = "user:info:"
	KeyPrefixLock           = "lock:"
)

func (s *RedisService) SaveUserSession(userID uint, accessToken, refreshToken string, expiresIn time.Duration) error {
	sessionData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       userID,
		"created_at":    time.Now().Unix(),
	}

	key := fmt.Sprintf("%s%d", KeyPrefixUserSession, userID)
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	return s.client.Set(key, data, expiresIn).Err()
}

func (s *RedisService) GetUserSession(userID uint) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixUserSession, userID)
	data, err := s.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var sessionData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &sessionData); err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (s *RedisService) DeleteUserSession(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserSession, userID)
	return s.client.Del(key).Err()
}

func (s *RedisService) SetUserOnline(userID uint, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	return s.client.Set(key, "1", expiresIn).Err()
}

func (s *RedisService) IsUserOnline(userID uint) bool {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	err := s.client.Get(key).Err()
	return err == nil
}

func (s *RedisService) SetUserOffline(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	return s.client.Del(key).Err()
}

func (s *RedisService) AddToBlacklist(token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%s", KeyPrefixTokenBlacklist, token)
	return s.client.Set(key, "1", expiresIn).Err()
}

func (s *RedisService) IsInBlacklist(token string) bool {
	key := fmt.Sprintf("%s%s", KeyPrefixTokenBlacklist, token)
	err := s.client.Get(key).Err()
	return err == nil
}

func addRandomExpiration(base time.Duration) time.Duration {
	percentage := 0.1 + 0.2*(float64(time.Now().UnixNano()%100)/100.0)
	randomDuration := time.Duration(float64(base) * percentage)
	return base + randomDuration
}

func (s *RedisService) CacheUserInfo(userID uint, userInfo interface{}, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	randomExpiresIn := addRandomExpiration(expiresIn)
	return s.client.Set(key, data, randomExpiresIn).Err()
}

func (s *RedisService) CacheEmptyUserInfo(userID uint, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	emptyUserInfo := make(map[string]interface{})
	data, err := json.Marshal(emptyUserInfo)
	if err != nil {
		return err
	}

	randomExpiresIn := addRandomExpiration(expiresIn)
	return s.client.Set(key, data, randomExpiresIn).Err()
}

func (s *RedisService) GetCachedUserInfo(userID uint) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	data, err := s.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal([]byte(data), &userInfo); err != nil {
		return nil, err
	}

	if len(userInfo) == 0 {
		return nil, nil
	}

	return userInfo, nil
}

func (s *RedisService) DeleteCachedUserInfo(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	return s.client.Del(key).Err()
}

func (s *RedisService) AcquireLock(key string, expiresIn time.Duration) bool {
	lockKey := fmt.Sprintf("%s%s", KeyPrefixLock, key)
	success, err := s.client.SetNX(lockKey, "1", expiresIn).Result()
	if err != nil {
		return false
	}
	return success
}

func (s *RedisService) ReleaseLock(key string) error {
	lockKey := fmt.Sprintf("%s%s", KeyPrefixLock, key)
	return s.client.Del(lockKey).Err()
}

func (s *RedisService) GetClient() *redis.Client {
	return s.client
}

func (s *RedisService) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}
