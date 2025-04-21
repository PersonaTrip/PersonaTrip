# PersonaTrip - AI定制旅游规划系统 | AI-Powered Custom Travel Planning System

[English Version](#english-version) | [中文版本](#中文版本)

---

## English Version

AI-powered custom travel planning backend system based on Go language and Eino large language model framework.

### Project Overview

PersonaTrip is an application that uses AI technology to provide personalized travel planning for users. The system combines large language models, user preference analysis, and travel knowledge graphs to achieve truly personalized custom travel services.

### Core Features

- Generate personalized travel plans based on user preferences
- Provide destination recommendations
- Manage user travel plans
- Offer intelligent travel suggestions using multiple large language models
- User authentication system (login, registration)

### Technology Stack

- Backend: Go
- Web Framework: Gin
- Databases:
  - MongoDB (travel data)
  - MySQL (user authentication data)
- Authentication: JWT (JSON Web Token)
- AI Integration: 
  - Eino framework with support for multiple LLMs
  - Compatible with OpenAI, Ollama (local models), Ark, and more
  - Flexible model selection at runtime
- API Documentation: Swagger
- Containerization: Docker

### Project Structure

```
PersonaTrip/
├── cmd/                # Command line entry points
│   └── root.go         # Main command implementation, server startup and graceful shutdown
├── internal/           # Internal packages
│   ├── api/            # API route definitions
│   │   └── routes.go   # Registration of all API routes
│   ├── config/         # Configuration management
│   │   └── config.go   # Application configuration structure and loading logic
│   ├── handlers/       # HTTP request handlers
│   │   ├── auth_handlers.go  # Authentication-related handlers
│   │   └── trip_handlers.go  # Travel plan-related handlers
│   ├── middleware/     # Middleware
│   │   └── auth.go     # Authentication middleware
│   ├── models/         # Data models
│   │   ├── models.go   # Travel plan-related models
│   │   └── user.go     # User-related models
│   ├── repository/     # Data storage layer
│   │   ├── errors.go           # Storage layer error definitions
│   │   ├── memory_store.go     # In-memory storage implementation (backup)
│   │   ├── mongodb.go          # MongoDB storage implementation
│   │   └── mysql.go            # MySQL storage implementation
│   └── services/       # Business logic layer
│       ├── auth_service.go     # Authentication service
│       └── eino_service.go     # Eino AI service
├── pkg/                # Exportable packages
│   └── einosdk/        # Eino SDK
│       └── einosdk.go  # Eino SDK implementation
├── .env                # Environment variables
├── API接口文档.md        # API documentation in Chinese
├── Dockerfile          # Docker build file
├── go.mod              # Go module definition
├── go.sum              # Dependency version lock
├── main.go             # Application entry point
└── README.md           # Project documentation
```

### Installation and Running

#### Prerequisites

- Go 1.21 or higher
- MongoDB
- MySQL
- Eino API key

#### Environment Variables Configuration

Create a `.env` file and configure the following environment variables:

```
APP_ENV=development
PORT=8080
MONGO_URI=mongodb://your-mongodb-host:port/personatrip
MYSQL_DSN=username:password@tcp(your-mysql-host:port)/personatrip?parseTime=true
JWT_SECRET=your-jwt-secret-key

# LLM Configuration (choose one or more based on your needs)
# OpenAI Configuration
# OPENAI_API_KEY=your-openai-api-key-here

# Ollama Configuration (for local deployment)
# OLLAMA_BASE_URL=http://localhost:11434

# Ark Configuration
# ARK_API_KEY=your-ark-api-key-here
```

#### Running the Application

```bash
# Install dependencies
go mod download

# Run the application
go run main.go

# Or build and run
go build -o personatrip
./personatrip
```

### Docker Deployment

```bash
# Build Docker image
docker build -t personatrip:latest .

# Run container
docker run -p 8080:8080 personatrip:latest
```

### API Documentation

API documentation is available in the `API接口文档.md` file, or you can access the Swagger API documentation after starting the application at:

```
http://localhost:8080/swagger/index.html
```

### Main API Endpoints

#### Travel Plan Related
- `POST /api/trips` - Create a travel plan
- `GET /api/trips/:id` - Get travel plan details
- `GET /api/trips` - Get all travel plans for a user
- `PUT /api/trips/:id` - Update a travel plan
- `DELETE /api/trips/:id` - Delete a travel plan
- `POST /api/trips/:id/ai-suggestions` - Get AI travel suggestions

#### Authentication Related
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/profile` - Get user profile

### Data Flow

1. Client sends HTTP request to API endpoint
2. Request passes through middleware processing (e.g., authentication)
3. Router dispatches the request to the corresponding handler
4. Handler calls service layer methods
5. Service layer interacts with storage layer and AI services
6. AI services communicate with selected LLM providers based on configuration
7. Results are processed and returned to the client

### LLM Integration

PersonaTrip supports multiple large language models through the Eino framework:

#### Supported Models

- **OpenAI**: GPT-3.5, GPT-4, etc.
- **Ollama**: Local deployment of open-source models like Llama, Mistral, etc.
- **Ark**: Cloud-based models from Volcano Engine
- **Mock**: For testing and development

#### Usage Examples

```go
// Using OpenAI
service := NewEinoServiceWithModel(
    einosdk.ModelTypeOpenAI,
    einosdk.WithAPIKey("your-openai-api-key"),
    einosdk.WithModel("gpt-4"),
)

// Using Ollama (locally deployed open-source models)
service := NewEinoServiceWithModel(
    einosdk.ModelTypeOllama,
    einosdk.WithBaseURL("http://localhost:11434"),
    einosdk.WithModel("llama2"),
)

// Using Ark
service := NewEinoServiceWithModel(
    einosdk.ModelTypeArk,
    einosdk.WithAPIKey("your-ark-api-key"),
)

// Using Mock model (for testing)
service := NewEinoServiceWithModel(einosdk.ModelTypeMock)
```

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
