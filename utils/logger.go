package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	// Log 全局日志对象
	Log *logrus.Logger
)

// InitLogger 初始化日志系统
func InitLogger(config *Config) {
	Log = logrus.New()

	// 设置日志级别
	level := config.Logging.Level
	switch strings.ToLower(level) {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	format := config.Logging.Format
	if strings.ToLower(format) == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置日志输出
	outputs := strings.Split(config.Logging.Output, ",")
	writers := []io.Writer{}

	for _, output := range outputs {
		switch strings.TrimSpace(strings.ToLower(output)) {
		case "console":
			writers = append(writers, os.Stdout)
		case "file":
			// 确保日志目录存在
			logDir := filepath.Dir(config.Logging.FilePath)
			if _, err := os.Stat(logDir); os.IsNotExist(err) {
				os.MkdirAll(logDir, 0755)
			}

			// 打开日志文件
			file, err := os.OpenFile(config.Logging.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				Log.Fatal("无法打开日志文件: ", err)
			}
			writers = append(writers, file)
		}
	}

	// 如果没有指定输出，默认输出到控制台
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	// 设置多输出
	Log.SetOutput(io.MultiWriter(writers...))

	Log.Info("日志系统初始化完成")
}

// LoggerMiddleware 日志中间件，用于记录每个请求
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		Log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
		}).Info()
	}
}

// 用于测试的日志函数
func LogTest() {
	Log.Debug("这是一条Debug日志")
	Log.Info("这是一条Info日志")
	Log.Warn("这是一条Warn日志")
	Log.Error("这是一条Error日志")
}
