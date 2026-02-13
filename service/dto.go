package service

import "time"

// 用户信息dto
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

// 待办事项dto
type TodoInfo struct {
	ID    uint   `json:"id"`
	Event string `json:"event"`
}

// 每日学习数据dto
type DailyStudyDataInfo struct {
	Date      time.Time `json:"date"`
	StudyTime int       `json:"study_time"`
	Tomatoes  int       `json:"tomatoes"`
}

// 每月学习数据dto
type MonthlyStudyDataInfo struct {
	Month     time.Time `json:"month"`
	StudyTime int       `json:"study_time"`
	Tomatoes  int       `json:"tomatoes"`
}

// 总学习数据dto
type TotalStudyDataInfo struct {
	StudyTime int `json:"study_time"`
	Tomatoes  int `json:"tomatoes"`
}
