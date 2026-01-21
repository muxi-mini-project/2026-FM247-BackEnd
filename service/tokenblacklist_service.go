package service

import (
	"2026-FM247-BackEnd/repositories"
	"time"
)

type TokenBlacklistService struct {
	blacklistRepo *repository.TokenBlacklistRepository
}

type ITokenBlacklistService interface {
	AddToBlacklist(jti string) error
	IsBlacklisted(jti string) (bool, error)
}

// 只是封装一下
func NewTokenBlacklistService(blacklistRepo *repository.TokenBlacklistRepository) *TokenBlacklistService {
	return &TokenBlacklistService{blacklistRepo: blacklistRepo}
}

func (s *TokenBlacklistService) AddToBlacklist(jti string) error {
	return s.blacklistRepo.AddToBlacklist(jti, time.Now())
}

func (s *TokenBlacklistService) IsBlacklisted(jti string) (bool, error) {
	return s.blacklistRepo.IsBlacklisted(jti)
}
