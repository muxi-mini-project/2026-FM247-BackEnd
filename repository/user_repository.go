package repository

import (
	"2026-FM247-BackEnd/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := r.db.Create(user)
	if result != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserInfo(updateInfo *model.UpdateUserInfo) error {
	var user *model.User
	user, err := r.GetUserByID(updateInfo.UserID)
	if err != nil {
		return err
	}
	if updateInfo.Username != "" {
		user.Username = updateInfo.Username
	}
	if updateInfo.Telenum != "" {
		user.Telenum = updateInfo.Telenum
	}

	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) UpdateUserEmail(updateemail *model.UpdateEmail) error {
	var user *model.User
	user, err := r.GetUserByID(updateemail.UserID)
	if err != nil {
		return err
	}
	user.Email = updateemail.NewEmail
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//上传用户头像,待办

func (r *UserRepository) UpdatePassword(userid uint, newpassword string) error {
	var user *model.User
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
