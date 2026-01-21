package storage

import (
	"context"
	"mime/multipart"
)

// Storage 接口定义
type Storage interface {
	UploadAvatar(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error)
	DeleteAvatar(ctx context.Context, avatarURL string) error // 确保方法名一致
}
