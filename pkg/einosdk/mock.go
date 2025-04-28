package einosdk

import (
	"context"
	"encoding/json"
	"fmt"
)

// generateTextMock 生成模拟文本响应
func (c *Client) generateTextMock(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 这是一个模拟实现，用于测试
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
