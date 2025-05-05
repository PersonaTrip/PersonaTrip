package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"personatrip/internal/app"
	"personatrip/internal/models"
	"personatrip/internal/repository"
	"personatrip/internal/utils/logger"
	"personatrip/pkg/appclient"
	"personatrip/pkg/mcp"
)

func main() {
	fmt.Println("PersonaTrip 统一客户端演示")
	fmt.Println("==========================")

	// 创建 context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	setupSignalHandler(cancel)

	// 初始化数据库连接
	mysqlDB, err := repository.NewMySQL("root:password@tcp(localhost:3306)/personatrip?parseTime=true")
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化仓库
	modelConfigRepo := repository.NewGormModelConfigRepository(mysqlDB.DB)
	adminRepo := repository.NewGormAdminRepository(mysqlDB.DB)

	// 创建统一客户端
	client := appclient.NewAppClient(modelConfigRepo, adminRepo)

	// 初始化客户端
	fmt.Println("\n初始化客户端...")
	if err := client.Init(ctx); err != nil {
		log.Printf("客户端初始化失败: %v", err)
		// 尝试创建默认配置
		fmt.Println("尝试创建默认配置...")
		createDefaultConfig(ctx, client)

		// 重新初始化
		if err := client.Init(ctx); err != nil {
			log.Fatalf("客户端重新初始化失败: %v", err)
		}
	}
	defer client.Close()

	// 打印当前活跃的模型信息
	printActiveModel(client)

	// 菜单循环
	for {
		printMenu()
		var choice int
		fmt.Print("请选择操作: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// 获取所有模型配置
			listAllModelConfigs(ctx, client)
		case 2:
			// 切换活跃模型
			switchActiveModel(ctx, client)
		case 3:
			// 生成文本
			generateText(ctx, client)
		case 4:
			// 使用地图工具生成文本
			generateTextWithMapTool(ctx, client)
		case 5:
			// 使用所有工具生成文本
			generateTextWithAllTools(ctx, client)
		case 6:
			// 创建模型配置
			createModelConfig(ctx, client)
		case 0:
			// 退出
			fmt.Println("感谢使用，再见！")
			return
		default:
			fmt.Println("无效的选择，请重试")
		}

		fmt.Println("\n按回车键继续...")
		fmt.Scanln()
	}
}

// 设置信号处理
func setupSignalHandler(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\n收到中断信号，正在退出...")
		cancel()
		os.Exit(0)
	}()
}

// 打印菜单
func printMenu() {
	fmt.Println("\n-- 主菜单 --")
	fmt.Println("1. 列出所有模型配置")
	fmt.Println("2. 切换活跃模型")
	fmt.Println("3. 生成文本")
	fmt.Println("4. 使用地图工具生成文本")
	fmt.Println("5. 使用所有工具生成文本")
	fmt.Println("6. 创建新的模型配置")
	fmt.Println("0. 退出")
}

// 打印当前活跃的模型信息
func printActiveModel(client *appclient.AppClient) {
	config := client.GetActiveModelConfig()
	if config == nil {
		fmt.Println("没有活跃的模型配置")
		return
	}

	fmt.Println("\n当前活跃的模型配置:")
	fmt.Printf("ID: %d\n", config.ID)
	fmt.Printf("名称: %s\n", config.Name)
	fmt.Printf("类型: %s\n", config.ModelType)
	fmt.Printf("模型: %s\n", config.ModelName)
	fmt.Printf("基础URL: %s\n", config.BaseUrl)
	fmt.Printf("Temperature: %.2f\n", config.Temperature)
	fmt.Printf("MaxTokens: %d\n", config.MaxTokens)
}

// 获取所有模型配置
func listAllModelConfigs(ctx context.Context, client *appclient.AppClient) {
	fmt.Println("\n获取所有模型配置...")
	configs, err := client.GetAllModelConfigs(ctx)
	if err != nil {
		fmt.Printf("获取模型配置失败: %v\n", err)
		return
	}

	fmt.Printf("找到 %d 个模型配置:\n", len(configs))
	for i, config := range configs {
		fmt.Printf("%d. %s (ID: %d) - %s %s - 活跃: %v\n",
			i+1, config.Name, config.ID, config.ModelType, config.ModelName, config.IsActive)
	}
}

// 切换活跃模型
func switchActiveModel(ctx context.Context, client *appclient.AppClient) {
	// 获取所有模型
	configs, err := client.GetAllModelConfigs(ctx)
	if err != nil {
		fmt.Printf("获取模型配置失败: %v\n", err)
		return
	}

	if len(configs) == 0 {
		fmt.Println("没有可用的模型配置")
		return
	}

	fmt.Println("\n-- 可用的模型配置 --")
	for i, config := range configs {
		fmt.Printf("%d. %s (ID: %d) - %s %s %v\n",
			i+1, config.Name, config.ID, config.ModelType, config.ModelName, config.IsActive)
	}

	var choice int
	fmt.Print("请选择要激活的模型 (1-", len(configs), "): ")
	fmt.Scanln(&choice)

	if choice < 1 || choice > len(configs) {
		fmt.Println("无效的选择")
		return
	}

	selectedConfig := configs[choice-1]
	err = client.UpdateActiveModel(ctx, selectedConfig.ID)
	if err != nil {
		fmt.Printf("更新活跃模型失败: %v\n", err)
		return
	}

	fmt.Printf("已将 %s (ID: %d) 设置为活跃模型\n", selectedConfig.Name, selectedConfig.ID)
	printActiveModel(client)
}

// 生成文本
func generateText(ctx context.Context, client *appclient.AppClient) {
	fmt.Print("\n请输入提示语: ")
	var prompt string
	fmt.Scanln(&prompt)
	fmt.Println("正在生成文本...")

	text, err := client.GenerateText(ctx, prompt)
	if err != nil {
		fmt.Printf("生成文本失败: %v\n", err)
		return
	}

	fmt.Println("\n--- 生成的文本 ---")
	fmt.Println(text)
	fmt.Println("--- 文本结束 ---")
}

// 使用地图工具生成文本
func generateTextWithMapTool(ctx context.Context, client *appclient.AppClient) {
	fmt.Print("\n请输入地图相关查询: ")
	var prompt string
	fmt.Scanln(&prompt)
	fmt.Println("正在使用地图工具生成文本...")

	text, err := client.GenerateTextWithMap(ctx, prompt)
	if err != nil {
		fmt.Printf("生成文本失败: %v\n", err)
		return
	}

	fmt.Println("\n--- 生成的文本 ---")
	fmt.Println(text)
	fmt.Println("--- 文本结束 ---")
}

// 使用所有工具生成文本
func generateTextWithAllTools(ctx context.Context, client *appclient.AppClient) {
	fmt.Print("\n请输入查询: ")
	var prompt string
	fmt.Scanln(&prompt)
	fmt.Println("正在使用所有工具生成文本...")

	text, err := client.GenerateTextWithAllTools(ctx, prompt)
	if err != nil {
		fmt.Printf("生成文本失败: %v\n", err)
		return
	}

	fmt.Println("\n--- 生成的文本 ---")
	fmt.Println(text)
	fmt.Println("--- 文本结束 ---")
}

// 创建模型配置
func createModelConfig(ctx context.Context, client *appclient.AppClient) {
	fmt.Println("\n创建新的模型配置")

	var config models.ModelConfig

	fmt.Print("模型名称: ")
	fmt.Scanln(&config.Name)

	fmt.Println("模型类型 (openai, ollama, ark, mock): ")
	fmt.Scanln(&config.ModelType)

	fmt.Print("模型: ")
	fmt.Scanln(&config.ModelName)

	fmt.Print("API Key: ")
	fmt.Scanln(&config.ApiKey)

	fmt.Print("Base URL: ")
	fmt.Scanln(&config.BaseUrl)

	fmt.Print("Temperature (0.0-1.0): ")
	fmt.Scanln(&config.Temperature)

	fmt.Print("Max Tokens: ")
	fmt.Scanln(&config.MaxTokens)

	fmt.Print("是否设为活跃 (true/false): ")
	fmt.Scanln(&config.IsActive)

	err := client.CreateModelConfig(ctx, &config)
	if err != nil {
		fmt.Printf("创建模型配置失败: %v\n", err)
		return
	}

	fmt.Println("模型配置创建成功")
}

// 创建默认配置
func createDefaultConfig(ctx context.Context, client *appclient.AppClient) {
	// 创建默认的ARK大模型配置
	defaultConfig := &models.ModelConfig{
		Name:        "默认ARK配置",
		ModelType:   "ark",
		ModelName:   "ep-20250408220714-wzgtv",
		ApiKey:      "106dc02a-bd3b-41cb-a739-4e7301f48385",
		BaseUrl:     "https://ark.cn-beijing.volces.com/api/v3",
		IsActive:    true,
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	err := client.CreateModelConfig(ctx, defaultConfig)
	if err != nil {
		logger.Errorf("创建默认大模型配置失败: %v", err)
	} else {
		logger.Info("默认大模型配置创建成功")
	}
}
