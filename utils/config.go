package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Environment string `json:"environment"` // 环境: development, production
	ServerPort  int    `json:"server_port"` // 服务器端口
	JWTSecret   string `json:"jwt_secret"`  // JWT 密钥
	JWTExpires  int    `json:"jwt_expires"` // JWT 过期时间（小时）

	// 数据库配置
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	DBSSLMode  string `json:"db_ssl_mode"`

	// Redis配置
	RedisHost     string `json:"redis_host"`
	RedisPort     int    `json:"redis_port"`
	RedisPassword string `json:"redis_password"`
	RedisDB       int    `json:"redis_db"`

	// 日志配置
	LogLevel  string `json:"log_level"`
	LogOutput string `json:"log_output"`
}

var cfg Config

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载.env文件
	godotenv.Load()

	cfg.Environment = getEnv("APP_ENV", "development")

	var err error
	cfg.ServerPort, err = strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		cfg.ServerPort = 8080
	}

	cfg.JWTSecret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	cfg.JWTExpires, err = strconv.Atoi(getEnv("JWT_EXPIRES", "24"))
	if err != nil {
		cfg.JWTExpires = 24
	}

	// 数据库配置
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort, err = strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		cfg.DBPort = 5432
	}
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "postgres")
	cfg.DBName = getEnv("DB_NAME", "traveler")
	cfg.DBSSLMode = getEnv("DB_SSL_MODE", "disable")

	// Redis配置
	cfg.RedisHost = getEnv("REDIS_HOST", "localhost")
	cfg.RedisPort, err = strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	if err != nil {
		cfg.RedisPort = 6379
	}
	cfg.RedisPassword = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB, err = strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		cfg.RedisDB = 0
	}

	// 日志配置
	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	cfg.LogOutput = getEnv("LOG_OUTPUT", "console")

	return &cfg, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &cfg
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDBConnString 获取数据库连接字符串
func (c *Config) GetDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// GetRedisConnString 获取Redis连接字符串
func (c *Config) GetRedisConnString() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}

// IsDevelopment 是否为开发环境
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Environment) == "development"
}

// IsProduction 是否为生产环境
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

// GetJWTExpirationTime 获取JWT过期时间
func (c *Config) GetJWTExpirationTime() time.Duration {
	return time.Duration(c.JWTExpires) * time.Hour
}
