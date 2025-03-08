package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// LLMClient 与OpenAI API交互的客户端
type LLMClient struct {
	client       *openai.Client
	model        string
	systemPrompt string
	messages     []openai.ChatCompletionMessage
}

// NewLLMClient 创建一个新的LLM客户端
func NewLLMClient(config *Config) *LLMClient {
	// 创建OpenAI客户端
	clientConfig := openai.DefaultConfig(config.OpenAI.APIKey)

	// 如果配置中有自定义端点，则使用自定义端点
	if config.OpenAI.Endpoint != "https://api.openai.com/v1" {
		clientConfig.BaseURL = config.OpenAI.Endpoint
	}

	client := openai.NewClientWithConfig(clientConfig)

	// 初始化消息历史，添加系统提示
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: config.OpenAI.SystemPrompt,
		},
	}

	return &LLMClient{
		client:       client,
		model:        config.OpenAI.Model,
		systemPrompt: config.OpenAI.SystemPrompt,
		messages:     messages,
	}
}

// AddSystemInfo 添加系统环境信息到消息历史
func (lc *LLMClient) AddSystemInfo(sysInfo SystemInfo) {
	systemInfoMsg := fmt.Sprintf("User's system environment information: %s", sysInfo.GetSystemDescription())
	lc.messages = append(lc.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemInfoMsg,
	})
}

// Ask 向LLM提问并获取回复
func (lc *LLMClient) Ask(question string) (string, error) {
	if strings.TrimSpace(question) == "" {
		return "", errors.New("Question cannot be empty")
	}

	// 添加用户问题到消息历史
	lc.messages = append(lc.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: question,
	})

	// 创建聊天完成请求
	resp, err := lc.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    lc.model,
			Messages: lc.messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// 检查回复是否为空
	if len(resp.Choices) == 0 {
		return "", errors.New("API did not return a valid response")
	}

	// 获取助手回复
	assistantReply := resp.Choices[0].Message.Content

	// 添加助手回复到消息历史
	lc.messages = append(lc.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: assistantReply,
	})

	return assistantReply, nil
}

// ResetConversation 重置对话历史
func (lc *LLMClient) ResetConversation() {
	lc.messages = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: lc.systemPrompt,
		},
	}
}