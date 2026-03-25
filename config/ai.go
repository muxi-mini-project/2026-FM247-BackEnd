package config

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func NewAIClient() (*openai.Client, error) {
	if AppConfig.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is not set")
	}
	config := openai.DefaultConfig(AppConfig.APIKey)
	config.BaseURL = AppConfig.AIBaseURL
	client := openai.NewClientWithConfig(config)
	return client, nil
}
