package einosdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
		"max_tokens":  req.MaxTokens,
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
