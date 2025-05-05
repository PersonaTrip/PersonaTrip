package amap

import (
	"context"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Provider 高德地图MCP工具提供者
type Provider struct {
	client    *client.Client
	tools     []mcp.Tool
	toolsLock sync.RWMutex
	apiKey    string
}

// NewProvider 创建高德地图MCP提供者
func NewProvider(apiKey string) *Provider {
	return &Provider{
		apiKey: apiKey,
	}
}

// Initialize 初始化高德地图MCP提供者
func (p *Provider) Initialize(ctx context.Context) error {
	// 创建高德地图MCP客户端
	cli, err := client.NewStdioMCPClient("npx", []string{"AMAP_MAPS_API_KEY=" + p.apiKey}, "-y", "@amap/amap-maps-mcp-server")
	if err != nil {
		return err
	}

	if err := cli.Start(ctx); err != nil {
		return err
	}

	// 发送初始化请求
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "personatrip-client",
		Version: "1.0.0",
	}

	if _, err := cli.Initialize(ctx, initRequest); err != nil {
		cli.Close()
		return err
	}

	p.client = cli

	// 获取并缓存工具列表
	return p.refreshTools(ctx)
}

// refreshTools 刷新工具列表
func (p *Provider) refreshTools(ctx context.Context) error {
	listResult, err := p.client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return err
	}

	p.toolsLock.Lock()
	p.tools = listResult.Tools
	p.toolsLock.Unlock()

	return nil
}

// GetTools 获取高德地图所有可用的工具
func (p *Provider) GetTools(ctx context.Context) ([]mcp.Tool, error) {
	p.toolsLock.RLock()
	defer p.toolsLock.RUnlock()

	// 如果工具列表为空，则刷新
	if len(p.tools) == 0 {
		p.toolsLock.RUnlock()
		if err := p.refreshTools(ctx); err != nil {
			return nil, err
		}
		p.toolsLock.RLock()
	}

	// 复制工具列表
	tools := make([]mcp.Tool, len(p.tools))
	copy(tools, p.tools)

	return tools, nil
}

// CallTool 调用指定的工具
func (p *Provider) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*mcp.CallToolResult, error) {
	request := mcp.CallToolRequest{}
	request.Params.Name = toolName
	request.Params.Arguments = arguments

	return p.client.CallTool(ctx, request)
}

// Close 关闭连接并释放资源
func (p *Provider) Close() error {
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}
