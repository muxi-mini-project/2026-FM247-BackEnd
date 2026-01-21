package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) Storage {
	// 确保目录存在
	os.MkdirAll(filepath.Join(basePath, "avatars"), 0755)

	return &LocalStorage{
		basePath: basePath,
	}
}

func (s *LocalStorage) UploadAvatar(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// 生成文件名
	filename := fmt.Sprintf("%d_%d%s",
		userID,
		time.Now().Unix(),
		filepath.Ext(fileHeader.Filename))

	filePath := filepath.Join(s.basePath, "avatars", filename)

	// 创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	// 返回URL路径
	return fmt.Sprintf("/uploads/avatars/%s", filename), nil
}

func (s *LocalStorage) DeleteAvatar(ctx context.Context, avatarURL string) error {
	// 实现删除逻辑
	return nil
}
