package models

import "time"

type Session struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"not null;index" json:"admin_id"`
	TokenHash string    `gorm:"size:64;not null;uniqueIndex" json:"tokenHash"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	UpdatedAt time.Time `gorm:"not null;autoCreateTime" json:"updated_at"`
}

func (Session) TableName() string {
	return "sessions"
}
