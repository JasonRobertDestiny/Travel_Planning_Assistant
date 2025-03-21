# 使用官方Go镜像作为构建环境
FROM golang:1.19-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o traveler_agent .

# 使用小型基础镜像运行应用
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/traveler_agent .
# 复制配置文件
COPY --from=builder /app/configs/config.json ./configs/config.json

# 暴露端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# 运行应用
CMD ["./traveler_agent"] 