package handle

import (
	"2026-FM247-BackEnd/model"
	"2026-FM247-BackEnd/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userservice *service.UserService
}

func NewUserHandler(userservice *service.UserService) *UserHandler {
	return &UserHandler{userservice: userservice}
}

// 修改密码
func (h *UserHandler) UpdatePasswordHandler(c *gin.Context) {
	var req model.UpdatePassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}
	principal, err := GetPrincipal(c)
	if principal == nil {
		FailWithMessage(c, "无法获取用户信息")
		return
	}

	msg := h.userservice.UpdateUserPassword(principal.UserID, req.OldPassword, req.NewPassword)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}

	OkWithMessage(c, "密码修改成功")
}

// 修改用户信息
func (h *UserHandler) UpdatedUserInfoHandler(c *gin.Context) {
	var req model.UpdateUserInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		FailWithMessage(c, "请求参数有误")
		return
	}
	if req.Username == "" && req.Telenum == "" {
		FailWithMessage(c, "用户名和电话号码不能同时为空")
		return
	}
	principal, err := GetPrincipal(c)
	if principal == nil {
		FailWithMessage(c, "无法通过上下文获取用户信息")
		return
	}
	msg := h.userservice.UpdateUserInfo(principal.UserID, req.Username, req.Telenum)
	if msg != "" {
		FailWithMessage(c, msg)
		return
	}
	OkWithMessage(c, "用户信息修改成功")
}
