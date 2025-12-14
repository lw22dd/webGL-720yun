# 第一阶段：构建Go应用
FROM golang:1.25.3-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod tidy

# 复制所有代码
COPY . .

# 构建应用
RUN go build -o main ./cmd/main.go

# 第二阶段：运行应用
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 复制构建好的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY config/config.dev.yaml config/config.dev.yaml
COPY config/config.prod.yaml config/config.prod.yaml

# 暴露端口
EXPOSE 7000

# 设置环境变量
ENV GIN_MODE=release

# 运行应用
CMD ["./main"]