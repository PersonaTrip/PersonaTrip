# PersonaTrip API 接口文档

## 概述

PersonaTrip 是一个AI定制旅游规划系统，提供用户认证和旅行计划管理功能。本文档详细说明了系统的API接口。

## 基础信息

- **基础URL**: `http://localhost:8080`
- **认证方式**: JWT Token (Bearer Authentication)
- **内容类型**: application/json

## 用户认证接口

### 用户注册

- **URL**: `/api/auth/register`
- **方法**: `POST`
- **描述**: 创建新用户账户
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "email": "邮箱地址"
  }
  ```
- **成功响应** (状态码: 201):
  ```json
  {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "created_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 400 Bad Request: 请求格式无效
  - 500 Internal Server Error: 服务器内部错误

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
- **成功响应** (状态码: 200):
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "用户名",
      "email": "邮箱地址"
    }
  }
  ```
- **错误响应**:
  - 400 Bad Request: 请求格式无效
  - 401 Unauthorized: 用户名或密码错误

### 获取用户资料

- **URL**: `/api/auth/profile`
- **方法**: `GET`
- **描述**: 获取当前登录用户的资料信息
- **认证**: 需要JWT令牌
- **成功响应** (状态码: 200):
  ```json
  {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "created_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 401 Unauthorized: 未授权访问或令牌无效

## 旅行计划接口

### 创建旅行计划

- **URL**: `/api/trips`
- **方法**: `POST`
- **描述**: 创建新的旅行计划
- **认证**: 需要JWT令牌
- **请求体**:
  ```json
  {
    "destination": "目的地",
    "start_date": "2025-05-01",
    "end_date": "2025-05-07",
    "preferences": {
      "budget": "中等",
      "activities": ["历史景点", "美食", "购物"],
      "accommodation_type": "酒店"
    }
  }
  ```
- **成功响应** (状态码: 201):
  ```json
  {
    "id": "60f7b0b9e6b3f12345678901",
    "user_id": 1,
    "destination": "目的地",
    "start_date": "2025-05-01",
    "end_date": "2025-05-07",
    "preferences": {
      "budget": "中等",
      "activities": ["历史景点", "美食", "购物"],
      "accommodation_type": "酒店"
    },
    "created_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 400 Bad Request: 请求格式无效
  - 401 Unauthorized: 未授权访问
  - 500 Internal Server Error: 服务器内部错误

### 获取旅行计划列表

- **URL**: `/api/trips`
- **方法**: `GET`
- **描述**: 获取当前用户的所有旅行计划
- **认证**: 需要JWT令牌
- **成功响应** (状态码: 200):
  ```json
  [
    {
      "id": "60f7b0b9e6b3f12345678901",
      "user_id": 1,
      "destination": "目的地1",
      "start_date": "2025-05-01",
      "end_date": "2025-05-07",
      "created_at": "2025-04-21T11:22:01+08:00"
    },
    {
      "id": "60f7b0b9e6b3f12345678902",
      "user_id": 1,
      "destination": "目的地2",
      "start_date": "2025-06-01",
      "end_date": "2025-06-07",
      "created_at": "2025-04-21T11:22:01+08:00"
    }
  ]
  ```
- **错误响应**:
  - 401 Unauthorized: 未授权访问
  - 500 Internal Server Error: 服务器内部错误

### 获取旅行计划详情

- **URL**: `/api/trips/:id`
- **方法**: `GET`
- **描述**: 获取指定ID的旅行计划详情
- **认证**: 需要JWT令牌
- **参数**:
  - `id`: 旅行计划ID (路径参数)
- **成功响应** (状态码: 200):
  ```json
  {
    "id": "60f7b0b9e6b3f12345678901",
    "user_id": 1,
    "destination": "目的地",
    "start_date": "2025-05-01",
    "end_date": "2025-05-07",
    "preferences": {
      "budget": "中等",
      "activities": ["历史景点", "美食", "购物"],
      "accommodation_type": "酒店"
    },
    "itinerary": [
      {
        "day": 1,
        "date": "2025-05-01",
        "activities": [
          {
            "time": "09:00-12:00",
            "description": "参观历史博物馆",
            "location": "市中心博物馆",
            "notes": "门票100元/人"
          },
          {
            "time": "12:30-14:00",
            "description": "午餐",
            "location": "当地特色餐厅",
            "notes": "推荐菜品：当地特色小吃"
          }
        ]
      },
      {
        "day": 2,
        "date": "2025-05-02",
        "activities": [
          {
            "time": "全天",
            "description": "海滩休闲",
            "location": "沙滩",
            "notes": "带上防晒霜"
          }
        ]
      }
    ],
    "created_at": "2025-04-21T11:22:01+08:00",
    "updated_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 401 Unauthorized: 未授权访问
  - 404 Not Found: 旅行计划不存在
  - 500 Internal Server Error: 服务器内部错误

### 更新旅行计划

- **URL**: `/api/trips/:id`
- **方法**: `PUT`
- **描述**: 更新指定ID的旅行计划
- **认证**: 需要JWT令牌
- **参数**:
  - `id`: 旅行计划ID (路径参数)
- **请求体**:
  ```json
  {
    "destination": "更新后的目的地",
    "start_date": "2025-05-03",
    "end_date": "2025-05-10",
    "preferences": {
      "budget": "高端",
      "activities": ["历史景点", "美食", "购物", "温泉"],
      "accommodation_type": "豪华酒店"
    }
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "id": "60f7b0b9e6b3f12345678901",
    "user_id": 1,
    "destination": "更新后的目的地",
    "start_date": "2025-05-03",
    "end_date": "2025-05-10",
    "preferences": {
      "budget": "高端",
      "activities": ["历史景点", "美食", "购物", "温泉"],
      "accommodation_type": "豪华酒店"
    },
    "updated_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 400 Bad Request: 请求格式无效
  - 401 Unauthorized: 未授权访问
  - 404 Not Found: 旅行计划不存在
  - 500 Internal Server Error: 服务器内部错误

### 删除旅行计划

- **URL**: `/api/trips/:id`
- **方法**: `DELETE`
- **描述**: 删除指定ID的旅行计划
- **认证**: 需要JWT令牌
- **参数**:
  - `id`: 旅行计划ID (路径参数)
- **成功响应** (状态码: 204): 无内容
- **错误响应**:
  - 401 Unauthorized: 未授权访问
  - 404 Not Found: 旅行计划不存在
  - 500 Internal Server Error: 服务器内部错误

### 生成AI旅行建议

- **URL**: `/api/trips/:id/ai-suggestions`
- **方法**: `POST`
- **描述**: 使用AI为指定的旅行计划生成建议
- **认证**: 需要JWT令牌
- **参数**:
  - `id`: 旅行计划ID (路径参数)
- **请求体**:
  ```json
  {
    "prompt": "我想要更多关于当地美食的建议",
    "preferences": {
      "budget": "经济",
      "interests": ["美食", "文化体验"]
    }
  }
  ```
- **成功响应** (状态码: 200):
  ```json
  {
    "trip_id": "60f7b0b9e6b3f12345678901",
    "suggestions": {
      "dining": [
        {
          "name": "当地特色餐厅1",
          "description": "提供正宗当地美食，价格适中",
          "location": "市中心",
          "price_range": "¥¥",
          "recommended_dishes": ["特色菜1", "特色菜2"]
        },
        {
          "name": "当地特色餐厅2",
          "description": "家庭式餐厅，提供家常菜",
          "location": "老城区",
          "price_range": "¥",
          "recommended_dishes": ["特色菜3", "特色菜4"]
        }
      ],
      "cultural_experiences": [
        {
          "name": "当地烹饪课程",
          "description": "学习制作当地特色美食",
          "location": "美食学校",
          "duration": "3小时",
          "price": "200元/人"
        },
        {
          "name": "美食街步行之旅",
          "description": "导游带领探索当地美食街",
          "location": "美食街",
          "duration": "2小时",
          "price": "150元/人"
        }
      ]
    },
    "created_at": "2025-04-21T11:22:01+08:00"
  }
  ```
- **错误响应**:
  - 400 Bad Request: 请求格式无效
  - 401 Unauthorized: 未授权访问
  - 404 Not Found: 旅行计划不存在
  - 500 Internal Server Error: 服务器内部错误

## 错误码说明

| 状态码 | 描述 | 可能原因 |
|--------|------|----------|
| 400 | 请求无效 | 请求参数格式错误或缺少必要参数 |
| 401 | 未授权 | 未提供认证令牌或令牌已过期/无效 |
| 403 | 禁止访问 | 用户无权访问请求的资源 |
| 404 | 资源不存在 | 请求的资源不存在 |
| 500 | 服务器内部错误 | 服务器处理请求时发生错误 |

## 认证说明

所有需要认证的API请求都应在HTTP头部包含以下字段：

```
Authorization: Bearer <your_jwt_token>
```

其中`<your_jwt_token>`是通过登录API获取的JWT令牌。令牌有效期为24小时，过期后需要重新登录获取新令牌。
