package model

import (
	"fmt"
	"time"
)

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TelegramID int64    `gorm:"uniqueIndex;not null" json:"telegram_id"`  // Telegram 用户 ID
	Username  string    `gorm:"size:255" json:"username"`                  // Telegram 用户名
	FirstName string    `gorm:"size:255" json:"first_name"`                // 名字
	LastName  string    `gorm:"size:255" json:"last_name"`                 // 姓氏
	Language  string    `gorm:"size:10;default:zh" json:"language"`        // 语言设置
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`             // 是否管理员
	IsActive  bool      `gorm:"default:true" json:"is_active"`             // 是否激活
	LastSeen  time.Time `json:"last_seen"`                                 // 最后活跃时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Favorites []Favorite `gorm:"foreignKey:UserID" json:"favorites,omitempty"`
	History   []History   `gorm:"foreignKey:UserID" json:"history,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// GetFullName 获取用户全名
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.Username != "" {
		return "@" + u.Username
	}
	return fmt.Sprintf("User_%d", u.TelegramID)
}
