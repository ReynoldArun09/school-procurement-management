package models

import "time"

type Category struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"not null;uniqueIndex" json:"name" validate:"required,min=4,max=100"`
	Slug      string     `gorm:"not null;uniqueIndex" json:"slug" validate:"required,min=4,max-100"`
	Order     int        `gorm:"not null;default:0" json:"order" validate:"gte=0"`
	Status    bool       `gorm:"not null;default:true" json:"status"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (Category) TableName() string {
	return "categories"
}
