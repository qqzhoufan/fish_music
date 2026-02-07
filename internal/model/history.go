package model

import (
	"time"
)

// History 播放历史模型
type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_created" json:"user_id"`
	SongID    uint      `gorm:"not null;index" json:"song_id"`
	CreatedAt time.Time `gorm:"index:idx_user_created" json:"created_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Song *Song `gorm:"foreignKey:SongID" json:"song,omitempty"`
}

// TableName 指定表名
func (History) TableName() string {
	return "history"
}
