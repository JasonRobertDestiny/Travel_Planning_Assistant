package configs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Config 存储应用配置信息
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	JWT      JWTConfig      `json:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        string `json:"port"`
	Environment string `json:"environment"` // development, production, testing
	Mode        string `json:"mode"`
	Timeout     int    `json:"timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string `json:"secret"`
	ExpiresIn int    `json:"expires_in"` // 过期时间（小时）
}

// 全局配置实例
var AppConfig Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 加载默认配置
	AppConfig = defaultConfig()

	// 如果指定了配置文件路径，则尝试从文件加载
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			return fmt.Errorf("打开配置文件失败: %w", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&AppConfig); err != nil {
			return fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 优先使用环境变量覆盖
	loadFromEnv()

	// 打印当前环境
	log.Printf("应用运行在 %s 环境", AppConfig.Server.Environment)

	return nil
}

// 加载默认配置
func defaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:        "8080",
			Environment: "development",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "traveler",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:    "your-secret-key",
			ExpiresIn: 24,
		},
	}
}

// 从环境变量加载配置
func loadFromEnv() {
	// 服务器配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		AppConfig.Server.Port = port
	}
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		AppConfig.Server.Environment = env
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		AppConfig.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		AppConfig.Database.Port = port
	}
	if user := os.Getenv("DB_USER"); user != "" {
		AppConfig.Database.User = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		AppConfig.Database.Password = password
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		AppConfig.Database.DBName = dbName
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		AppConfig.Database.SSLMode = sslMode
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		AppConfig.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		AppConfig.Redis.Port = port
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		AppConfig.Redis.Password = password
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		AppConfig.JWT.Secret = secret
	}
}

// GetDatabaseDSN 获取数据库连接字符串
func GetDatabaseDSN() string {
	db := AppConfig.Database
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

// GetRedisAddr 获取Redis连接地址
func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", AppConfig.Redis.Host, AppConfig.Redis.Port)
}

// IsProduction 检查是否为生产环境
func IsProduction() bool {
	return strings.ToLower(AppConfig.Server.Environment) == "production"
}

// IsDevelopment 检查是否为开发环境
func IsDevelopment() bool {
	return strings.ToLower(AppConfig.Server.Environment) == "development"
}

// IsTesting 检查是否为测试环境
func IsTesting() bool {
	return strings.ToLower(AppConfig.Server.Environment) == "testing"
}

// PostgresConnectionString 获取PostgreSQL连接字符串
func PostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.DBName,
		AppConfig.Database.SSLMode,
	)
}

// RedisConnectionString 获取Redis连接字符串
func RedisConnectionString() string {
	return fmt.Sprintf(
		"%s:%s",
		AppConfig.Redis.Host,
		AppConfig.Redis.Port,
	)
}

// ServerAddress 获取服务器监听地址
func ServerAddress() string {
	return fmt.Sprintf(":%s", AppConfig.Server.Port)
}

// SetupTimeout 设置超时时间
func SetupTimeout() time.Duration {
	return time.Duration(AppConfig.Server.Timeout) * time.Second
}

// GetJWTExpirationTime 获取JWT过期时间
func GetJWTExpirationTime() time.Duration {
	return time.Duration(AppConfig.JWT.ExpiresIn) * time.Hour
}

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() string {
	return AppConfig.JWT.Secret
}
