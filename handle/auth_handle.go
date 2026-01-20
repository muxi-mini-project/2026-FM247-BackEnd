package handle

import (
	"2026-FM247-BackEnd/model"
	"2026-FM247-BackEnd/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	tokenservice *service.TokenBlacklistService
	userservice  *service.UserService
}

func NewAuthHandler(tokenservice *service.TokenBlacklistService, userservice *service.UserService) *AuthHandler {
	return &AuthHandler{
		tokenservice: tokenservice,
		userservice:  userservice,
	}
}

// LoginHandler 登录
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req model.LoginUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}

	token, msg := h.userservice.Login(req.Email, req.Password)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+token)

	OkWithMessage(c, msg)
}

// LogoutHandler 登出
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	p, err := GetPrincipal(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息")
		return
	}
	err = h.tokenservice.AddToBlacklist(p.Jti)
	if err != nil {
		FailWithMessage(c, "登出失败")
		return
	}
	OkWithMessage(c, "登出成功")
}

// CancelHandler 注销用户
func (h *AuthHandler) CancelHandler(c *gin.Context) {
	var req model.CancelUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := GetPrincipal(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息")
		return
	}
	err, msg := h.userservice.CancelUser(principal.UserID, req.Password)
	if err != nil {
		FailWithMessage(c, msg)
		return
	}

	OkWithMessage(c, "注销成功")
}

// RegisterUserHandler 注册新用户
func (h *AuthHandler) RegisterUserHandler(c *gin.Context) {
	var req model.RegisterUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}

	err, msg := h.userservice.Register(req.Username, req.Password, req.Email)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}
	OkWithMessage(c, "注册成功,请牢记账户和密码")
}
