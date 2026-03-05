package utils

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestValidateUsername(t *testing.T) {
	// 1. 定义测试用例表
	tests := []struct {
		name     string // 用例描述
		username string // 输入
		want     bool   // 期望结果
	}{
		{"正常用户名", "User_123", true},
		{"过短", "a", false},
		{"非法字符", "User@123", false},
		{"中文支持", "测试用户", true},
	}

	// 2. 遍历执行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUsername(tt.username)
			// 3. 断言结果 (使用 testify 库更简洁)
			assert.Equal(t, tt.want, got)
		})
	}
}
