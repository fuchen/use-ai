package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	appName    = "use-ai"
	appVersion = "0.1.0"
)

func main() {
	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// 检测系统环境
	sysInfo := DetectSystem()

	// 创建LLM客户端
	llmClient := NewLLMClient(config)
	llmClient.AddSystemInfo(sysInfo)

	// 创建输入读取器
	reader := bufio.NewReader(os.Stdin)

	// 开始交互循环
	fmt.Printf("Enter your question about <%s> (press Enter to exit):\n", sysInfo.ShellType)
	for {
		fmt.Print("> ")
		question, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		// 去除输入中的空白字符
		question = strings.TrimSpace(question)

		// 检查是否要退出
		if question == "" {
			break
		}

		// 向LLM提问
		fmt.Println("AI thinking...")
		answer, err := llmClient.Ask(question)
		if err != nil {
			fmt.Printf("Failed to get answer: %v\n", err)
			continue
		}

		// 展示回答
		fmt.Println("\n" + answer + "\n")
	}
}