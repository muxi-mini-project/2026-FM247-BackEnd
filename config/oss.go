package config

import (
	"os"
)

type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

func LoadOSSConfig() *OSSConfig {
	return &OSSConfig{
		Endpoint:        getOssEnv("OSS_ENDPOINT"),
		AccessKeyID:     getOssEnv("OSS_ACCESS_KEY_ID"),
		AccessKeySecret: getOssEnv("OSS_ACCESS_KEY_SECRET"),
		BucketName:      getOssEnv("OSS_BUCKET_NAME"),
	}
}

// IsValid 验证配置是否完整
func (c *OSSConfig) IsValid() bool {
	return c.Endpoint != "" &&
		c.AccessKeyID != "" &&
		c.AccessKeySecret != "" &&
		c.BucketName != ""
}

func getOssEnv(key string) string {
	return os.Getenv(key)
}
