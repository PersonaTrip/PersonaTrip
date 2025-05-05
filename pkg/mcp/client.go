package mcp

import (
	"context"
	"fmt"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

// Client 是MCPClient接口的实现
type Client struct {
	providers map[string]MCPProvider
	mu        sync.RWMutex
}

// NewClient 创建一个新的MCP客户端
func NewClient() *Client {
	return &Client{
		providers: make(map[string]MCPProvider),
	}
}

// AddProvider 添加一个MCP工具提供者
func (c *Client) AddProvider(name string, provider MCPProvider) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.providers[name] = provider
}

// Initialize 初始化所有提供者
func (c *Client) Initialize(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for name, provider := range c.providers {
		if err := provider.Initialize(ctx); err != nil {
			return fmt.Errorf("初始化提供者 %s 失败: %w", name, err)
		}
	}

	return nil
}

// GetTools 获取指定提供者的所有工具
func (c *Client) GetTools(ctx context.Context, providerName string) ([]mcp.Tool, error) {
	c.mu.RLock()
	provider, ok := c.providers[providerName]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("提供者 %s 不存在", providerName)
	}

	return provider.GetTools(ctx)
}

// GetAllTools 获取所有提供者的所有工具
func (c *Client) GetAllTools(ctx context.Context) (map[string][]mcp.Tool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string][]mcp.Tool)
	var errs []error

	for name, provider := range c.providers {
		tools, err := provider.GetTools(ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("获取提供者 %s 的工具失败: %w", name, err))
			continue
		}
		result[name] = tools
	}

	if len(errs) > 0 {
		return result, fmt.Errorf("获取工具时发生错误: %v", errs)
	}

	return result, nil
}

// CallTool 调用指定提供者的指定工具
func (c *Client) CallTool(ctx context.Context, providerName, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	c.mu.RLock()
	provider, ok := c.providers[providerName]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("提供者 %s 不存在", providerName)
	}

	return provider.CallTool(ctx, toolName, arguments)
}

// Close 关闭所有提供者连接
func (c *Client) Close() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var errs []error
	for name, provider := range c.providers {
		if err := provider.Close(); err != nil {
			errs = append(errs, fmt.Errorf("关闭提供者 %s 失败: %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("关闭提供者时发生错误: %v", errs)
	}

	return nil
}
