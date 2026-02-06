package service

import "time"

type UserInfo struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Gender     string    `json:"gender"`
	Telenum    string    `json:"telenum"`
	Avatar     string    `json:"avatar_path"`
	Experience int       `json:"experience"`
	Level      int       `json:"level"`
	CreatedAt  time.Time `json:"created_at"`
}
