package handler

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
	Event string `json:"event" binding:"required"`
}

// UpdateTodoRequest 更新待办事项请求
type UpdateTodoRequest struct {
	Event string `json:"event"`
}

//============学习数据请求结构体=============
type AddStudyDataRequest struct {
	Date      string `json:"date" binding:"required"` // 日期字符串，格式为 "YYYY-MM-DD"
	StudyTime int    `json:"study_time" binding:"required"`
	Tomatoes  int    `json:"tomatoes" binding:"required"`
}
