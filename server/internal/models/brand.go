package models

import "time"

type Brand struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"not null" json:"name" validate:"min=4,max=100"`
	Slug      *string    `gorm:"uniqueIndex" json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
	Status    bool       `gorm:"not null; default:true" json:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (Brand) TableName() string {
	return "brands"
}
