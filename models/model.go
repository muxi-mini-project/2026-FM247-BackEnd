package models

import (
	"time"

	"gorm.io/gorm"
)

type TotalStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type DailyStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	Date      time.Time      `json:"date"`
	EndAt     time.Time      `json:"end_at"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type MonthlyStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	Month     time.Time      `json:"month"`
	EndAt     time.Time      `json:"end_at"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type Note struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id"`
	Title     string         `json:"title"`
	Date      time.Time      `json:"date"`
	Content   string         `json:"content"`
}

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username   string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email      string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Telenum    string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"` // 密码哈希，不返回给前端
	Gender     string    `gorm:"type:varchar(10);default:'草履虫'" json:"gender"`
	Experience int       `gorm:"default:0" json:"experience"` // 经验值
	Level      int       `gorm:"default:1" json:"level"`      // 等级
	Avatar     string    `gorm:"type:varchar(500);default:'/default-avatar.png'" json:"avatar_path"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 扩展字段（根据需求添加）
	LastLoginAt *time.Time `json:"last_login_at"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	Settings    string     `gorm:"type:json" json:"settings"` // 用户设置，JSON格式存储
}

// Todo 待办事项表
type Todo struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"type:varchar(20);default:'未完成'" json:"status"` // pending, in_progress, completed, archived
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	StartTime   *time.Time `json:"start_time"`
	DDL         *time.Time `json:"ddl"`
	CompletedAt *time.Time `json:"completed_at"`
	User        User       `gorm:"foreignKey:UserID" json:"-"`
}

// Audio 音频表 - 对应图片1的"音频"表
type Audio struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	FilePath  string    `gorm:"type:varchar(500);not null" json:"file_path"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// UserAudioFavorite 用户音频收藏关联表
type UserAudioFavorite struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	AudioID   uint      `gorm:"primaryKey" json:"audio_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	User  User  `gorm:"foreignKey:UserID" json:"-"`
	Audio Audio `gorm:"foreignKey:AudioID" json:"-"`
}

type TokenBlacklist struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Jti       string         `json:"jti" gorm:"uniqueIndex"`
	ExpiresAt time.Time      `json:"expires_at"`
}
