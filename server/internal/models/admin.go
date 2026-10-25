package models

import "time"

type Admin struct {
	ID            int       `gorm:"primaryKey;autoincrement" json:"id"`
	Name          string    `gorm:"not null" json:"name" validate:"required,min=2,max=100"`
	Username      string    `gorm:"not null;uniqueIndex" json:"username" validate:"required,min=3,max=100"`
	Email         string    `gorm:"not null;uniqueIndex" json:"email" validate:"required,email,max=255"`
	Password      string    `gorm:"not null" json:"-"`
	Image         *string   `json:"image,omitempty"`
	RememberToken *string   `json:"remember_token,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Admin) TableName() string {
	return "admin"
}
