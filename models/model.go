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
	UserID    uint           `json:"user_id" gorm:"uniqueIndex:idx_user_date"` //联合索引,保证每天一条记录
	Date      time.Time      `json:"date" gorm:"uniqueIndex:idx_user_date"`
	StudyTime int            `json:"study_time"`
	Tomatoes  int            `json:"tomatoes"`
}

type MonthlyStudyData struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex:idx_user_month"` //联合索引,保证每月一条记录
	Month     time.Time      `json:"month" gorm:"uniqueIndex:idx_user_month"`
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
	UserID      uint       `gorm:"not null;index" json:"user_id"`                    // 关联用户ID
	Title       string     `gorm:"type:varchar(255);not null" json:"title"`          // 标题
	Description string     `gorm:"type:text" json:"description"`                     // 描述
	Status      string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // 完成状态: pending, in_progress, completed
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间
	StartTime   *time.Time `gorm:"index" json:"start_time"`                          // 开始时间（可选）
	DDL         *time.Time `gorm:"index" json:"deadline"`                            // DDL（截止时间，可选）
	CompletedAt *time.Time `json:"completed_at"`                                     // 完成时间（可选）
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

// 番茄钟实体定义
// type Tomato struct {
// 	ID     string `json:"id"`      // 番茄钟ID
// 	UserID uint   `json:"user_id"` // 用户ID

// 	// 时间配置
// 	StudyDuration int `json:"study_duration"` // 学习时长（分钟）
// 	BreakDuration int `json:"break_duration"` // 休息时长（分钟）

// 	// 循环配置
// 	AutoContinue bool `json:"auto_continue"` // 是否自动继续下一个番茄钟
// 	TotalCycles  int  `json:"total_cycles"`  // 总循环次数
// 	CurrentCycle int  `json:"current_cycle"` // 当前循环次数

// 	// 状态信息
// 	Status       string `json:"status"`        // 当前状态(running, paused, completed)
// 	Phase        string `json:"phase"`         // 当前阶段(study, break)
// 	RemainingSec int    `json:"remaining_sec"` // 剩余秒数
// 	ElapsedSec   int    `json:"elapsed_sec"`   // 已过秒数
// 	PausedSec    int    `json:"paused_sec"`    // 暂停秒数
// 	TotalElapsed int    `json:"total_elapsed"` // 总学习分钟:(已过秒数-暂停秒数)/60

// 	// 时间戳
// 	StartedAt   time.Time  `json:"started_at"`
// 	CompletedAt *time.Time `json:"completed_at,omitempty"`
// }
