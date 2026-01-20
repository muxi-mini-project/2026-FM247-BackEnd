package middleware

import (
	"2026-FM247-BackEnd/handle"
	"2026-FM247-BackEnd/model"
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenblacklistservice *service.TokenBlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handle.FailWithMessage(c, "请先登录")
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))

		claims, err := utils.ValidateToken(token)
		if err != nil {
			handle.FailWithMessage(c, "token无效或已过期")
			c.Abort()
			return
		}

		ok, err := tokenblacklistservice.IsBlacklisted(claims.Jti)
		if err != nil {
			handle.FailWithMessage(c, "服务器内部错误")
			c.Abort()
			return
		}
		if ok {
			handle.FailWithMessage(c, "token已被注销,请重新登录")
			c.Abort()
			return
		}

		principal := &model.Principal{
			UserID:   claims.UserID,
			Username: claims.Username,
			Jti:      claims.Jti,
		}
		c.Set("principal", principal)
		c.Next()
	}
}
