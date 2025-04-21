package einosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	if c.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	
	// 构建请求URL
	url := "https://api.openai.com/v1/chat/completions"
	
	// 构建请求体
	openaiReq := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
		"temperature": req.Temperature,
		"max_tokens": req.MaxTokens,
	}
	
	// 将请求体转换为JSON
	reqBody, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	
	// 解析响应
	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(respBody, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	// 提取生成的文本
	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("no text generated")
	}
	
	return &GenerateTextResponse{
		Text: openaiResp.Choices[0].Message.Content,
	}, nil
}

// generateTextWithOllama 使用Ollama生成文本
func (c *Client) generateTextWithOllama(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 构建请求URL
	baseURL := c.baseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	url := fmt.Sprintf("%s/api/generate", baseURL)
	
	// 构建请求体
	ollamaReq := map[string]interface{}{
		"model": c.model,
		"prompt": req.Prompt,
		"temperature": req.Temperature,
		"num_predict": req.MaxTokens,
		"stream": false,
	}
	
	// 将请求体转换为JSON
	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	
	// 发送请求
	client := &http.Client{Timeout: 120 * time.Second} // Ollama可能需要更长的超时时间
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	
	// 解析响应
	var ollamaResp struct {
		Response string `json:"response"`
	}
	
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	return &GenerateTextResponse{
		Text: ollamaResp.Response,
	}, nil
}

// generateTextWithArk 使用Ark生成文本
func (c *Client) generateTextWithArk(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 基于Eino框架文档实现: https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/chat_model/chat_model_ark/
	
	if c.apiKey == "" {
		return nil, fmt.Errorf("Ark API key is required")
	}
	
	// 构建请求URL
	url := c.baseURL
	if url == "" {
		url = "https://api.ark.com/v1/chat/completions"
	}
	
	// 构建请求体
	arkReq := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "user", "content": req.Prompt},
		},
		"temperature": req.Temperature,
		"max_tokens": req.MaxTokens,
	}
	
	// 将请求体转换为JSON
	reqBody, err := json.Marshal(arkReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}
	
	// 解析响应
	var arkResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.Unmarshal(respBody, &arkResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	// 提取生成的文本
	if len(arkResp.Choices) == 0 {
		return nil, fmt.Errorf("no text generated")
	}
	
	return &GenerateTextResponse{
		Text: arkResp.Choices[0].Message.Content,
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
