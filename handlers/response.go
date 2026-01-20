package handler

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewResponse(c *gin.Context, code int, message string, data any) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// SuccessResponse 成功响应
func Ok(c *gin.Context, message string, data any) {
	NewResponse(c, 200, message, data)
}

func OkWithData(c *gin.Context, data any) {
	Ok(c, "操作成功", data)
}

func OkWithMessage(c *gin.Context, message string) {
	Ok(c, message, gin.H{})
}

// ErrorResponse 失败响应
func Fail(c *gin.Context, code int, message string) {
	NewResponse(c, code, message, gin.H{})
}

func FailWithMessage(c *gin.Context, message string) {
	Fail(c, 400, message)
}

func FailWithCode(c *gin.Context, code int) {
	Fail(c, code, "操作失败")
}
