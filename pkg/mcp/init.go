package mcp

import (
	"context"
	"fmt"
	"personatrip/internal/config"
	"personatrip/internal/utils/logger"

	"personatrip/pkg/mcp/providers/amap"
)

// InitOptions 初始化选项
type InitOptions struct {
	// AMapAPIKey 高德地图API密钥，如果为空则不初始化高德地图提供者
	AMapAPIKey string
	// 其他提供者的配置可以在这里添加
}

// NewInitOptionsFromConfig 从应用配置创建初始化选项
func NewInitOptionsFromConfig(cfg *config.Config) *InitOptions {
	if cfg == nil || cfg.MCPConfig == nil {
		return &InitOptions{}
	}

	return &InitOptions{
		AMapAPIKey: cfg.MCPConfig.AMapAPIKey,
	}
}

// InitMCPClient 初始化MCP客户端及所有配置的提供者
func InitMCPClient(ctx context.Context, opts *InitOptions) (*Client, error) {
	if opts == nil {
		// 如果未提供选项，尝试从配置文件加载
		cfg, err := config.Load()
		if err != nil {
			logger.Warnf("加载配置文件失败，使用默认选项: %v", err)
			opts = &InitOptions{}
		} else {
			opts = NewInitOptionsFromConfig(cfg)
		}
	}

	client := NewClient()

	// 初始化高德地图提供者
	if opts.AMapAPIKey != "" {
		amapProvider := amap.NewProvider(opts.AMapAPIKey)
		client.AddProvider(ProviderAMap, amapProvider)
	}

	// 这里可以添加其他提供者的初始化

	// 初始化所有提供者
	if err := client.Initialize(ctx); err != nil {
		return client, fmt.Errorf("初始化MCP客户端失败: %w", err)
	}

	logger.Infof("MCP客户端初始化成功")

	// 打印所有已加载的工具
	PrintLoadedTools(ctx, client)

	return client, nil
}

// PrintLoadedTools 打印所有已加载的工具
func PrintLoadedTools(ctx context.Context, client *Client) {
	tools, err := client.GetAllTools(ctx)
	if err != nil {
		logger.Infof("获取工具列表失败: %v\n", err)
		return
	}

	fmt.Println("已加载的MCP工具:")
	for provider, providerTools := range tools {
		fmt.Println("提供者: %s\n", provider)
		for _, tool := range providerTools {
			toolInfo, err := tool.Info(ctx)
			if err != nil {
				logger.Errorf("tool info failure: %v", err)
			}
			fmt.Printf("  - %s: %s\n", toolInfo.Name, toolInfo.Desc)
		}
	}
}
