# PersonaTrip - AI定制旅游规划系统

基于Go语言和Eino大模型框架的AI定制旅游规划后端系统。

## 项目概述

PersonaTrip是一个利用AI技术为用户提供个性化旅游规划的应用。系统结合大语言模型、用户偏好分析和旅游知识图谱，实现真正"千人千面"的定制旅游服务。

### 核心功能

- 基于用户偏好生成个性化旅游计划
- 提供目的地推荐
- 管理用户旅行计划
- 结合Eino大模型提供智能旅游建议

## 技术栈

- 后端：Go语言
- Web框架：Gin
- 数据库：MongoDB
- AI模型：Eino大模型框架
- API文档：Swagger

## 项目结构

```
PersonaTrip/
├── cmd/                # 命令行入口
├── internal/           # 内部包
│   ├── api/            # API路由定义
│   ├── config/         # 配置管理
│   ├── handlers/       # 请求处理器
│   ├── middleware/     # 中间件
│   ├── models/         # 数据模型
│   ├── repository/     # 数据存储层
│   ├── services/       # 业务逻辑层
│   └── utils/          # 工具函数
├── pkg/                # 可导出的包
├── .env                # 环境变量配置
├── go.mod              # Go模块定义
├── go.sum              # 依赖版本锁定
├── main.go             # 应用入口
└── README.md           # 项目说明
```

## 安装与运行

### 前置条件

- Go 1.21或更高版本
- MongoDB
- Eino API密钥

### 环境变量配置

创建`.env`文件并配置以下环境变量：

```
APP_ENV=development
PORT=8080
MONGO_URI=mongodb://localhost:27017/personatrip
EINO_API_KEY=your_eino_api_key
```

### 运行应用

```bash
# 安装依赖
go mod download

# 运行应用
go run main.go
```

## API文档

启动应用后，可以通过以下URL访问Swagger API文档：

```
http://localhost:8080/swagger/index.html
```

## 主要API端点

- `POST /api/trips/generate` - 生成AI旅行计划
- `GET /api/trips/:id` - 获取旅行计划详情
- `GET /api/trips/user` - 获取用户的所有旅行计划
- `PUT /api/trips/:id` - 更新旅行计划
- `DELETE /api/trips/:id` - 删除旅行计划
- `POST /api/recommendations/destinations` - 获取目的地推荐
