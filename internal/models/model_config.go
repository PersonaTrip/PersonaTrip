package models

import (
	"time"

	"personatrip/pkg/einosdk"
)

// ModelConfig 表示大模型配置
type ModelConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	ModelType   string    `json:"model_type" gorm:"size:50;not null"` // openai, ollama, ark, mock
	ModelName   string    `json:"model_name" gorm:"size:100;not null"`
	ApiKey      string    `json:"api_key,omitempty" gorm:"size:255"`
	BaseUrl     string    `json:"base_url,omitempty" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"default:false"`
	Temperature float64   `json:"temperature" gorm:"default:0.7"`
	MaxTokens   int       `json:"max_tokens" gorm:"default:2000"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToEinoModelType 将字符串类型转换为einosdk.ModelType
func (m *ModelConfig) ToEinoModelType() einosdk.ModelType {
	switch m.ModelType {
	case "openai":
		return einosdk.ModelTypeOpenAI
	case "ollama":
		return einosdk.ModelTypeOllama
	case "ark":
		return einosdk.ModelTypeArk
	case "mock":
		return einosdk.ModelTypeMock
	default:
		return einosdk.ModelTypeMock
	}
}

// GetEinoOptions 获取Eino客户端选项
func (m *ModelConfig) GetEinoOptions() []einosdk.ClientOption {
	options := []einosdk.ClientOption{
		einosdk.WithModel(m.ModelName),
	}

	if m.ApiKey != "" {
		options = append(options, einosdk.WithAPIKey(m.ApiKey))
	}

	if m.BaseUrl != "" {
		options = append(options, einosdk.WithBaseURL(m.BaseUrl))
	}

	return options
}

// ModelConfigResponse 是模型配置的响应格式
type ModelConfigResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	ModelType   string  `json:"model_type"`
	ModelName   string  `json:"model_name"`
	BaseUrl     string  `json:"base_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ToResponse 将ModelConfig转换为ModelConfigResponse
func (m *ModelConfig) ToResponse() ModelConfigResponse {
	return ModelConfigResponse{
		ID:          m.ID,
		Name:        m.Name,
		ModelType:   m.ModelType,
		ModelName:   m.ModelName,
		BaseUrl:     m.BaseUrl,
		IsActive:    m.IsActive,
		Temperature: m.Temperature,
		MaxTokens:   m.MaxTokens,
	}
}

// ModelConfigCreateRequest 是创建模型配置的请求格式
type ModelConfigCreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	ModelType   string  `json:"model_type" binding:"required,oneof=openai ollama ark mock"`
	ModelName   string  `json:"model_name" binding:"required"`
	ApiKey      string  `json:"api_key"`
	BaseUrl     string  `json:"base_url"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ModelConfigUpdateRequest 是更新模型配置的请求格式
type ModelConfigUpdateRequest struct {
	Name        string  `json:"name"`
	ModelType   string  `json:"model_type" binding:"omitempty,oneof=openai ollama ark mock"`
	ModelName   string  `json:"model_name"`
	ApiKey      string  `json:"api_key"`
	BaseUrl     string  `json:"base_url"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ModelConfigListItem 模型配置列表项
type ModelConfigListItem struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	ModelType   string  `json:"model_type"`
	ModelName   string  `json:"model_name"`
	BaseUrl     string  `json:"base_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// ToListItem 转换为列表项
func (m *ModelConfig) ToListItem() ModelConfigListItem {
	return ModelConfigListItem{
		ID:          m.ID,
		Name:        m.Name,
		ModelType:   m.ModelType,
		ModelName:   m.ModelName,
		BaseUrl:     m.BaseUrl,
		IsActive:    m.IsActive,
		Temperature: m.Temperature,
		MaxTokens:   m.MaxTokens,
	}
}

// CreateModelConfigRequest 创建模型配置请求
type CreateModelConfigRequest struct {
	Name        string  `json:"name" binding:"required"`
	ModelType   string  `json:"model_type" binding:"required,oneof=openai ollama ark mock"`
	ModelName   string  `json:"model_name" binding:"required"`
	ApiKey      string  `json:"api_key"`
	BaseUrl     string  `json:"base_url"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// UpdateModelConfigRequest 更新模型配置请求
type UpdateModelConfigRequest struct {
	Name        string  `json:"name"`
	ModelType   string  `json:"model_type" binding:"omitempty,oneof=openai ollama ark mock"`
	ModelName   string  `json:"model_name"`
	ApiKey      string  `json:"api_key"`
	BaseUrl     string  `json:"base_url"`
	IsActive    bool    `json:"is_active"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}
