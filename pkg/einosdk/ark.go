package einosdk

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

// generateTextWithArk 使用Ark生成文本
func (c *Client) generateTextWithArk(ctx context.Context, req *GenerateTextRequest) (*GenerateTextResponse, error) {
	// 检查API密钥
	apiKey := c.apiKey
	if apiKey == "" {
		return nil, fmt.Errorf("ARK API密钥是必需的")
	}

	fmt.Println("正在准备调用ARK模型:", c.model)

	// 初始化模型
	fmt.Println(apiKey, c.model, c.baseURL)
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey:  apiKey,
		Model:   c.model,
		BaseURL: c.baseURL,
	})
	if err != nil {
		fmt.Printf("初始化ARK模型失败: %v\n", err)
		fmt.Println("切换到模拟模式...")
		return c.generateTextMock(ctx, req)
	}

	// 准备消息
	messages := []*schema.Message{
		schema.SystemMessage("你是一个旅行规划助手，帮助用户规划旅行计划，请返回有效的JSON格式数据，不要添加任何代码块反引号(```)或其他标记。"),
		schema.UserMessage(req.Prompt),
	}

	fmt.Println("连接ARK模型流式API...")

	// 获取流式回复
	reader, err := model.Stream(ctx, messages)
	if err != nil {
		fmt.Printf("连接流式API失败: %v\n", err)
		fmt.Println("切换到模拟模式...")
		return c.generateTextMock(ctx, req)
	}
	defer reader.Close() // 确保关闭流

	// 处理流式内容
	fmt.Println("\n--- ARK模型开始生成回复 ---")
	var fullResponse strings.Builder

	for {
		chunk, err := reader.Recv()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("接收数据时出错: %v\n", err)
			}
			break
		}

		// 打印并累积响应内容
		if chunk.Content != "" {
			fmt.Print(chunk.Content)
			fullResponse.WriteString(chunk.Content)
		}
	}

	fmt.Println("\n--- ARK模型回复结束 ---")

	// 如果没有收到任何内容，返回错误
	if fullResponse.Len() == 0 {
		fmt.Println("未收到任何内容，切换到模拟模式...")
		return c.generateTextMock(ctx, req)
	}

	return &GenerateTextResponse{
		Text: fullResponse.String(),
	}, nil
}
