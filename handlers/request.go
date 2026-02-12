package handler

import "time"

//============用户认证请求结构体=============
type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserInfo struct {
	Username string `json:"username"`
	Telenum  string `json:"telenum"`
	Gender   string `json:"gender" binding:"omitempty,oneof=男 女 草履虫"`
}

type UpdateEmail struct {
	NewEmail string `json:"new_email"`
	Password string `json:"password"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CancelUser struct {
	Password string `json:"password"`
}

//============待办事项请求结构体=============
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline" binding:"required"` // DDL为必填
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTodoStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed"`
}
