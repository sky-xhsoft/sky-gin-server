# 使用官方 Go 镜像作为构建环境
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .

# 下载依赖
RUN go mod download

# 编译为静态二进制
RUN CGO_ENABLED=0 GOOS=linux go build -o sky-server main.go

# 使用极简运行镜像
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/sky-server .

# 配置文件等（如 config.yaml）
COPY config/ ./config/
COPY tmp/ ./tmp/

# 启动服务
CMD ["./sky-server"]