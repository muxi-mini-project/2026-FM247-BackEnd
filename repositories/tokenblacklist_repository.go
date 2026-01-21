package repository

import (
	"2026-FM247-BackEnd/models"
	"time"

	"gorm.io/gorm"
)

type TokenBlacklistRepository struct {
	db *gorm.DB
}

type ITokenBlacklistRepository interface {
	AddToBlacklist(jti string, expiresAt time.Time) error
	IsBlacklisted(jti string) (bool, error)
}

func NewTokenBlacklistRepository(db *gorm.DB) *TokenBlacklistRepository {
	return &TokenBlacklistRepository{db: db}
}

func (r *TokenBlacklistRepository) AddToBlacklist(jti string, expiresAt time.Time) error {
	blacklistEntry := &models.TokenBlacklist{
		Jti:       jti,
		ExpiresAt: expiresAt,
	}
	result := r.db.Create(blacklistEntry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TokenBlacklistRepository) IsBlacklisted(jti string) (bool, error) {
	result := r.db.Where("jti = ?", jti).First(&models.TokenBlacklist{})
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

//待办：定时删除过期的黑名单记录
