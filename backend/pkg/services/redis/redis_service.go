package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"webGL-720yun/config"

	"github.com/go-redis/redis"
)

// RedisService Redis服务
type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisService 创建Redis服务
func NewRedisService(config *config.RedisConfig) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	ctx := context.Background()

	// 测试连接
	if err := client.Ping().Err(); err != nil {
		// Redis连接失败，记录错误但不中断程序
		fmt.Printf("警告: Redis连接失败: %v\n", err)
		fmt.Println("程序将继续运行，但缓存功能将不可用")
		// 返回一个空的RedisService实例，后续操作会失败但不会panic
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

// 键名前缀定义
const (
	KeyPrefixUserSession    = "user:session:"    // 用户会话
	KeyPrefixUserOnline     = "user:online:"     // 用户在线状态
	KeyPrefixTokenBlacklist = "token:blacklist:" // 令牌黑名单
	KeyPrefixUserInfo       = "user:info:"       // 用户信息缓存
)

// SaveUserSession 保存用户会话
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

// GetUserSession 获取用户会话
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

// DeleteUserSession 删除用户会话
func (s *RedisService) DeleteUserSession(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserSession, userID)
	return s.client.Del(key).Err()
}

// SetUserOnline 设置用户在线状态
func (s *RedisService) SetUserOnline(userID uint, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	return s.client.Set(key, "1", expiresIn).Err()
}

// IsUserOnline 检查用户是否在线
func (s *RedisService) IsUserOnline(userID uint) bool {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	err := s.client.Get(key).Err()
	return err == nil
}

// SetUserOffline 设置用户离线
func (s *RedisService) SetUserOffline(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserOnline, userID)
	return s.client.Del(key).Err()
}

// AddToBlacklist 将令牌加入黑名单
func (s *RedisService) AddToBlacklist(token string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%s", KeyPrefixTokenBlacklist, token)
	return s.client.Set(key, "1", expiresIn).Err()
}

// IsInBlacklist 检查令牌是否在黑名单中
func (s *RedisService) IsInBlacklist(token string) bool {
	key := fmt.Sprintf("%s%s", KeyPrefixTokenBlacklist, token)
	err := s.client.Get(key).Err()
	return err == nil
}

// CacheUserInfo 缓存用户信息
func (s *RedisService) CacheUserInfo(userID uint, userInfo interface{}, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	data, err := json.Marshal(userInfo)
	if err != nil {
		return err
	}

	return s.client.Set(key, data, expiresIn).Err()
}

// GetCachedUserInfo 获取缓存的用户信息
func (s *RedisService) GetCachedUserInfo(userID uint) (map[string]interface{}, error) {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	data, err := s.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal([]byte(data), &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// DeleteCachedUserInfo 删除缓存的用户信息
func (s *RedisService) DeleteCachedUserInfo(userID uint) error {
	key := fmt.Sprintf("%s%d", KeyPrefixUserInfo, userID)
	return s.client.Del(key).Err()
}

// GetClient 获取Redis客户端
func (s *RedisService) GetClient() *redis.Client {
	return s.client
}

// Close 关闭连接
func (s *RedisService) Close() error {
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}
