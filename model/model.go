package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type User struct {
	BaseModel
	Username     string `json:"username" gorm:"column:username"`
	ProfilePhoto string `json:"profile_photo"` // 头像URL
	Password     string `json:"-"`
	Email        string `json:"email"`
	Telenum      string `json:"telenum"`
	Level        string `json:"level"`
	TotalExp     int    `json:"total_exp"` //总经验值
}

type TotalStudyData struct {
	BaseModel
	UserID    uint `json:"user_id"`
	StudyTime int  `json:"study_time"`
	Tomatoes  int  `json:"tomatoes"`
}

type DailyStudyData struct {
	BaseModel
	UserID    uint      `json:"user_id"`
	Date      time.Time `json:"date"`
	EndAt     time.Time `json:"end_at"`
	StudyTime int       `json:"study_time"`
	Tomatoes  int       `json:"tomatoes"`
}

type MonthlyStudyData struct {
	BaseModel
	UserID    uint      `json:"user_id"`
	Month     time.Time `json:"month"`
	EndAt     time.Time `json:"end_at"`
	StudyTime int       `json:"study_time"`
	Tomatoes  int       `json:"tomatoes"`
}

type Todo struct {
	BaseModel
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Deadline    time.Time `json:"deadline"`
}

type Note struct {
	BaseModel
	UserID  uint      `json:"user_id"`
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}
