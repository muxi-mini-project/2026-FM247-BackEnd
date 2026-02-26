package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"2026-FM247-BackEnd/models"
	"2026-FM247-BackEnd/utils"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUserInfo(userID uint, username, telenum, gender string) error
	UpdateUserEmail(userid uint, newEmail string) error
	UpdatePassword(userid uint, newpassword string) error
	DeleteUser(userid uint) error
	UpdateAvatarURL(userID uint, avatarURL string) error
}

type UserService struct {
	userRepo  UserRepository
	storage   Storage
	tokenRepo TokenBlacklistRepository
}

func NewUserService(userRepo UserRepository, tokenRepo TokenBlacklistRepository, storage Storage) *UserService {
	return &UserService{
		userRepo:  userRepo,
		storage:   storage,
		tokenRepo: tokenRepo,
	}
}

// 注册
func (u *UserService) Register(username, password, email string) (err error, message string) {
	if username == "" || password == "" || email == "" {
		return errors.New("用户名、密码和邮箱不能为空"), "用户名、密码和邮箱不能为空"
	}

	if username != "" && !utils.ValidateUsername(username) {
		return errors.New("用户名格式不正确"), "用户名格式不正确, 只能包含汉字、字母、数字、下划线，长度为2到20个字符"
	}

	_, err = u.userRepo.GetUserByEmail(email)
	if err == nil {
		return errors.New("用户已存在"), "用户已存在"
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("密码加密失败"), "密码加密失败"
	}

	newUser := &models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Gender:   "草履虫",
	}
	err = u.userRepo.CreateUser(newUser)
	if err != nil {
		return errors.New("注册失败"), "注册失败"
	}
	return nil, "注册成功"
}

// 登录
func (u *UserService) Login(email, password string) (token string, message string) {
	if email == "" || password == "" {
		return "", "邮箱和密码不能为空"
	}
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", "用户不存在"
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "密码错误"
	}
	token, err = utils.GenerateToken(user)
	if err != nil {
		return "", "生成token失败"
	}
	return token, "登录成功"
}

// 登出
func (u *UserService) Logout(jti string) (err error, message string) {
	err = u.tokenRepo.AddToBlacklist(jti, time.Now())
	if err != nil {
		return err, "登出失败"
	}
	return nil, "登出成功"
}

// 注销
func (u *UserService) CancelUser(userID uint, password string) (err error, message string) {
	if password == "" {
		return nil, "密码不能为空"
	}
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return err, "用户不存在"
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return errors.New("密码错误"), "密码错误"
	}
	err = u.userRepo.DeleteUser(userID)
	if err != nil {
		return err, "注销失败"
	}
	return nil, "注销成功"
}

// 更新用户信息
func (u *UserService) UpdateUserInfo(userID uint, username, telenum, gender string) (message string) {
	if username == "" && telenum == "" && gender == "" {
		return "未修改任何信息"
	}
	err := u.userRepo.UpdateUserInfo(userID, username, telenum, gender)
	if err != nil {
		return "更新用户信息失败"
	}
	return "更新用户信息成功"
}

// 更新用户邮箱
func (u *UserService) UpdateUserEmail(userID uint, newEmail string, password string) (message string) {
	_, err := u.userRepo.GetUserByEmail(newEmail)
	if err == nil {
		return "该邮箱已被使用"
	}

	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return "用户不存在"
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "密码错误"
	}
	err = u.userRepo.UpdateUserEmail(userID, newEmail)
	if err != nil {
		return "更新用户邮箱失败"
	}
	return "更新用户邮箱成功"
}

// 更新用户密码
func (u *UserService) UpdateUserPassword(userID uint, oldPassword, newPassword string) (message string) {
	if oldPassword == "" || newPassword == "" {
		return "旧密码和新密码不能为空"
	}
	if oldPassword == newPassword {
		return "新旧密码不能相同"
	}
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return "用户不存在"
	}
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return "旧密码错误"
	}
	hashedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return "新密码加密失败"
	}
	err = u.userRepo.UpdatePassword(userID, hashedNewPassword)
	if err != nil {
		return "更新用户密码失败"
	}
	return "更新用户密码成功"
}

// ID获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*UserInfo, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	finalAvatarURL := user.Avatar
	if finalAvatarURL != "http://localhost:8080/uploads/default-avatar.png" {
		// 调用 storage 接口获取完整链接
		// 假如 user.Avatar 是 "avatars/1.jpg"
		// 本地存储会变成 "http://localhost:8080/uploads/avatars/1.jpg"
		// OSS 存储会变成 "https://oss-cn-xxx.../avatars/1.jpg"
		fullURL, err := s.storage.GetURL(user.Avatar)
		if err == nil {
			finalAvatarURL = fullURL
		}
	}

	userInfo := &UserInfo{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Gender:     user.Gender,
		Telenum:    user.Telenum,
		Avatar:     finalAvatarURL,
		Experience: user.Experience,
		Level:      user.Level,
		CreatedAt:  user.CreatedAt,
	}
	return userInfo, nil
}

// 上传头像方法。返回文件完整url
func (s *UserService) UploadAvatar(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	//验证文件类型
	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", errors.New("只允许上传图片文件")
	}

	// 生成文件路径，格式为 avatars/{userID}_{timestamp}{ext}
	ext := filepath.Ext(fileHeader.Filename)
	path := fmt.Sprintf("avatars/%d_%d%s", userID, time.Now().Unix(), ext)

	// 保存文件
	url, err := s.storage.Upload(ctx, path, file, fileHeader.Size, contentType)
	if err != nil {
		return "", fmt.Errorf("文件存储失败: %w", err)
	}

	// 更新数据库中的头像URL
	err = s.userRepo.UpdateAvatarURL(userID, url)
	if err != nil {
		return "", fmt.Errorf("更新数据库失败: %w", err)
	}

	// 返回头像URL
	fullURL, err := s.storage.GetURL(url)
	if err != nil {
		return "", fmt.Errorf("获取文件完整URL失败: %w", err)
	}

	return fullURL, nil
}
