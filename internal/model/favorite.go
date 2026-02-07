package model

import (
	"time"
)

// Favorite 收藏模型
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	SongID    uint      `gorm:"not null;index" json:"song_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Song *Song `gorm:"foreignKey:SongID" json:"song,omitempty"`
}

// TableName 指定表名
func (Favorite) TableName() string {
	return "favorites"
}
