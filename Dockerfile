# 使用官方的 Go 语言镜像作为构建阶段基础镜像
FROM golang:1.22.5 AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录，并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 编译 Go 应用
RUN CGO_ENABLED=0 GOOS=linux go build -o myblog ./cmd/main.go

# 第二阶段：运行阶段
FROM alpine:3.18

# 安装证书以支持 HTTPS 请求（如需）
RUN apk --no-cache add ca-certificates

# 创建非 root 用户，提升容器安全性
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 将编译好的二进制文件复制到新镜像中
COPY --from=builder /app/myblog /myblog

# 添加执行权限
RUN chmod +x /myblog

# 切换到非 root 用户
USER appuser

# 暴露端口（根据应用情况修改）
EXPOSE 1323

# 设置容器启动时执行的命令
CMD ["/myblog"]
