# 第一阶段：构建阶段
FROM golang:1.21-alpine AS builder

# 安装基本依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod ./
COPY go.sum ./

# 下载所有依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o personatrip

# 第二阶段：运行阶段
FROM alpine:latest

# 安装基本工具和CA证书（用于HTTPS请求）
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为亚洲/上海
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN adduser -D -g '' appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的应用
COPY --from=builder /app/personatrip .
# 复制环境配置文件（如果需要）
COPY --from=builder /app/.env .

# 使用非root用户运行应用
USER appuser

# 暴露应用端口
EXPOSE 8080

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# 运行应用
CMD ["./personatrip"]
