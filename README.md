# PersonaTrip - AI定制旅游规划系统

## 项目概述

PersonaTrip是一个基于Go语言和Eino大模型框架开发的AI定制旅游规划系统。该系统结合大语言模型、用户偏好分析和旅游知识图谱，为用户提供真正个性化的定制旅游服务。

## 核心功能

- 基于用户偏好生成个性化旅游计划
- 提供目的地推荐
- 管理用户旅行计划
- 结合多种大语言模型提供智能旅游建议
- 用户认证系统（登录、注册）
- 管理员后台系统（模型配置和管理）


## 技术栈

- 后端：Go语言
- Web框架：Gin
- 数据库：
  - MongoDB（旅行数据）
  - MySQL（用户认证数据和管理员系统）
- 认证：JWT (JSON Web Token)
- AI集成：
  - Eino框架（支持多种大模型）
  - 兼容OpenAI、Ollama（本地模型）、Ark等多种模型
  - 支持管理员后台配置模型
- API文档：Swagger
- 容器化：Docker

## 项目结构

```
personatrip/
├── cmd/                # 命令行入口
│   └── root.go         # 主服务器初始化
├── internal/           # 内部包
│   ├── api/            # API路由定义
│   │   ├── routes.go      # 主API路由
│   │   └── admin_routes.go # 管理员API路由
│   ├── config/         # 配置管理
│   │   └── config.go      # 应用配置
│   ├── handlers/       # 请求处理器
│   │   ├── auth_handler.go    # 认证处理器
│   │   ├── trip_handler.go    # 旅行计划处理器
│   │   ├── admin_handler.go   # 管理员处理器
│   │   └── model_config_handler.go # 模型配置处理器
│   ├── middleware/     # 中间件组件
│   │   ├── auth.go        # 认证中间件
│   │   └── admin_auth.go  # 管理员认证中间件
│   ├── models/         # 数据模型
│   │   ├── user.go        # 用户模型
│   │   ├── trip.go        # 旅行模型
│   │   ├── admin.go       # 管理员模型
│   │   └── model_config.go # 模型配置模型
│   ├── repository/     # 数据存储层
│   │   ├── mongodb.go     # MongoDB存储
│   │   ├── mysql.go       # MySQL存储
│   │   ├── admin_repository.go # 管理员存储
│   │   └── model_config_repository.go # 模型配置存储
│   └── services/       # 业务逻辑层
│       ├── auth_service.go     # 认证服务
│       ├── eino_service.go     # Eino AI服务
│       ├── admin_service.go    # 管理员服务
│       └── model_config_service.go # 模型配置服务
├── pkg/                # 可导出的包
│   └── einosdk/        # Eino SDK
│       └── einosdk.go  # Eino SDK实现
├── .env                # 环境变量
├── API接口文档.md        # API文档（中文）
├── Dockerfile          # Docker构建文件
├── go.mod              # Go模块定义
├── go.sum              # 依赖版本锁
├── main.go             # 应用入口点
└── README.md           # 项目文档
```

## 数据流

1. 客户端发送HTTP请求到API端点
2. 请求经过中间件处理（如认证）
3. 路由将请求分发到对应的处理器
4. 处理器调用服务层方法
5. 服务层与存储层和AI服务交互
6. AI服务根据配置与选定的大模型提供商通信
7. 处理结果并返回给客户端

## 大模型集成

PersonaTrip通过Eino框架支持多种大语言模型，并提供管理员后台进行配置和管理：

### 支持的模型

- **OpenAI**：GPT-3.5、GPT-4等
- **Ollama**：本地部署的开源模型，如Llama、Mistral等
- **Ark**：火山引擎提供的云端模型
- **Mock**：用于测试和开发

### 环境变量配置

在`.env`文件中配置以下环境变量：

```
APP_ENV=development
PORT=8080
MONGO_URI=mongodb://your-mongodb-host:port/personatrip
MYSQL_DSN=username:password@tcp(your-mysql-host:port)/personatrip?parseTime=true
JWT_SECRET=your-jwt-secret-key

# 自动迁移数据库和创建超级管理员
# AUTO_MIGRATE=true
# CREATE_SUPER_ADMIN=true
# SUPER_ADMIN_USERNAME=admin
# SUPER_ADMIN_PASSWORD=admin123
# SUPER_ADMIN_EMAIL=admin@personatrip.com

# 大模型配置（可选，优先使用数据库配置）
# OpenAI配置
# OPENAI_API_KEY=your-openai-api-key-here

# Ollama配置（本地部署）
# OLLAMA_BASE_URL=http://localhost:11434

# Ark配置
# ARK_API_KEY=your-ark-api-key-here
```

## 管理员系统

系统包含一个完整的管理员后台，用于管理和配置大模型。

### 管理员角色

- **超级管理员**（super_admin）：可以管理其他管理员和所有模型配置
- **普通管理员**（admin）：可以管理模型配置

### 管理员API端点

- `POST /api/admin/login` - 管理员登录

#### 管理员管理（需要超级管理员权限）

- `POST /api/admin/admins` - 创建新管理员
- `GET /api/admin/admins` - 获取所有管理员
- `GET /api/admin/admins/:id` - 获取特定管理员
- `PUT /api/admin/admins/:id` - 更新管理员
- `DELETE /api/admin/admins/:id` - 删除管理员

#### 模型配置管理

- `POST /api/admin/models` - 创建新的模型配置
- `GET /api/admin/models` - 获取所有模型配置
- `GET /api/admin/models/active` - 获取当前活跃的模型配置
- `GET /api/admin/models/:id` - 获取特定模型配置
- `PUT /api/admin/models/:id` - 更新模型配置
- `DELETE /api/admin/models/:id` - 删除模型配置
- `POST /api/admin/models/:id/activate` - 设置指定模型为活跃
- `POST /api/admin/models/:id/test` - 测试指定模型

### 模型配置字段

每个模型配置包含以下字段：

- **名称**：配置的显示名称
- **模型类型**：openai、ollama、ark或mock
- **模型名称**：具体的模型名称（如gpt-4、llama2等）
- **API密钥**：如果需要，提供模型的API密钥
- **基础URL**：如果需要，提供模型的API基础URL
- **是否活跃**：标记该配置是否当前活跃
- **温度**：生成文本的温度参数
- **最大令牌数**：生成文本的最大令牌数

## 安装和运行

### 前置条件

- Go 1.21或更高版本
- MongoDB
- MySQL

### 运行应用

```bash
# 安装依赖
# 注意：需要添加以下依赖到go.mod
# github.com/golang-jwt/jwt/v4
# gorm.io/gorm
go mod download

# 运行应用
go run main.go

# 或构建并运行
go build -o personatrip
./personatrip
```

### Docker部署

```bash
# 构建Docker镜像
docker build -t personatrip:latest .

# 运行容器
docker run -p 8080:8080 --env-file .env personatrip:latest
```

## API文档

API文档可在`API接口文档.md`文件中查看，或在启动应用后访问：

```
http://localhost:8080/swagger/index.html
```

## 主要API端点

### 旅行计划相关
- `POST /api/trips` - 创建旅行计划
- `GET /api/trips/:id` - 获取旅行计划详情
- `GET /api/trips` - 获取用户的所有旅行计划
- `PUT /api/trips/:id` - 更新旅行计划
- `DELETE /api/trips/:id` - 删除旅行计划
- `POST /api/trips/:id/ai-suggestions` - 获取AI旅行建议

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取用户信息

---

## 中文版本

基于Go语言和Eino大模型框架的AI定制旅游规划后端系统。

### 项目概述

PersonaTrip是一个利用AI技术为用户提供个性化旅游规划的应用。系统结合大语言模型、用户偏好分析和旅游知识图谱，实现真正"千人千面"的定制旅游服务。

### 核心功能

- 基于用户偏好生成个性化旅游计划
- 提供目的地推荐
- 管理用户旅行计划
- 结合多种大语言模型提供智能旅游建议
- 用户认证系统（登录、注册）

### 技术栈

- 后端：Go语言
- Web框架：Gin
- 数据库：
  - MongoDB（旅行数据）
  - MySQL（用户认证数据）
- 认证：JWT (JSON Web Token)
- AI集成：
  - Eino框架（支持多种大模型）
  - 兼容OpenAI、Ollama（本地模型）、Ark等多种模型
  - 支持运行时灵活切换模型
- API文档：Swagger
- 容器化：Docker

### 项目结构

```
PersonaTrip/
├── cmd/                # 命令行入口
│   └── root.go         # 主命令实现，服务器启动和优雅关闭逻辑
├── internal/           # 内部包
│   ├── api/            # API路由定义
│   │   └── routes.go   # 注册所有API路由
│   ├── config/         # 配置管理
│   │   └── config.go   # 应用配置结构和加载逻辑
│   ├── handlers/       # HTTP请求处理器
│   │   ├── auth_handlers.go  # 认证相关处理器
│   │   └── trip_handlers.go  # 旅行计划相关处理器
│   ├── middleware/     # 中间件
│   │   └── auth.go     # 认证中间件
│   ├── models/         # 数据模型
│   │   ├── models.go   # 旅行计划相关模型
│   │   └── user.go     # 用户相关模型
│   ├── repository/     # 数据存储层
│   │   ├── errors.go           # 存储层错误定义
│   │   ├── memory_store.go     # 内存存储实现（备用）
│   │   ├── mongodb.go          # MongoDB存储实现
│   │   └── mysql.go            # MySQL存储实现
│   └── services/       # 业务逻辑层
│       ├── auth_service.go     # 认证服务
│       └── eino_service.go     # Eino AI服务
├── pkg/                # 可导出的包
│   └── einosdk/        # Eino SDK
│       └── einosdk.go  # Eino SDK实现
├── .env                # 环境变量配置
├── API接口文档.md        # API接口中文文档
├── Dockerfile          # Docker构建文件
├── go.mod              # Go模块定义
├── go.sum              # 依赖版本锁定
├── main.go             # 应用入口
└── README.md           # 项目说明
```

### 安装与运行

#### 前置条件

- Go 1.21或更高版本
- MongoDB
- MySQL
- Eino API密钥

#### 环境变量配置

创建`.env`文件并配置以下环境变量：

```
APP_ENV=development
PORT=8080
MONGO_URI=mongodb://your-mongodb-host:port/personatrip
MYSQL_DSN=username:password@tcp(your-mysql-host:port)/personatrip?parseTime=true
JWT_SECRET=your-jwt-secret-key

# 大模型配置（根据需要选择一种或多种）
# OpenAI配置
# OPENAI_API_KEY=your-openai-api-key-here

# Ollama配置（本地部署）
# OLLAMA_BASE_URL=http://localhost:11434

# Ark配置
# ARK_API_KEY=your-ark-api-key-here
```

#### 运行应用

```bash
# 安装依赖
go mod download

# 运行应用
go run main.go

# 或构建后运行
go build -o personatrip
./personatrip
```

### Docker部署

```bash
# 构建Docker镜像
docker build -t personatrip:latest .

# 运行容器
docker run -p 8080:8080 personatrip:latest
```

### API文档

API文档可在`API接口文档.md`文件中查看，或者启动应用后通过以下URL访问Swagger API文档：

```
http://localhost:8080/swagger/index.html
```

### 主要API端点

#### 旅行计划相关
- `POST /api/trips` - 创建旅行计划
- `GET /api/trips/:id` - 获取旅行计划详情
- `GET /api/trips` - 获取用户的所有旅行计划
- `PUT /api/trips/:id` - 更新旅行计划
- `DELETE /api/trips/:id` - 删除旅行计划
- `POST /api/trips/:id/ai-suggestions` - 获取AI旅行建议

#### 用户认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取用户资料

### 数据流

1. 客户端发送HTTP请求到API端点
2. 请求经过中间件处理（如认证）
3. 路由将请求分发到对应的处理器
4. 处理器调用服务层方法
5. 服务层与存储层和AI服务交互
6. AI服务根据配置与选定的大模型提供商通信
7. 处理结果并返回给客户端

### 大模型集成

PersonaTrip通过Eino框架支持多种大语言模型：

#### 支持的模型

- **OpenAI**：GPT-3.5、GPT-4等
- **Ollama**：本地部署的开源模型，如Llama、Mistral等
- **Ark**：火山引擎提供的云端模型
- **Mock**：用于测试和开发

#### 使用示例

```go
// 使用OpenAI
service := NewEinoServiceWithModel(
    einosdk.ModelTypeOpenAI,
    einosdk.WithAPIKey("your-openai-api-key"),
    einosdk.WithModel("gpt-4"),
)

// 使用Ollama（本地部署的开源模型）
service := NewEinoServiceWithModel(
    einosdk.ModelTypeOllama,
    einosdk.WithBaseURL("http://localhost:11434"),
    einosdk.WithModel("llama2"),
)

// 使用Ark
service := NewEinoServiceWithModel(
    einosdk.ModelTypeArk,
    einosdk.WithAPIKey("your-ark-api-key"),
)

// 使用Mock模型（用于测试）
service := NewEinoServiceWithModel(einosdk.ModelTypeMock)
```
