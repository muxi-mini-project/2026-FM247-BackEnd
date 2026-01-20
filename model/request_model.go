package model

// 这里是由handler层传入service层的结构体
type UpdateUserInfo struct {
	Username string `json:"username"`
	Telenum  string `json:"telenum"`
}

type RegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Telenum  string `json:"telenum"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateEmail struct {
	NewEmail string `json:"new_email"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CancelUser struct {
	Password string `json:"password"`
}
