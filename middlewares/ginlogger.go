package middleware

import (
	"time"

	"2026-FM247-BackEnd/logger"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			logger.Log.Errorf("[GIN] %d | %v | %s | %s %s | %s",
				statusCode, latency, clientIP, method, path, errorMsg)
		} else if statusCode >= 400 {
			logger.Log.Warnf("[GIN] %d | %v | %s | %s %s | %s",
				statusCode, latency, clientIP, method, path, errorMsg)
		} else {
			logger.Log.Infof("[GIN] %d | %v | %s | %s %s",
				statusCode, latency, clientIP, method, path)
		}
	}
}
