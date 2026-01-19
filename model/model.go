package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `json:"id"`
	UserName  string `json:"username" gorm:"column:username"`
	Password  string `json:"-"`
	Telenum   string `json:"telenum"`
	Level     string `json:"level"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Todo struct {
	Id          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Note struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
