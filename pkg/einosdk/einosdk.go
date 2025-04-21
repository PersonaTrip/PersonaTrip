package einosdk

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// ModelType 表示支持的大模型类型
type ModelType string

const (
	ModelTypeOpenAI ModelType = "openai"
	ModelTypeOllama ModelType = "ollama"
	ModelTypeArk    ModelType = "ark"
	ModelTypeMock   ModelType = "mock" // 用于测试
)

// Client 是Eino API的客户端
type Client struct {
	modelType ModelType
	apiKey    string
	baseURL   string
	model     string
}

// ClientOption 是Client的配置选项
type ClientOption func(*Client)

// WithAPIKey 设置API密钥
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithModel 设置模型名称
func WithModel(model string) ClientOption {
	return func(c *Client) {
		c.model = model
	}
}

// NewClient 创建一个新的Eino客户端
func NewClient(modelType ModelType, opts ...ClientOption) *Client {
	c := &Client{
		modelType: modelType,
		model:     getDefaultModel(modelType),
	}
	
	// 应用选项
	for _, opt := range opts {
		opt(c)
	}
	
	// 如果没有设置API密钥，尝试从环境变量获取
	if c.apiKey == "" {
		c.apiKey = getAPIKeyFromEnv(modelType)
	}
	
	// 如果没有设置基础URL，使用默认值
	if c.baseURL == "" {
		c.baseURL = getDefaultBaseURL(modelType)
	}
	
	return c
}

// getDefaultModel 根据模型类型返回默认模型名称
func getDefaultModel(modelType ModelType) string {
	switch modelType {
	case ModelTypeOpenAI:
		return "gpt-3.5-turbo"
	case ModelTypeOllama:
		return "llama2"
	case ModelTypeArk:
		return "ark-large"
	case ModelTypeMock:
		return "mock-model"
	default:
		return "gpt-3.5-turbo"
	}
}

// getAPIKeyFromEnv 从环境变量获取API密钥
func getAPIKeyFromEnv(modelType ModelType) string {
	switch modelType {
	case ModelTypeOpenAI:
		return os.Getenv("OPENAI_API_KEY")
	case ModelTypeArk:
		return os.Getenv("ARK_API_KEY")
	default:
		return ""
	}
}

// getDefaultBaseURL 获取默认的基础URL
func getDefaultBaseURL(modelType ModelType) string {
	switch modelType {
	case ModelTypeOllama:
		return "http://localhost:11434"
	case ModelTypeArk:
		return "https://api.ark.com/v1"
	default:
		return ""
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
	// 根据不同的模型类型调用不同的API
	switch c.modelType {
	case ModelTypeOpenAI:
		return c.generateTextWithOpenAI(ctx, req)
	case ModelTypeOllama:
		return c.generateTextWithOllama(ctx, req)
	case ModelTypeArk:
		return c.generateTextWithArk(ctx, req)
	case ModelTypeMock:
		return c.generateTextMock(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported model type: %s", c.modelType)
	}
}

// generateTextWithOpenAI 使用OpenAI API生成文本
func (c *Client) generateTextWithOpenAI(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 实际项目中应该调用OpenAI API
	// 这里仅作为示例
	return &GenerateTextResponse{
		Text: "这是OpenAI模型生成的回复",
	}, nil
}

// generateTextWithOllama 使用Ollama生成文本
func (c *Client) generateTextWithOllama(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 实际项目中应该调用Ollama API
	// 这里仅作为示例
	return &GenerateTextResponse{
		Text: "这是Ollama模型生成的回复",
	}, nil
}

// generateTextWithArk 使用Ark生成文本
func (c *Client) generateTextWithArk(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 实际项目中应该调用Ark API
	// 这里仅作为示例
	return &GenerateTextResponse{
		Text: "这是Ark模型生成的回复",
	}, nil
}

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
