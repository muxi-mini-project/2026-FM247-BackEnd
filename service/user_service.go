package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"2026-FM247-BackEnd/config"
	"2026-FM247-BackEnd/models"
	repository "2026-FM247-BackEnd/repositories"
	"2026-FM247-BackEnd/storage"
	"2026-FM247-BackEnd/utils"
)

type UserService struct {
	userRepo  *repository.UserRepository
	storage   storage.Storage
	tokenRepo *repository.TokenBlacklistRepository
}

type IUserService interface {
	Register(username, password string) error
	Login(username, password string) (*models.User, error)
	CancelUser(userID uint, password string) (error, string)
	UpdateUserPassword(userID uint, oldPassword, newPassword string) string
	UpdateUserInfo(userID uint, username, telenum string) string
	GetUserByID(userID uint) (*models.User, error)
}

func NewUserService(userRepo *repository.UserRepository, tokenRepo *repository.TokenBlacklistRepository) *UserService {
	// 1. 加载OSS配置
	ossConfig := config.LoadOSSConfig()

	// 2. 创建存储实例
	var storageImpl storage.Storage
	var err error

	if ossConfig.IsValid() {
		storageImpl, err = storage.NewOSSStorage(ossConfig)
		if err != nil {
			fmt.Printf("警告：OSS初始化失败，将使用本地存储: %v\n", err)
			storageImpl = NewSimpleLocalStorage()
		} else {
			fmt.Println("已启用OSS存储")
		}
	} else {
		storageImpl = NewSimpleLocalStorage()
		fmt.Println("已启用本地存储（开发模式）")
	}

	return &UserService{
		userRepo:  userRepo,
		storage:   storageImpl,
		tokenRepo: tokenRepo,
	}
}

type localStorageImpl struct{}

func NewSimpleLocalStorage() storage.Storage {
	return &localStorageImpl{}
}

// 注册
func (u *UserService) Register(username, password, email, gender string) (err error, message string) {
	if username == "" || password == "" || email == "" {
		return errors.New("用户名、密码和邮箱不能为空"), "用户名、密码和邮箱不能为空"
	}

	if gender != "" && gender != "男" && gender != "女" && gender != "草履虫" {
		return errors.New("性别参数无效"), "性别参数无效，可选值：男、女、草履虫"
	}

	// 如果没有传入性别，设置默认值
	if gender == "" {
		gender = "草履虫"
	}

	_, err = u.userRepo.GetUserByEmail(email)
	if err == nil {
		return errors.New("用户已存在"), "用户已存在"
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("密码加密失败"), "密码加密失败"
	}

	defaultAvatar := u.getDefaultAvatarURL()

	newUser := &models.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Gender:   gender,
		Avatar:   defaultAvatar,
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
func (u *UserService) UpdateUserEmail(userid uint, newEmail string) (message string) {
	_, err := u.userRepo.GetUserByEmail(newEmail)
	if err == nil {
		return "该邮箱已被使用"
	}
	err = u.userRepo.UpdateUserEmail(userid, newEmail)
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

// ID查找用户
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}

// 上传头像方法
func (s *localStorageImpl) UploadAvatar(ctx context.Context, userID uint, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// 确保上传目录存在
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// 生成文件名
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return "/uploads/avatars/" + filename, nil
}

func (s *localStorageImpl) DeleteAvatar(ctx context.Context, avatarURL string) error {
	return nil
}

// 上传头像方法
func (s *UserService) UploadAvatar(userID uint, file multipart.File, fileHeader *multipart.FileHeader) (avatarURL string, err error) {
	// 验证用户是否存在
	_, err = s.userRepo.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("用户不存在: %v", err)
	}

	// 上传到存储
	avatarURL, err = s.storage.UploadAvatar(context.Background(), userID, file, fileHeader)
	if err != nil {
		return "", fmt.Errorf("上传失败: %v", err)
	}

	// 更新数据库
	err = s.UpdateUserAvatar(userID, avatarURL)
	if err != nil {
		return "", fmt.Errorf("更新数据库失败: %v", err)
	}

	return avatarURL, nil
}

// 更新用户头像URL到数据库
func (s *UserService) UpdateUserAvatar(userID uint, avatarURL string) error {
	// 直接调用数据库更新
	return s.userRepo.UpdateAvatarURL(userID, avatarURL)
}

func (s *UserService) getDefaultAvatarURL() string {
	// 使用接口检查而不是具体类型断言
	_, isLocal := s.storage.(*localStorageImpl)

	if !isLocal {
		cfg := config.LoadOSSConfig()
		if cfg.IsValid() {
			return fmt.Sprintf("https://%s.%s/default-avatar.png",
				cfg.BucketName, cfg.Endpoint)
		}
	}

	// 本地默认头像
	return "/static/default-avatar.png"
}
