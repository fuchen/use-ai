package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// Config 保存用户配置
type Config struct {
	OpenAI struct {
		Endpoint    string `json:"endpoint"`
		Model       string `json:"model"`
		APIKey      string `json:"api_key"`
		SystemPrompt string `json:"system_prompt"`
	} `json:"openai"`
}

// 默认配置
const defaultConfig = `{
  "openai": {
    "endpoint": "https://api.openai.com/v1",
    "model": "gpt-4-turbo",
    "api_key": "YOUR_API_KEY_HERE",
    "system_prompt": "You are a professional command line assistant. Based on the user's query, only return the appropriate shell command without explanation. Make sure the command is suitable for the user's operating system and shell environment."
  }
}`

// LoadConfig 加载用户配置文件
func LoadConfig() (*Config, error) {
	// 获取用户主目录
	home, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get home directory: %w", err)
	}

	// 配置文件路径
	configPath := filepath.Join(home, ".ask-ai.json")

	// 检查配置文件是否存在
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		// 创建默认配置文件
		err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
		if err != nil {
			return nil, fmt.Errorf("Failed to create default config file: %w", err)
		}
		fmt.Printf("Default configuration file created at %s. Please edit this file to add your OpenAI API key, then run the program again.\n", configPath)
		os.Exit(0)
	} else if err != nil {
		return nil, fmt.Errorf("Error checking config file: %w", err)
	}

	// 直接读取文件内容
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file: %w", err)
	}

	// 解析JSON
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse config file: %w", err)
	}

	// 检查必要字段
	if config.OpenAI.Endpoint == "" {
		return nil, fmt.Errorf("Missing 'endpoint' field in config file")
	}
	if config.OpenAI.Model == "" {
		return nil, fmt.Errorf("Missing 'model' field in config file")
	}
	if config.OpenAI.APIKey == "" {
		return nil, fmt.Errorf("Missing 'api_key' field in config file")
	}
	if config.OpenAI.SystemPrompt == "" {
		return nil, fmt.Errorf("Missing 'system_prompt' field in config file")
	}

	return &config, nil
}