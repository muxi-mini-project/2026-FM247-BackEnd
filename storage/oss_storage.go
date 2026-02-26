package storage

// OSSStorage OSS存储实现
// type OSSStorage struct {
// 	bucket *oss.Bucket
// 	config *config.OSSConfig
// }

// // NewOSSStorage 初始化OSS存储
// func NewOSSStorage(cfg *config.OSSConfig) (Storage, error) {
// 	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bucket, err := client.Bucket(cfg.BucketName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &OSSStorage{
// 		bucket: bucket,
// 		config: cfg,
// 	}, nil
// }

// // UploadAvatar 上传头像
// func (s *OSSStorage) UploadAvatar(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (avatarURL string, err error) {
// 	// 1. 验证文件
// 	if err := s.validateFile(fileHeader); err != nil {
// 		return "", err
// 	}

// 	// 2. 生成文件名
// 	objectName := s.generateObjectName(userID, fileHeader.Filename)

// 	// 3. 上传到OSS
// 	err = s.bucket.PutObject(objectName, file)
// 	if err != nil {
// 		return "", fmt.Errorf("上传到OSS失败: %v", err)
// 	}

// 	// 4. 返回访问URL
// 	return s.generateURL(objectName), nil
// }

// // 验证文件
// func (s *OSSStorage) validateFile(fileHeader *multipart.FileHeader) error {
// 	// 文件大小：不超过5MB
// 	if fileHeader.Size > 5*1024*1024 {
// 		return errors.New("图片大小不能超过5MB")
// 	}

// 	// 文件类型
// 	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
// 	allowedExts := map[string]bool{
// 		".jpg":  true,
// 		".jpeg": true,
// 		".png":  true,
// 		".gif":  true,
// 	}

// 	if !allowedExts[ext] {
// 		return errors.New("只支持 jpg, jpeg, png, gif 格式的图片")
// 	}

// 	return nil
// }

// // 生成对象名
// func (s *OSSStorage) generateObjectName(userID uint, filename string) string {
// 	ext := filepath.Ext(filename)
// 	timestamp := time.Now().Format("20060102_150405")
// 	return fmt.Sprintf("avatars/%d_%s%s", userID, timestamp, ext)
// }

// // 生成访问URL
// func (s *OSSStorage) generateURL(objectName string) string {
// 	return fmt.Sprintf("https://%s.%s/%s",
// 		s.config.BucketName,
// 		s.config.Endpoint,
// 		objectName)
// }

// // 测试连接
// func (s *OSSStorage) TestConnection() error {
// 	// 尝试上传一个测试文件
// 	testContent := strings.NewReader("test")
// 	err := s.bucket.PutObject("test-connection.txt", testContent)
// 	if err != nil {
// 		return fmt.Errorf("OSS连接测试失败: %v", err)
// 	}

// 	// 清理测试文件
// 	s.bucket.DeleteObject("test-connection.txt")
// 	return nil
// }

// // DeleteAvatar 实现接口方法（必须添加）
// func (s *OSSStorage) DeleteAvatar(ctx context.Context, avatarURL string) error {
// 	// 从完整URL中提取OSS对象名
// 	// 例如: https://bucket.oss-cn-hangzhou.aliyuncs.com/avatars/1_test.png
// 	// 提取出: avatars/1_test.png

// 	parts := strings.Split(avatarURL, "/")
// 	if len(parts) < 4 {
// 		return fmt.Errorf("无效的头像URL")
// 	}

// 	objectName := strings.Join(parts[3:], "/")
// 	return s.bucket.DeleteObject(objectName)
// }
