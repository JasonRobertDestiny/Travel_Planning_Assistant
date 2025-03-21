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
	Logging struct {
		Level    string `json:"level"`     // 日志级别: debug, info, warn, error
		Format   string `json:"format"`    // 日志格式: text, json
		Output   string `json:"output"`    // 日志输出: console, file, console,file
		FilePath string `json:"file_path"` // 日志文件路径
	} `json:"logging"`

	// CORS配置
	CORS struct {
		AllowedOrigins   []string `json:"allowed_origins"`
		AllowedMethods   []string `json:"allowed_methods"`
		AllowedHeaders   []string `json:"allowed_headers"`
		ExposeHeaders    []string `json:"expose_headers"`
		AllowCredentials bool     `json:"allow_credentials"`
		MaxAge           int      `json:"max_age"`
	} `json:"cors"`
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
	cfg.Logging.Level = getEnv("LOG_LEVEL", "info")
	cfg.Logging.Format = getEnv("LOG_FORMAT", "text")
	cfg.Logging.Output = getEnv("LOG_OUTPUT", "console")
	cfg.Logging.FilePath = getEnv("LOG_FILE_PATH", "./logs/traveler_agent.log")

	// CORS配置
	cfg.CORS.AllowedOrigins = strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"), ",")
	cfg.CORS.AllowedMethods = strings.Split(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"), ",")
	cfg.CORS.AllowedHeaders = strings.Split(getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"), ",")
	cfg.CORS.ExposeHeaders = strings.Split(getEnv("CORS_EXPOSE_HEADERS", "Content-Length"), ",")
	cfg.CORS.AllowCredentials = getEnv("CORS_ALLOW_CREDENTIALS", "true") == "true"
	cfg.CORS.MaxAge, err = strconv.Atoi(getEnv("CORS_MAX_AGE", "86400"))
	if err != nil {
		cfg.CORS.MaxAge = 86400
	}

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
