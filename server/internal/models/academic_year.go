package models

import "time"

type AcademicYear struct {
	ID        int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"type:varchar(20);not null; uniqueIndex" json:"name" validate:"required,min=9,max=20"`
	StartDate *time.Time `gorm:"type:date" json:"start_date,omitempty"`
	EndDate   *time.Time `gorm:"type:date" json:"end_date,omitempty"`
	IsCurrent bool       `gorm:"not null;default:false" json:"is_current"`
	Status    bool       `gorm:"not null;default:true" json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// method for renaming the default table name.
func (AcademicYear) TableName() string {
	return "academic_years"
}
