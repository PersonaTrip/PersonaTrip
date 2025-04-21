package einosdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// Client 是Eino API的客户端
type Client struct {
	apiKey string
}

// NewClient 创建一个新的Eino客户端
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

// GenerateTextRequest 是生成文本的请求参数
type GenerateTextRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

// GenerateTextResponse 是生成文本的响应
type GenerateTextResponse struct {
	Text string `json:"text"`
}

// GenerateText 调用Eino API生成文本
func (c *Client) GenerateText(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 这是一个模拟实现，实际项目中应该调用真正的API
	// 在这里，我们只是返回一个简单的旅行计划JSON作为示例
	
	// 根据提示词中的目的地生成不同的响应
	var destination string
	if req.Prompt != "" {
		// 尝试从提示词中提取目的地信息
		var destMap map[string]string
		if err := json.Unmarshal([]byte(`{"destination":"东京"}`), &destMap); err == nil {
			destination = destMap["destination"]
		}
		
		if destination == "" {
			destination = "东京" // 默认目的地
		}
	}
	
	// 生成示例旅行计划
	planJSON := fmt.Sprintf(`{
		"title": "%s旅行计划",
		"days": [
			{
				"day": 1,
				"date": "2025-05-01",
				"activities": [
					{
						"name": "参观%s塔",
						"type": "景点",
						"location": {
							"name": "%s塔",
							"address": "%s中心区",
							"city": "%s",
							"country": "日本"
						},
						"start_time": "10:00",
						"end_time": "12:00",
						"description": "欣赏%s全景",
						"cost": 2000
					}
				],
				"meals": [
					{
						"type": "午餐",
						"venue": "%s特色餐厅",
						"description": "品尝当地美食",
						"cost": 1500
					}
				],
				"accommodation": {
					"name": "%s中心酒店",
					"type": "酒店",
					"description": "位于市中心的舒适酒店",
					"cost": 8000
				}
			}
		],
		"budget": {
			"currency": "JPY",
			"total_estimate": 50000,
			"accommodation": 24000,
			"transportation": 10000,
			"food": 10000,
			"activities": 5000,
			"other": 1000
		},
		"notes": "这是一个AI生成的%s旅行计划示例"
	}`, destination, destination, destination, destination, destination, destination, destination, destination, destination)
	
	return &GenerateTextResponse{
		Text: planJSON,
	}, nil
}
