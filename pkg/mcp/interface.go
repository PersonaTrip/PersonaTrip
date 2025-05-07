package mcp

import (
	"context"
	"github.com/cloudwego/eino/components/tool"

	"github.com/mark3labs/mcp-go/mcp"
)

// MCPProvider 是MCP工具提供者接口
type MCPProvider interface {
	// GetTools 返回该提供者支持的所有工具
	GetTools(ctx context.Context) ([]tool.BaseTool, error)
	// Initialize 初始化该提供者
	Initialize(ctx context.Context) error
	// Close 关闭连接并释放资源
	Close() error
	// CallTool 调用指定的工具
	CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error)
}
