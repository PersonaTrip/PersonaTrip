package services

import (
	"context"

	"personatrip/internal/models"
	"personatrip/pkg/einosdk"
)

// NewEinoServiceWithConfig 创建使用特定配置的Eino服务实例
func NewEinoServiceWithConfig(config *models.ModelConfig) *EinoService {
	// 创建客户端
	client := einosdk.NewClient(
		config.ToEinoModelType(),
		config.GetEinoOptions()...,
	)

	// 创建服务
	return &EinoService{
		client: client,
		defaultOptions: &einosdk.GenerateTextRequest{
			MaxTokens:   config.MaxTokens,
			Temperature: config.Temperature,
		},
	}
}

// TestGenerateText 测试生成文本
func (s *EinoService) TestGenerateText(ctx context.Context, prompt string) (string, error) {
	// 调用Eino API
	response, err := s.client.GenerateText(ctx, &einosdk.GenerateTextRequest{
		Prompt:      prompt,
		MaxTokens:   s.defaultOptions.MaxTokens,
		Temperature: s.defaultOptions.Temperature,
	})
	if err != nil {
		return "", err
	}

	return response.Text, nil
}
