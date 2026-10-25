package repositories

import (
	"context"
	"fmt"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type AcademicYearRepository interface {
	Create(ctx context.Context, year *models.AcademicYear) error
	// GetByID(ctx context.Context, id uint) (*models.AcademicYear, error)
	// List(ctx context.Context, search string, page, pageSize int) ([]models.AcademicYear, int64, error)
	// Update(ctx context.Context, year *models.AcademicYear) error
	// Delete(ctx context.Context, id uint) error
	// SetCurrent(ctx context.Context, id uint) error
}

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
