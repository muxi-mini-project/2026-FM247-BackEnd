package handle

import (
	"2026-FM247-BackEnd/model"
	"errors"

	"github.com/gin-gonic/gin"
)

func GetPrincipal(c *gin.Context) (*model.Principal, error) {
	principal, exists := c.Get("principal")
	if !exists {
		return nil, errors.New("用户信息不存在")
	}
	p, ok := principal.(*model.Principal)
	if !ok {
		return nil, errors.New("用户信息类型错误")
	}
	return p, nil
}
