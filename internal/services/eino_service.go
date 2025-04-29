package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"personatrip/internal/models"
	"personatrip/pkg/einosdk"
)

// EinoService 是大模型服务的实现
type EinoService struct {
	client         *einosdk.Client
	configService  ModelConfigService
	activeConfig   *models.ModelConfig
	defaultOptions *einosdk.GenerateTextRequest
}

// NewEinoService 创建新的Eino服务实例
func NewEinoService(configService ModelConfigService) *EinoService {
	service := &EinoService{
		configService: configService,
		defaultOptions: &einosdk.GenerateTextRequest{
			MaxTokens:   8000,
			Temperature: 0.7,
		},
	}

	// 初始化时尝试加载激活的模型配置
	service.RefreshModelConfig(context.Background())

	return service
}

// NewEinoServiceWithConfig 根据指定的模型配置创建新的Eino服务实例
func NewEinoServiceWithConfig(config *models.ModelConfig) *EinoService {
	service := &EinoService{
		activeConfig: config,
		defaultOptions: &einosdk.GenerateTextRequest{
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
		},
	}

	// 创建Eino客户端
	service.client = einosdk.NewClient(
		config.ToEinoModelType(),
		config.GetEinoOptions()...,
	)

	return service
}

// TestGenerateText 测试生成文本
func (s *EinoService) TestGenerateText(ctx context.Context, prompt string) (string, error) {
	// 调用Eino API
	response, err := s.client.GenerateText(ctx, &einosdk.GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   s.defaultOptions.MaxTokens,
		Temperature: s.defaultOptions.Temperature,
	})
	if err != nil {
		return "", err
	}

	return response.Text, nil
}

// RefreshModelConfig 刷新模型配置
func (s *EinoService) RefreshModelConfig(ctx context.Context) error {
	// 从数据库获取激活的模型配置
	config, err := s.configService.GetActiveModelConfig(ctx)
	if err != nil {
		// 如果没有激活的配置，使用Mock模型
		s.client = einosdk.NewClient(einosdk.ModelTypeArk)
		return err
	}

	// 更新激活的配置
	s.activeConfig = config

	// 更新默认选项
	s.defaultOptions.MaxTokens = config.MaxTokens
	s.defaultOptions.Temperature = config.Temperature

	// 创建Eino客户端
	s.client = einosdk.NewClient(
		config.ToEinoModelType(),
		config.GetEinoOptions()...,
	)

	return nil
}

// GenerateTripPlan 根据用户请求生成旅行计划
func (s *EinoService) GenerateTripPlan(ctx context.Context, req *models.PlanRequest) (*models.TripPlan, error) {
	// 刷新模型配置，确保使用最新的配置
	s.RefreshModelConfig(ctx)

	// 构建提示词
	prompt := buildTripPlanPrompt(req)

	// 调用Eino API
	response, err := s.client.GenerateText(ctx, &einosdk.GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   8000,
		Temperature: s.defaultOptions.Temperature,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate trip plan: %w", err)
	}
	plan, err := parseTripPlanResponse(response.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip plan: %w", err)
	}

	// 填充请求中的基本信息
	plan.Destination = req.Destination
	plan.StartDate = req.StartDate.String()
	plan.EndDate = req.EndDate.String()

	return plan, nil
}

// 构建旅行计划的提示词
func buildTripPlanPrompt(req *models.PlanRequest) string {
	return fmt.Sprintf(`
你是一个专业的旅游规划助手。请为以下旅行需求创建一个详细的旅行计划:

目的地: %s
开始日期: %s
结束日期: %s
预算: %s
旅行风格: %v
住宿偏好: %v
交通偏好: %v
活动偏好: %v
饮食偏好: %v
特殊要求: %s

请提供一个包含以下内容的详细旅行计划:
1. 每天的行程安排，包括景点、活动、餐饮和住宿
2. 每个活动的大致时间安排
3. 每个活动和住宿的估计费用
4. 交通建议
5. 当地特色美食推荐
6. 天气信息和穿衣建议
7. 当地文化和习俗提示
8. 安全和健康建议
9. 必备物品清单
10. 紧急联系信息

请以JSON格式返回，格式如下:
{
  "title": "旅行计划标题",
  "destination_info": {
    "name": "目的地名称",
    "country": "国家",
    "language": "当地语言",
    "currency": "当地货币",
    "time_zone": "时区",
    "best_time_to_visit": "最佳旅游时间"
  },
  "travel_info": {
    "visa_required": true/false,
    "visa_tips": "签证信息",
    "passport_validity": "护照有效期要求",
    "vaccination_required": ["疫苗1", "疫苗2"],
    "local_customs": "当地习俗简介",
    "etiquette_tips": "礼仪提示",
    "safety_tips": "安全提示",
    "health_tips": "健康建议",
    "electrical_socket_type": "电源插座类型",
    "internet_availability": "网络可用性说明",
    "language_phrases": [
      {"phrase": "你好", "pronunciation": "Ni Hao", "meaning": "Hello"}
    ]
  },
  "weather_forecast": {
    "climate_overview": "季节性气候概况",
    "daily_forecast": [
      {
        "date": "YYYY-MM-DD",
        "temperature": {
          "min": 最低温度,
          "max": 最高温度,
          "unit": "摄氏/华氏"
        },
        "conditions": "天气状况",
        "precipitation_chance": 降水几率,
        "clothing_suggestions": ["穿衣建议1", "穿衣建议2"]
      }
    ]
  },
  "packing_list": {
    "essentials": ["必备物品1", "必备物品2"],
    "clothing": ["衣物1", "衣物2"],
    "toiletries": ["洗漱用品1", "洗漱用品2"],
    "electronics": ["电子设备1", "电子设备2"],
    "documents": ["文档1", "文档2"],
    "other": ["其他物品1", "其他物品2"]
  },
  "emergency_contacts": {
    "local_emergency": "当地紧急电话",
    "police": "警察电话",
    "ambulance": "救护车电话",
    "fire": "消防电话",
    "embassy": "使馆信息",
    "hospitals": [
      {
        "name": "医院名称",
        "address": "地址",
        "phone": "电话",
        "has_english_speaking_staff": true/false
      }
    ]
  },
  "days": [
    {
      "day": 1,
      "date": "YYYY-MM-DD",
      "weather": {
        "temperature": {
          "morning": 早晨温度,
          "day": 白天温度,
          "evening": 傍晚温度,
          "unit": "摄氏/华氏"
        },
        "conditions": "天气状况",
        "clothing_suggestion": "穿衣建议"
      },
      "activities": [
        {
          "name": "活动名称",
          "type": "活动类型",
          "location": {
            "name": "地点名称",
            "address": "地址",
            "city": "城市",
            "country": "国家",
            "coordinates": {
              "latitude": 纬度,
              "longitude": 经度
            }
          },
          "start_time": "HH:MM",
          "end_time": "HH:MM",
          "description": "活动描述",
          "cost": 费用数值,
          "booking_required": true/false,
          "booking_tips": "预订提示",
          "crowd_level": "人群水平预期",
          "suitable_weather": "适合的天气条件",
          "indoor_outdoor": "室内/室外",
          "accessibility": "无障碍设施情况",
          "rating": 评分,
          "photos": ["照片URL1", "照片URL2"],
          "tips": ["小贴士1", "小贴士2"]
        }
      ],
      "meals": [
        {
          "type": "餐食类型",
          "venue": "餐厅名称",
          "cuisine": "菜系",
          "description": "描述",
          "specialties": ["特色菜1", "特色菜2"],
          "dietary_options": ["素食", "无麸质"],
          "address": "地址",
          "booking_required": true/false,
          "cost": 费用数值,
          "tips": "用餐提示"
        }
      ],
      "accommodation": {
        "name": "住宿名称",
        "type": "住宿类型",
        "address": "地址",
        "description": "描述",
        "amenities": ["设施1", "设施2"],
        "check_in": "入住时间",
        "check_out": "退房时间",
        "cost": 费用数值,
        "booking_reference": "预订参考信息",
        "contact": "联系方式",
        "nearest_landmarks": ["地标1", "地标2"],
        "transportation_options": ["交通选项1", "交通选项2"]
      },
      "transportation": [
        {
          "type": "交通类型",
          "from": "出发地",
          "to": "目的地",
          "departure_time": "出发时间",
          "arrival_time": "到达时间",
          "cost": 费用数值,
          "booking_reference": "预订参考信息",
          "notes": "交通备注"
        }
      ],
      "tips": ["当天提示1", "当天提示2"]
    }
  ],
  "budget": {
    "currency": "货币",
    "exchange_rate": "汇率",
    "total_estimate": 总预算,
    "accommodation": 住宿预算,
    "transportation": 交通预算,
    "food": 餐饮预算,
    "activities": 活动预算,
    "shopping": 购物预算,
    "other": 其他预算,
    "daily_breakdown": [
      {
        "day": 1,
        "date": "YYYY-MM-DD",
        "total": 当天总花费,
        "details": {
          "accommodation": 住宿费用,
          "transportation": 交通费用,
          "food": 餐饮费用,
          "activities": 活动费用,
          "other": 其他费用
        }
      }
    ],
    "payment_tips": {
      "credit_cards_accepted": true/false,
      "atm_availability": "ATM可用性",
      "tipping_culture": "小费文化",
      "recommended_payment_methods": ["建议支付方式1", "建议支付方式2"]
    }
  },
  "local_attractions": [
    {
      "name": "景点名称",
      "category": "景点类别",
      "description": "描述",
      "must_see": true/false,
      "address": "地址",
      "opening_hours": "开放时间",
      "cost": 费用数值,
      "time_required": "建议游览时间",
      "best_time_to_visit": "最佳游览时间",
      "tips": ["小贴士1", "小贴士2"]
    }
  ],
  "local_cuisine": [
    {
      "name": "美食名称",
      "description": "描述",
      "must_try": true/false,
      "where_to_find": ["地点1", "地点2"],
      "price_range": "价格范围",
      "photos": ["照片URL1", "照片URL2"]
    }
  ],
  "shopping": {
    "recommended_items": ["推荐购买物品1", "推荐购买物品2"],
    "markets_and_malls": [
      {
        "name": "商场/市场名称",
        "type": "类型",
        "address": "地址",
        "specialty": "特色",
        "opening_hours": "营业时间"
      }
    ],
    "souvenirs": ["纪念品1", "纪念品2"]
  },
  "cultural_events": [
    {
      "name": "文化活动名称",
      "date": "日期",
      "description": "描述",
      "location": "地点",
      "cost": 费用数值,
      "tips": "参与提示"
    }
  ],
  "practical_information": {
    "local_transportation": {
      "options": ["选项1", "选项2"],
      "recommended": "推荐方式",
      "cost": "费用信息",
      "passes": "交通通行证信息",
      "apps": ["推荐应用1", "推荐应用2"]
    },
    "communication": {
      "local_sim": "当地SIM卡信息",
      "wifi_availability": "WiFi可用性",
      "useful_apps": ["有用的应用1", "有用的应用2"]
    }
  },
  "notes": "额外注意事项",
  "suggested_modifications": "根据天气或其他因素可能需要的计划调整建议"
}
`,
		req.Destination,
		req.StartDate,
		req.EndDate,
		req.Budget,
		req.TravelStyle,
		req.Accommodation,
		req.Transportation,
		req.Activities,
		req.FoodPreferences,
		req.SpecialRequests,
	)
}

// 解析大模型返回的旅行计划
func parseTripPlanResponse(response string) (*models.TripPlan, error) {
	var plan models.TripPlan

	// 尝试直接解析JSON
	err := json.Unmarshal([]byte(response), &plan)
	if err != nil {
		// 如果直接解析失败，尝试提取JSON部分
		jsonStart := strings.Index(response, "{")
		if jsonStart >= 0 {
			// 寻找匹配的大括号结束位置
			braceCount := 0
			jsonEnd := -1

			for i := jsonStart; i < len(response); i++ {
				if response[i] == '{' {
					braceCount++
				} else if response[i] == '}' {
					braceCount--
					if braceCount == 0 {
						jsonEnd = i + 1
						break
					}
				}
			}

			if jsonEnd > jsonStart {
				jsonStr := response[jsonStart:jsonEnd]
				// 清理JSON字符串，替换中文标点等
				jsonStr = cleanJSONString(jsonStr)

				// 处理JSON中destination字段映射到DestinationInfo
				var rawJSON map[string]interface{}
				if err := json.Unmarshal([]byte(jsonStr), &rawJSON); err == nil {
					if dest, ok := rawJSON["destination"].(map[string]interface{}); ok {
						// 将destination对象内容复制到destination_info
						rawJSON["destination_info"] = dest

						// 保存原始destination字符串值
						if name, ok := dest["name"].(string); ok {
							rawJSON["destination"] = name
						}

						// 重新序列化修正后的JSON
						if newJSON, err := json.Marshal(rawJSON); err == nil {
							jsonStr = string(newJSON)
						}
					}
				}

				err = json.Unmarshal([]byte(jsonStr), &plan)
				if err == nil {
					return &plan, nil
				}
			}
		}

		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// 检查是否需要处理destination字段
	var rawJSON map[string]interface{}
	if json.Unmarshal([]byte(response), &rawJSON) == nil {
		if dest, ok := rawJSON["destination"].(map[string]interface{}); ok {
			// 如果destination是对象而不是字符串，则需要转换
			if name, ok := dest["name"].(string); ok {
				plan.Destination = name
			}

			// 将destination对象映射到DestinationInfo
			destJSON, _ := json.Marshal(dest)
			json.Unmarshal(destJSON, &plan.DestinationInfo)
		}
	}

	return &plan, nil
}

// cleanJSONString 尝试修复常见的JSON格式问题
func cleanJSONString(jsonStr string) string {
	// 替换非标准引号 (使用Unicode码点)
	jsonStr = strings.ReplaceAll(jsonStr, "\u201c", "\"") // 左双引号
	jsonStr = strings.ReplaceAll(jsonStr, "\u201d", "\"") // 右双引号
	jsonStr = strings.ReplaceAll(jsonStr, "\u2018", "\"") // 左单引号
	jsonStr = strings.ReplaceAll(jsonStr, "\u2019", "\"") // 右单引号

	// 替换中文标点符号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff0c", ",")  // 中文逗号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff1a", ":")  // 中文冒号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff1b", ";")  // 中文分号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3001", ",")  // 中文顿号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff01", "!")  // 中文感叹号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff1f", "?")  // 中文问号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3002", ".")  // 中文句号
	jsonStr = strings.ReplaceAll(jsonStr, "\u300c", "\"") // 中文左引号
	jsonStr = strings.ReplaceAll(jsonStr, "\u300d", "\"") // 中文右引号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff08", "(")  // 中文左括号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff09", ")")  // 中文右括号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3010", "[")  // 中文左方括号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3011", "]")  // 中文右方括号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3008", "<")  // 中文左尖括号
	jsonStr = strings.ReplaceAll(jsonStr, "\u3009", ">")  // 中文右尖括号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff0b", "+")  // 中文加号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff0d", "-")  // 中文减号
	jsonStr = strings.ReplaceAll(jsonStr, "\uff0f", "/")  // 中文斜杠
	jsonStr = strings.ReplaceAll(jsonStr, "\uff3c", "\\") // 中文反斜杠

	// 移除特殊控制字符
	var cleanBytes []byte
	for _, c := range []byte(jsonStr) {
		if c >= 32 && c <= 126 || c >= 128 {
			cleanBytes = append(cleanBytes, c)
		}
	}

	return string(cleanBytes)
}

// GenerateDestinationRecommendations 根据用户偏好生成目的地推荐
func (s *EinoService) GenerateDestinationRecommendations(ctx context.Context, preferences *models.UserPreferences) ([]string, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`
基于以下用户偏好，推荐5个最适合的旅游目的地:

旅行风格: %v
预算: %s
住宿偏好: %v
交通偏好: %v
活动偏好: %v
饮食偏好: %v

请以JSON数组格式返回5个推荐目的地，每个目的地包含名称和简短理由。
`,
		preferences.TravelStyle,
		preferences.Budget,
		preferences.Accommodation,
		preferences.Transportation,
		preferences.Activities,
		preferences.FoodPreferences,
	)

	// 调用Eino API
	response, err := s.client.GenerateText(ctx, &einosdk.GenerateTextRequest{
		Model:       "eino-large", // 使用适当的模型名称
		Prompt:      prompt,
		MaxTokens:   2000,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// 解析响应
	var recommendations []string
	err = json.Unmarshal([]byte(response.Text), &recommendations)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recommendations: %w", err)
	}

	return recommendations, nil
}
