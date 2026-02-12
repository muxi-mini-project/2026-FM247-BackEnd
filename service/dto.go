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

type TodoInfo struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline"`
	CompletedAt *time.Time `json:"completed_at"`
}
