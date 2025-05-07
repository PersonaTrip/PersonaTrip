package einosdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"io"
	"personatrip/internal/utils/logger"
	"strings"
)

// generateTextWithArk 使用Ark生成文本
func (c *Client) generateTextWithArk(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 检查API密钥
	apiKey := c.apiKey
	if apiKey == "" {
		return nil, fmt.Errorf("ARK API密钥是必需的")
	}
	logger.Info("正在准备调用ARK模型:", c.model)
	// 初始化模型
	logger.Debugf("ARK配置: APIKey=%s, Model=%s, BaseURL=%s", apiKey, c.model, c.baseURL)
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:    apiKey,
		Model:     c.model,
		BaseURL:   c.baseURL,
		MaxTokens: &req.MaxTokens,
	})
	if err != nil {
		logger.Errorf("初始化ARK模型失败: %v", err)
		logger.Info("切换到模拟模式...")
		return c.generateTextMock(ctx, req)
	}

	// 准备消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个旅行规划助手，帮助用户规划旅行计划，请返回有效的JSON格式数据，不要添加任何代码块反引号(```)或其他标记。"),
		schema.UserMessage(req.Prompt),
	}

	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		Model:            model,
		ToolCallingModel: model,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: req.Tools,
		},
		MaxStep:               80,
		StreamToolCallChecker: ARKToolCallChecker,
	})
	if err != nil {
		logger.Errorf("react new agent failure: %v", err)
		return nil, err
	}

	reader, err := ragent.Stream(ctx, messages)
	if err != nil {
		logger.Errorf("ragent generate failure: %v", err)
		return nil, err
	}

	// 处理流式内容
	logger.Info("--- ARK模型开始生成回复 ---")
	var fullResponse strings.Builder
	for {
		chunk, err := reader.Recv()
		if err != nil {
			if err != io.EOF {
				logger.Errorf("接收数据时出错: %v", err)
			}
			break
		}
		if chunk.Extra["ark-reasoning-content"] != nil {
			fmt.Printf("%v", chunk.Extra["ark-reasoning-content"])
		} else if chunk.Content != "" {
			fullResponse.WriteString(chunk.Content)
		}
	}

	logger.Info("--- ARK模型回复结束 ---")

	// 如果没有收到任何内容，返回错误
	if fullResponse.Len() == 0 {
		logger.Info("未收到任何内容，切换到模拟模式...")
		return c.generateTextMock(ctx, req)
	}
	return &GenerateTextResponse{
		Text: fullResponse.String(),
	}, nil
}

func ARKToolCallChecker(ctx context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()
	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// 流结束，未检测到工具调用
				break
			}
			return false, err
		}
		if msg.Content != "" {
			fmt.Printf(msg.Content)
		} else if msg.Extra["ark-reasoning-content"] != "" {
			fmt.Printf("%v", msg.Extra["ark-reasoning-content"])
		}
		if len(msg.ToolCalls) > 0 {
			return true, nil
		}
	}

	return false, nil
}
