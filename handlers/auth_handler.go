package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(username, password, email string) (err error, message string)
	Login(email, password string) (token string, message string)
	Logout(jti string) (err error, message string)
	CancelUser(userID uint, password string) (err error, message string)
	UpdateUserPassword(userID uint, oldPassword, newPassword string) (message string)
	UpdateUserEmail(userID uint, newEmail string, password string) (message string)
	UpdateUserInfo(userID uint, username, telenum, gender string) (message string)
	GetUserInfo(userID uint) (*service.UserInfo, error)
}

type TokenService interface {
	AddToBlacklist(jti string) error
	IsBlacklisted(jti string) (bool, error)
}

type AuthHandler struct {
	Tokenservice TokenService
	Userservice  UserService
}

func NewAuthHandler(tokenservice TokenService, userservice UserService) *AuthHandler {
	return &AuthHandler{
		Tokenservice: tokenservice,
		Userservice:  userservice,
	}
}

// RegisterUserHandler 注册新用户
// @Router /api/auth/register [post]
func (h *AuthHandler) RegisterUserHandler(c *gin.Context) {
	var req RegisterUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}

	err, msg := h.Userservice.Register(req.Username, req.Password, req.Email)
	if msg != "注册成功" {
		FailWithMessage(c, msg)
		return
	}
	OkWithMessage(c, "注册成功,请牢记账户和密码")
}

// LoginHandler 登录
// @Router /api/auth/login [post]
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}

	token, msg := h.Userservice.Login(req.Email, req.Password)
	if msg != "登录成功" {
		FailWithMessage(c, msg)
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+token)

	Ok(c, msg, gin.H{
		"Authorization": "Bearer " + token,
	})

}

// LogoutHandler 登出
// @Router /api/auth/logout [post]
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	// 方法1：使用GetClaimsFromContext获取完整的claims
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	err = h.Tokenservice.AddToBlacklist(claims.Jti)
	if err != nil {
		FailWithMessage(c, "登出失败: "+err.Error())
		return
	}

	OkWithMessage(c, "登出成功")
}

// CancelHandler 注销用户
// @Router /api/auth/cancel [post]
func (h *AuthHandler) CancelHandler(c *gin.Context) {
	var req CancelUser
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}

	// 获取完整的claims
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	err, msg := h.Userservice.CancelUser(claims.UserID, req.Password)
	if err != nil {
		FailWithMessage(c, msg)
		return
	}

	// 注销成功后，吊销当前令牌
	_ = h.Tokenservice.AddToBlacklist(claims.Jti)

	OkWithMessage(c, "注销成功")
}

// UpdatePasswordHandler 修改密码
// @Router /api/user/update_password [post]
func (h *AuthHandler) UpdatePasswordHandler(c *gin.Context) {
	var req UpdatePassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}

	// 使用 GetClaimsFromContext
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	msg := h.Userservice.UpdateUserPassword(claims.UserID, req.OldPassword, req.NewPassword)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}

	OkWithMessage(c, "密码修改成功")
}

// UpdateEmailHandler 修改邮箱
// @Router /api/user/update_email [post]
func (h *AuthHandler) UpdateEmailHandler(c *gin.Context) {
	var req UpdateEmail
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}

	// 使用 GetClaimsFromContext
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	msg := h.Userservice.UpdateUserEmail(claims.UserID, req.NewEmail, req.Password)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}

	OkWithMessage(c, "密码修改成功")
}

// UpdateUserInfoHandler 修改用户信息
// @Router /api/user/update_info [post]
func (h *AuthHandler) UpdateUserInfoHandler(c *gin.Context) {

	// 使用 GetClaimsFromContext
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	var req UpdateUserInfo
	err = c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误: "+err.Error())
		return
	}

	// 验证请求参数
	if req.Username == "" && req.Telenum == "" && req.Gender == "" {
		FailWithMessage(c, "未修改任何信息")
		return
	}
	//验证用户名格式
	if req.Username != "" && !utils.ValidateUsername(req.Username) {
		FailWithMessage(c, "用户名格式不正确")
		return
	}

	// 验证电话号码格式
	if req.Telenum != "" && !utils.ValidatePhoneNumber(req.Telenum) {
		FailWithMessage(c, "电话号码格式不正确")
		return
	}

	msg := h.Userservice.UpdateUserInfo(claims.UserID, req.Username, req.Telenum, req.Gender)
	if msg != "更新用户信息成功" {
		FailWithMessage(c, msg)
		return
	}

	OkWithMessage(c, msg)
}

// GetUserInfoHandler 获取当前用户信息
// @Router /api/user/info [get]
func (h *AuthHandler) GetUserInfoHandler(c *gin.Context) {
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "无法获取用户信息: "+err.Error())
		return
	}

	// 从数据库获取完整的用户信息
	userinfo, err := h.Userservice.GetUserInfo(claims.UserID)
	if err != nil {
		FailWithMessage(c, "获取用户信息失败")
		return
	}

	OkWithData(c, userinfo)
}
