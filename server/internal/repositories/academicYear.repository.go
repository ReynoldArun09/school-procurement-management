package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AcademicYearRepository interface {
	Create(ctx context.Context, year *models.AcademicYear) error
	GetByID(ctx context.Context, id uint) (*models.AcademicYear, error)
	// List(ctx context.Context, search string, page, pageSize int) ([]models.AcademicYear, int64, error)
	Update(ctx context.Context, year *models.AcademicYear) error
	Delete(ctx context.Context, id uint) error
	SetCurrent(ctx context.Context, id uint) error
}

var ErrNotFound = errors.New("academic year not found")

type GormAcademicYearRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormAcademicYearRepository {
	return &GormAcademicYearRepository{db: db}
}

func (r *GormAcademicYearRepository) Create(ctx context.Context, year *models.AcademicYear) error {
	if err := r.db.WithContext(ctx).Create(year).Error; err != nil {
		return fmt.Errorf("Create academic year: %w", err)
	}
	return nil
}

func (r *GormAcademicYearRepository) GetByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
	var year models.AcademicYear

	err := r.db.WithContext(ctx).First(&year, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get academic year by ID: %w", err)
	}

	return &year, nil
}

func (r *GormAcademicYearRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.AcademicYear{}, id)

	if result.Error != nil {
		return fmt.Errorf("delete academic year: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *GormAcademicYearRepository) Update(ctx context.Context, year *models.AcademicYear) error {
	result := r.db.WithContext(ctx).Model(&models.AcademicYear{}).Where("id = ?", year.ID).Select("name", "start_date", "end_date", "status").Updates(year)

	if result.Error != nil {
		return fmt.Errorf("delete academic year: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *GormAcademicYearRepository) SetCurrent(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var year models.AcademicYear

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&year, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("find academic year to set academic year: %w", err)
		}

		if err := tx.Model(&models.AcademicYear{}).Where("is_current = ?", true).Update("is_current", false).Error; err != nil {
			return fmt.Errorf("clear current academic year: %w", err)
		}

		if err := tx.Model(&models.AcademicYear{}).Where("id = ?", id).Update("is_current", true).Error; err != nil {
			return fmt.Errorf("set current academic year: %w", err)
		}

		return nil
	})

}
