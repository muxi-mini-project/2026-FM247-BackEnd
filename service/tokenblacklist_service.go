package service

import (
	"time"
)

type TokenBlacklistRepository interface {
	AddToBlacklist(jti string, expiresAt time.Time) error
	IsBlacklisted(jti string) (bool, error)
}

type TokenBlacklistService struct {
	blacklistRepo TokenBlacklistRepository
}

// 只是封装一下
func NewTokenBlacklistService(blacklistRepo TokenBlacklistRepository) *TokenBlacklistService {
	return &TokenBlacklistService{blacklistRepo: blacklistRepo}
}

func (s *TokenBlacklistService) AddToBlacklist(jti string) error {
	return s.blacklistRepo.AddToBlacklist(jti, time.Now())
}

func (s *TokenBlacklistService) IsBlacklisted(jti string) (bool, error) {
	return s.blacklistRepo.IsBlacklisted(jti)
}
