package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"traveler_agent/models"
	"traveler_agent/routers"
	"traveler_agent/utils"
)

func main() {
	// 定义配置文件路径标志
	configPath := flag.String("config", "configs/config.json", "配置文件路径")
	flag.Parse()

	// 加载配置
	// 使用命令行指定的配置文件路径
	log.Printf("使用配置文件: %s", *configPath)
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	utils.InitLogger(config)
	utils.Log.Info("日志系统初始化完成")
	utils.Log.Infof("使用配置文件: %s", *configPath)
	utils.Log.Info("旅行者助手后端服务启动中...")

	// 尝试初始化数据库连接
	// 数据库启动失败不会阻止应用启动，因为在开发环境中可能没有数据库
	if err := models.InitDB(); err != nil {
		utils.Log.Warnf("警告: 数据库初始化失败: %v", err)
	} else {
		utils.Log.Info("数据库连接成功")
		defer models.CloseDB()
	}

	// 尝试初始化Redis连接
	// Redis启动失败不会阻止应用启动，因为Redis主要用于缓存
	if err := utils.InitRedis(); err != nil {
		utils.Log.Warnf("警告: Redis初始化失败: %v", err)
	} else {
		utils.Log.Info("Redis连接成功")
		defer utils.CloseRedis()
	}

	// 设置路由
	r := routers.SetupRouter()

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
		Handler: r,
	}

	// 在单独的协程中启动服务器
	go func() {
		utils.Log.Infof("服务器启动在 http://localhost:%d", config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 处理优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	utils.Log.Info("正在关闭服务器...")

	// 创建一个5秒的上下文用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		utils.Log.Fatalf("服务器强制关闭: %v", err)
	}

	utils.Log.Info("服务器已关闭")
}
