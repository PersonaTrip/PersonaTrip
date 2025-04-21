# PersonaTrip API 文档

本文档详细介绍了PersonaTrip系统的所有API接口。

## 目录

- [认证相关](#认证相关)
- [旅行计划相关](#旅行计划相关)
- [目的地推荐相关](#目的地推荐相关)
- [管理员系统相关](#管理员系统相关)
- [模型配置相关](#模型配置相关)

## 基本信息

- **基础URL**: `http://14.103.232.255:32601`
- **认证方式**: JWT令牌（Bearer Token）
- **数据格式**: JSON

所有接口均支持跨域请求（CORS）。

---

## 认证相关

### 注册用户

- **URL**: `/api/auth/register`
- **方法**: `POST`
- **描述**: 注册新用户
- **请求体**:
  ```json
  {
    "username": "用户名",
    "email": "邮箱地址",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "created_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 用户登录

- **URL**: `/api/auth/login`
- **方法**: `POST`
- **描述**: 用户登录并获取JWT令牌
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "token": "JWT令牌",
    "user": {
      "id": 1,
      "username": "用户名",
      "email": "邮箱地址"
    }
  }
  ```

### 获取用户资料

- **URL**: `/api/auth/profile`
- **方法**: `GET`
- **描述**: 获取当前登录用户的资料
- **认证**: 需要JWT令牌
- **响应**:
  ```json
  {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "created_at": "2025-04-21T13:52:02+08:00"
  }
  ```

---

## 旅行计划相关

### 生成旅行计划

- **URL**: `/api/trips/generate`
- **方法**: `POST`
- **描述**: 根据用户偏好生成旅行计划
- **认证**: 需要JWT令牌
- **请求体**:
  ```json
  {
    "destination": "目的地",
    "start_date": "2025-05-01",
    "end_date": "2025-05-05",
    "budget": "中等",
    "travel_style": ["文化", "美食"],
    "accommodation": ["酒店", "民宿"],
    "transportation": ["公共交通", "步行"],
    "activities": ["观光", "购物"],
    "food_preferences": ["当地美食", "素食"],
    "special_requests": "其他特殊要求"
  }
  ```
- **响应**:
  ```json
  {
    "id": "旅行计划ID",
    "title": "旅行计划标题",
    "destination": "目的地",
    "start_date": "2025-05-01",
    "end_date": "2025-05-05",
    "days": [
      {
        "day": 1,
        "date": "2025-05-01",
        "activities": [
          {
            "name": "活动名称",
            "type": "活动类型",
            "location": {
              "name": "地点名称",
              "address": "地址",
              "city": "城市",
              "country": "国家"
            },
            "start_time": "09:00",
            "end_time": "11:00",
            "description": "活动描述",
            "cost": 100
          }
        ],
        "meals": [
          {
            "type": "午餐",
            "venue": "餐厅名称",
            "description": "描述",
            "cost": 50
          }
        ],
        "accommodation": {
          "name": "住宿名称",
          "type": "酒店",
          "description": "描述",
          "cost": 300
        }
      }
    ],
    "budget": {
      "currency": "CNY",
      "total_estimate": 2000,
      "accommodation": 800,
      "transportation": 300,
      "food": 500,
      "activities": 300,
      "other": 100
    },
    "notes": "额外注意事项",
    "created_at": "2025-04-21T13:52:02+08:00",
    "user_id": 1
  }
  ```

### 获取旅行计划

- **URL**: `/api/trips/:id`
- **方法**: `GET`
- **描述**: 获取指定ID的旅行计划
- **参数**: 
  - `id`: 旅行计划ID
- **响应**: 与生成旅行计划接口的响应格式相同

### 获取用户的所有旅行计划

- **URL**: `/api/trips/user`
- **方法**: `GET`
- **描述**: 获取当前登录用户的所有旅行计划
- **认证**: 需要JWT令牌
- **响应**:
  ```json
  [
    {
      "id": "旅行计划ID",
      "title": "旅行计划标题",
      "destination": "目的地",
      "start_date": "2025-05-01",
      "end_date": "2025-05-05",
      "created_at": "2025-04-21T13:52:02+08:00"
    }
  ]
  ```

### 更新旅行计划

- **URL**: `/api/trips/:id`
- **方法**: `PUT`
- **描述**: 更新指定ID的旅行计划
- **认证**: 需要JWT令牌
- **参数**: 
  - `id`: 旅行计划ID
- **请求体**: 与生成旅行计划接口的请求体格式相同
- **响应**: 与生成旅行计划接口的响应格式相同

### 删除旅行计划

- **URL**: `/api/trips/:id`
- **方法**: `DELETE`
- **描述**: 删除指定ID的旅行计划
- **认证**: 需要JWT令牌
- **参数**: 
  - `id`: 旅行计划ID
- **响应**:
  ```json
  {
    "message": "旅行计划已删除"
  }
  ```

---

## 目的地推荐相关

### 生成目的地推荐

- **URL**: `/api/recommendations/destinations`
- **方法**: `POST`
- **描述**: 根据用户偏好生成目的地推荐
- **请求体**:
  ```json
  {
    "interests": ["文化", "美食", "自然"],
    "budget": "中等",
    "duration": "5-7天",
    "season": "夏季",
    "travel_style": ["休闲", "探险"]
  }
  ```
- **响应**:
  ```json
  {
    "destinations": [
      {
        "name": "目的地名称",
        "country": "国家",
        "description": "描述",
        "highlights": ["亮点1", "亮点2"],
        "best_time_to_visit": "最佳旅游时间",
        "estimated_cost": "估计费用"
      }
    ]
  }
  ```

---

## 管理员系统相关

### 管理员登录

- **URL**: `/api/admin/login`
- **方法**: `POST`
- **描述**: 管理员登录并获取JWT令牌
- **请求体**:
  ```json
  {
    "username": "admin",
    "password": "gr123465!"
  }
  ```
- **响应**:
  ```json
  {
    "token": "JWT令牌",
    "admin": {
      "id": 1,
      "username": "admin",
      "email": "admin@personatrip.com",
      "role": "super_admin"
    }
  }
  ```

### 创建管理员

- **URL**: `/api/admin/admins`
- **方法**: `POST`
- **描述**: 创建新的管理员（仅超级管理员可访问）
- **认证**: 需要管理员JWT令牌
- **请求体**:
  ```json
  {
    "username": "新管理员用户名",
    "password": "密码",
    "email": "邮箱地址",
    "role": "admin"
  }
  ```
- **响应**:
  ```json
  {
    "id": 2,
    "username": "新管理员用户名",
    "email": "邮箱地址",
    "role": "admin",
    "created_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 获取所有管理员

- **URL**: `/api/admin/admins`
- **方法**: `GET`
- **描述**: 获取所有管理员（仅超级管理员可访问）
- **认证**: 需要管理员JWT令牌
- **响应**:
  ```json
  [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@personatrip.com",
      "role": "super_admin",
      "created_at": "2025-04-21T13:52:02+08:00"
    },
    {
      "id": 2,
      "username": "新管理员用户名",
      "email": "邮箱地址",
      "role": "admin",
      "created_at": "2025-04-21T13:52:02+08:00"
    }
  ]
  ```

### 获取管理员详情

- **URL**: `/api/admin/admins/:id`
- **方法**: `GET`
- **描述**: 获取指定ID的管理员详情（仅超级管理员可访问）
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 管理员ID
- **响应**:
  ```json
  {
    "id": 2,
    "username": "管理员用户名",
    "email": "邮箱地址",
    "role": "admin",
    "created_at": "2025-04-21T13:52:02+08:00",
    "updated_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 更新管理员

- **URL**: `/api/admin/admins/:id`
- **方法**: `PUT`
- **描述**: 更新指定ID的管理员信息（仅超级管理员可访问）
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 管理员ID
- **请求体**:
  ```json
  {
    "email": "新邮箱地址",
    "password": "新密码",
    "role": "admin"
  }
  ```
- **响应**:
  ```json
  {
    "id": 2,
    "username": "管理员用户名",
    "email": "新邮箱地址",
    "role": "admin",
    "updated_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 删除管理员

- **URL**: `/api/admin/admins/:id`
- **方法**: `DELETE`
- **描述**: 删除指定ID的管理员（仅超级管理员可访问）
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 管理员ID
- **响应**:
  ```json
  {
    "message": "管理员已删除"
  }
  ```

---

## 模型配置相关

### 创建模型配置

- **URL**: `/api/admin/models`
- **方法**: `POST`
- **描述**: 创建新的模型配置
- **认证**: 需要管理员JWT令牌
- **请求体**:
  ```json
  {
    "name": "模型配置名称",
    "model_type": "openai",
    "model_name": "gpt-4",
    "api_key": "API密钥",
    "base_url": "基础URL",
    "is_active": false,
    "temperature": 0.7,
    "max_tokens": 2000
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "模型配置名称",
    "model_type": "openai",
    "model_name": "gpt-4",
    "api_key": "API密钥",
    "base_url": "基础URL",
    "is_active": false,
    "temperature": 0.7,
    "max_tokens": 2000,
    "created_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 获取所有模型配置

- **URL**: `/api/admin/models`
- **方法**: `GET`
- **描述**: 获取所有模型配置
- **认证**: 需要管理员JWT令牌
- **响应**:
  ```json
  [
    {
      "id": 1,
      "name": "OpenAI GPT-4",
      "model_type": "openai",
      "model_name": "gpt-4",
      "api_key": "API密钥",
      "base_url": "基础URL",
      "is_active": true,
      "temperature": 0.7,
      "max_tokens": 2000,
      "created_at": "2025-04-21T13:52:02+08:00"
    },
    {
      "id": 2,
      "name": "Ollama本地模型",
      "model_type": "ollama",
      "model_name": "llama3",
      "base_url": "http://localhost:11434",
      "is_active": false,
      "temperature": 0.8,
      "max_tokens": 1500,
      "created_at": "2025-04-21T13:52:02+08:00"
    }
  ]
  ```

### 获取活跃的模型配置

- **URL**: `/api/admin/models/active`
- **方法**: `GET`
- **描述**: 获取当前活跃的模型配置
- **认证**: 需要管理员JWT令牌
- **响应**:
  ```json
  {
    "id": 1,
    "name": "OpenAI GPT-4",
    "model_type": "openai",
    "model_name": "gpt-4",
    "api_key": "API密钥",
    "base_url": "基础URL",
    "is_active": true,
    "temperature": 0.7,
    "max_tokens": 2000,
    "created_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 获取模型配置详情

- **URL**: `/api/admin/models/:id`
- **方法**: `GET`
- **描述**: 获取指定ID的模型配置详情
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 模型配置ID
- **响应**:
  ```json
  {
    "id": 1,
    "name": "OpenAI GPT-4",
    "model_type": "openai",
    "model_name": "gpt-4",
    "api_key": "API密钥",
    "base_url": "基础URL",
    "is_active": true,
    "temperature": 0.7,
    "max_tokens": 2000,
    "created_at": "2025-04-21T13:52:02+08:00",
    "updated_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 更新模型配置

- **URL**: `/api/admin/models/:id`
- **方法**: `PUT`
- **描述**: 更新指定ID的模型配置
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 模型配置ID
- **请求体**:
  ```json
  {
    "name": "更新后的名称",
    "api_key": "新的API密钥",
    "temperature": 0.8,
    "max_tokens": 1800
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "name": "更新后的名称",
    "model_type": "openai",
    "model_name": "gpt-4",
    "api_key": "新的API密钥",
    "base_url": "基础URL",
    "is_active": true,
    "temperature": 0.8,
    "max_tokens": 1800,
    "updated_at": "2025-04-21T13:52:02+08:00"
  }
  ```

### 删除模型配置

- **URL**: `/api/admin/models/:id`
- **方法**: `DELETE`
- **描述**: 删除指定ID的模型配置
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 模型配置ID
- **响应**:
  ```json
  {
    "message": "模型配置已删除"
  }
  ```

### 设置活跃模型配置

- **URL**: `/api/admin/models/:id/activate`
- **方法**: `POST`
- **描述**: 将指定ID的模型配置设置为活跃
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 模型配置ID
- **响应**:
  ```json
  {
    "message": "模型配置已设置为活跃"
  }
  ```

### 测试模型配置

- **URL**: `/api/admin/models/:id/test`
- **方法**: `POST`
- **描述**: 测试指定ID的模型配置
- **认证**: 需要管理员JWT令牌
- **参数**: 
  - `id`: 模型配置ID
- **请求体**:
  ```json
  {
    "prompt": "测试提示词"
  }
  ```
- **响应**:
  ```json
  {
    "result": "模型生成的结果文本"
  }
  ```

---

## 错误响应

所有API在发生错误时会返回相应的HTTP状态码和错误信息：

```json
{
  "error": "错误信息"
}
```

常见的HTTP状态码：

- `200 OK`: 请求成功
- `201 Created`: 资源创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或认证失败
- `403 Forbidden`: 没有权限访问资源
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误
