package service

import (
	"errors"
	"time"

	"2026-FM247-BackEnd/model"
	"2026-FM247-BackEnd/repository"
	"2026-FM247-BackEnd/utils"
)

type UserService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenBlacklistRepository
}

func NewUserService(userRepo *repository.UserRepository, tokenRepo *repository.TokenBlacklistRepository) *UserService {
	return &UserService{userRepo: userRepo, tokenRepo: tokenRepo}
}

// 注册
func (u *UserService) Register(username, password, email string) (err error, message string) {
	if username == "" || password == "" || email == "" {
		return errors.New("用户名、密码和邮箱不能为空"), "用户名、密码和邮箱不能为空"
	}
	_, err = u.userRepo.GetUserByEmail(email)
	if err == nil {
		return errors.New("用户已存在"), "用户已存在"
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("密码加密失败"), "密码加密失败"
	}
	newUser := &model.User{
		Username: username,
		Password: hashedPassword,
		Email:    email,
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
func (u *UserService) UpdateUserInfo(userID uint, username, telenum string) (message string) {
	if username == "" && telenum == "" {
		return "未修改任何信息"
	}
	err := u.userRepo.UpdateUserInfo(userID, username, telenum)
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

// 待办：更新头像
