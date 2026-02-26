package storage

import (
	"context"
	"io"
)

type Storage interface {
	// Upload 上传文件
	// path: 存储路径（例如 "avatars/1.jpg" 或 "music/2026/song.mp3"）
	// content: 文件内容流
	// fileSize: 文件大小（某些OSS SDK需要预知大小）
	// contentType: 文件类型（如 "image/jpeg", "audio/mpeg"）
	Upload(ctx context.Context, path string, content io.Reader, fileSize int64, contentType string) (string, error)

	// Delete 删除文件
	Delete(ctx context.Context, path string) error

	// GetURL 获取文件的完整访问地址（如果是公有读OSS，直接拼接域名；如果是私有，可能需要生成签名URL）
	GetURL(path string) (string, error)
}

func InitStorage(baseURL string) Storage {
	//暂时初始化本地存储，后续可以根据配置切换到OSS等云存储
	return NewLocalStorage("./uploads", baseURL)
}
