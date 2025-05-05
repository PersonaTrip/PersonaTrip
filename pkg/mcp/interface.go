package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// MCPProvider 是MCP工具提供者接口
type MCPProvider interface {
	// GetTools 返回该提供者支持的所有工具
	GetTools(ctx context.Context) ([]mcp.Tool, error)
	// Initialize 初始化该提供者
	Initialize(ctx context.Context) error
	// Close 关闭连接并释放资源
	Close() error
	// CallTool 调用指定的工具
	CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error)
}

// MCPClient 管理所有MCP提供者的客户端
type MCPClient interface {
	// AddProvider 添加一个MCP工具提供者
	AddProvider(name string, provider MCPProvider)
	// Initialize 初始化所有提供者
	Initialize(ctx context.Context) error
	// GetTools 获取指定提供者的所有工具
	GetTools(ctx context.Context, providerName string) ([]mcp.Tool, error)
	// GetAllTools 获取所有提供者的所有工具
	GetAllTools(ctx context.Context) (map[string][]mcp.Tool, error)
	// CallTool 调用指定提供者的指定工具
	CallTool(ctx context.Context, providerName, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error)
	// Close 关闭所有提供者连接
	Close() error
}
