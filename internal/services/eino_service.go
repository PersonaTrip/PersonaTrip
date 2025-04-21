package services

import (
	"context"
	"encoding/json"
	"fmt"

	"personatrip/internal/models"
	"personatrip/pkg/einosdk"
)

// EinoService 处理与Eino大模型的交互
type EinoService struct {
	client *einosdk.Client
}

// NewEinoService 创建新的Eino服务实例
func NewEinoService(apiKey string) *EinoService {
	// 默认使用Mock模型，方便测试
	client := einosdk.NewClient(einosdk.ModelTypeMock, einosdk.WithAPIKey(apiKey))
	return &EinoService{
		client: client,
	}
}

// NewEinoServiceWithModel 创建指定模型类型的Eino服务实例
func NewEinoServiceWithModel(modelType einosdk.ModelType, options ...einosdk.ClientOption) *EinoService {
	client := einosdk.NewClient(modelType, options...)
	return &EinoService{
		client: client,
	}
}

// GenerateTripPlan 根据用户请求生成旅行计划
func (s *EinoService) GenerateTripPlan(ctx context.Context, req *models.PlanRequest) (*models.TripPlan, error) {
	// 构建提示词
	prompt := buildTripPlanPrompt(req)

	// 调用Eino API
	response, err := s.client.GenerateText(ctx, &einosdk.GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   2000,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate trip plan: %w", err)
	}

	// 解析响应
	plan, err := parseTripPlanResponse(response.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trip plan: %w", err)
	}

	// 填充请求中的基本信息
	plan.Destination = req.Destination
	plan.StartDate = req.StartDate
	plan.EndDate = req.EndDate

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

请以JSON格式返回，格式如下:
{
  "title": "旅行计划标题",
  "days": [
    {
      "day": 1,
      "date": "YYYY-MM-DD",
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
          "start_time": "HH:MM",
          "end_time": "HH:MM",
          "description": "活动描述",
          "cost": 费用数值
        }
      ],
      "meals": [
        {
          "type": "餐食类型",
          "venue": "餐厅名称",
          "description": "描述",
          "cost": 费用数值
        }
      ],
      "accommodation": {
        "name": "住宿名称",
        "type": "住宿类型",
        "description": "描述",
        "cost": 费用数值
      }
    }
  ],
  "budget": {
    "currency": "货币",
    "total_estimate": 总预算,
    "accommodation": 住宿预算,
    "transportation": 交通预算,
    "food": 餐饮预算,
    "activities": 活动预算,
    "other": 其他预算
  },
  "notes": "额外注意事项"
}
`,
		req.Destination,
		req.StartDate.Format("2006-01-02"),
		req.EndDate.Format("2006-01-02"),
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
		// 这里可以添加更复杂的逻辑来处理非标准JSON响应
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	return &plan, nil
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
		Model:    "eino-large",  // 使用适当的模型名称
		Prompt:   prompt,
		MaxTokens: 1000,
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
