package handler

import (
	"2026-FM247-BackEnd/service"
	"2026-FM247-BackEnd/utils"

	"github.com/gin-gonic/gin"
)

type AvatarHandler struct {
	userService *service.UserService
}

func NewAvatarHandler(userService *service.UserService) *AvatarHandler {
	return &AvatarHandler{
		userService: userService,
	}
}

// UploadAvatar 上传头像
// @Router /api/user/avatar [post]
func (h *AvatarHandler) UploadAvatar(c *gin.Context) {
	// 1. 验证登录
	claims, err := utils.GetClaimsFromContext(c)
	if err != nil {
		FailWithMessage(c, "请先登录")
		return
	}

	// 2. 获取上传的文件
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		FailWithMessage(c, "请选择头像文件")
		return
	}
	defer file.Close()

	// 3. 检查文件大小
	if header.Size > 5*1024*1024 { // 5MB
		FailWithMessage(c, "图片大小不能超过5MB")
		return
	}

	// 4. 上传头像
	avatarURL, err := h.userService.UploadAvatar(c.Request.Context(), claims.UserID, file, header)
	if err != nil {
		FailWithMessage(c, "上传失败: "+err.Error())
		return
	}

	// 5. 返回结果
	OkWithData(c, gin.H{"avatar_url": avatarURL})
}
