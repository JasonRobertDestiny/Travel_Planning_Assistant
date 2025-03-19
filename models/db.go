package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// PostgreSQL驱动
	"traveler_agent/configs"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 获取数据库连接字符串
	dsn := configs.GetDatabaseDSN()

	// 打开数据库连接
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %w", err)
	}

	// 配置连接池
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	log.Println("数据库连接成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("数据库连接已关闭")
	}
}
