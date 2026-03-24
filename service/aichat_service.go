package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

const MaxChatHistoryMessages = 20

type AIChatRepository interface {
	SaveChatHistory(ctx *gin.Context, sessionID uint, newMessages ...openai.ChatCompletionMessage) error
	GetChatHistory(ctx *gin.Context, sessionID uint) ([]openai.ChatCompletionMessage, error)
	TrimChatHistory(ctx *gin.Context, sessionID uint) error
}

type AIChatService struct {
	repo AIChatRepository
	ai   openai.Client
}

func NewAIChatService(repo AIChatRepository, ai openai.Client) *AIChatService {
	return &AIChatService{repo: repo, ai: ai}
}

func (s *AIChatService) Chat(ctx *gin.Context, userID uint, content string) (string, error) {
	userMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: content,
	}

	err := s.repo.SaveChatHistory(ctx, userID, userMsg)
	if err != nil {
		return "", fmt.Errorf("存入聊天记录失败: %w", err)
	}

	aiRequestSuccess := false
	defer func() {
		if !aiRequestSuccess {
			s.repo.TrimChatHistory(ctx, userID)
		}
	}()

	chatHistory, err := s.repo.GetChatHistory(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("获取聊天记录失败: %w", err)
	}

	finalMessages := make([]openai.ChatCompletionMessage, 0, MaxChatHistoryMessages+1)
	systemPrompt := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "你是一个线上自习室的学姐",
	}
	finalMessages = append(finalMessages, systemPrompt)
	startIndex := 0
	if len(chatHistory) > MaxChatHistoryMessages {
		startIndex = len(chatHistory) - MaxChatHistoryMessages
	}
	finalMessages = append(finalMessages, chatHistory[startIndex:]...)

	resp, err := s.ai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       "deepseek-chat",
		Messages:    finalMessages,
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("调用AI接口失败: %w", err)
	}

	aiMsg := resp.Choices[0].Message
	err = s.repo.SaveChatHistory(ctx, userID, aiMsg)
	if err != nil {
		return "", fmt.Errorf("保存AI回复失败: %w", err)
	}
	aiRequestSuccess = true
	return aiMsg.Content, nil
}

func (s *AIChatService) GetChatHistory(ctx *gin.Context, userID uint) ([]openai.ChatCompletionMessage, error) {
	chatHistory, err := s.repo.GetChatHistory(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取聊天记录失败: %w", err)
	}
	return chatHistory, nil
}
