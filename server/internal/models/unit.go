package models

import "time"

type Unit struct {
	ID        uint       `gorm:"primaryKey,autoIncrement" json:"id"`
	Name      string     `gorm:"not null" json:"name" validate:"required,min=4,max=100"`
	Code      *string    `json:"code,omitempty" validate:"omitempty,max=50"`
	Status    bool       `gorm:"not null, default:true" json:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (Unit) TableName() string {
	return "units"
}
