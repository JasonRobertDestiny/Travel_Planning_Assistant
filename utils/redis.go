package utils

import (
	"context"
	"log"
	"time"

	"traveler_agent/configs"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	// 获取Redis连接配置
	addr := configs.GetRedisAddr()
	password := configs.AppConfig.Redis.Password
	db := configs.AppConfig.Redis.DB

	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis连接已关闭")
	}
}

// SetCache 设置缓存
func SetCache(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// DeleteCache 删除缓存
func DeleteCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// ClearCache 清除缓存（通过模式匹配）
func ClearCache(ctx context.Context, pattern string) error {
	keys, err := RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}

	return nil
}
