package einosdk

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
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

// getDefaultMaxToken 获取默认的上下文长度
func getDefaultMaxToken() int {
	return 8000
}

// GenerateTextRequest 是生成文本的请求参数
type GenerateTextRequest struct {
	Model       string          `json:"model"`
	Prompt      string          `json:"prompt"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float32         `json:"temperature"`
	Tools       []tool.BaseTool `json:"tools"`
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
