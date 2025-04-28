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
		"model":       c.model,
		"prompt":      req.Prompt,
		"temperature": req.Temperature,
		"num_predict": req.MaxTokens,
		"stream":      false,
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
