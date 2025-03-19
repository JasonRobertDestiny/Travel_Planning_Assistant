package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"traveler_agent/configs"
	"traveler_agent/models"
	"traveler_agent/routers"
	"traveler_agent/utils"
)

func main() {
	// 定义配置文件路径标志
	configPath := flag.String("config", "configs/config.json", "配置文件路径")
	flag.Parse()

	// 加载配置
	if err := configs.LoadConfig(*configPath); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 尝试初始化数据库连接
	// 数据库启动失败不会阻止应用启动，因为在开发环境中可能没有数据库
	if err := models.InitDB(); err != nil {
		log.Printf("警告: 数据库初始化失败: %v", err)
	} else {
		defer models.CloseDB()
	}

	// 尝试初始化Redis连接
	// Redis启动失败不会阻止应用启动，因为Redis主要用于缓存
	if err := utils.InitRedis(); err != nil {
		log.Printf("警告: Redis初始化失败: %v", err)
	} else {
		defer utils.CloseRedis()
	}

	// 设置路由
	r := routers.SetupRouter()

	// 在单独的协程中启动服务器
	go func() {
		port := configs.AppConfig.Server.Port
		if port == "" {
			port = "8080" // 默认端口
		}

		log.Printf("服务器启动在 http://localhost:%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 处理优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 这里可以添加其他清理工作
	log.Println("服务器已关闭")
}
