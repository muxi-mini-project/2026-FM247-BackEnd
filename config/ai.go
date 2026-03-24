package config

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func NewAIClient() (*openai.Client, error) {
	if AppConfig.ApiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is not set")
	}
	config := openai.DefaultConfig(AppConfig.ApiKey)
	config.BaseURL = AppConfig.AiBaseURL
	client := openai.NewClientWithConfig(config)
	return client, nil
}
