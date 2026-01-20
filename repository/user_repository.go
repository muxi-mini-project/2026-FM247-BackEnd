package repository

import (
	"2026-FM247-BackEnd/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserInfo(userID uint, username, telenum string) error {
	var user *models.User
	user, err := r.GetUserByID(userID)
	if err != nil {
		return err
	}
	if username != "" {
		user.Username = username
	}
	if telenum != "" {
		user.Telenum = telenum
	}

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) UpdateUserEmail(userid uint, newEmail string) error {
	var user *models.User
	user, err := r.GetUserByID(userid)
	if err != nil {
		return err
	}
	user.Email = newEmail
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 待办：上传用户头像
func (r *UserRepository) UpdatePassword(userid uint, newpassword string) error {
	var user *models.User
	user, err := r.GetUserByID(userid)
	if err != nil {
		return err
	}
	user.Password = newpassword
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) DeleteUser(userid uint) error {
	result := r.db.Delete(&models.User{}, userid)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
