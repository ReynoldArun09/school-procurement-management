package repositories

import (
	"context"
	"fmt"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type UnitRepository interface {
	Create(ctx context.Context, input *models.Unit) error
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, input *models.Unit) error
}

type GormUnitRepository struct {
	db *gorm.DB
}

func NewGormUnitRepository(db *gorm.DB) *GormUnitRepository {
	return &GormUnitRepository{
		db: db,
	}
}

func (r *GormUnitRepository) Create(ctx context.Context, input *models.Unit) error {
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return fmt.Errorf("failed to Create Brand: %w", err)
	}
	return nil
}

func (r *GormUnitRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Unit{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}

func (r *GormUnitRepository) Update(ctx context.Context, input *models.Unit) error {
	result := r.db.WithContext(ctx).Model(&models.Unit{}).Where("id = ?", input.ID).Select("name", "slug", "code", "status").Updates(input)

	if result.Error != nil {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete: %w", result.Error)
	}
	return nil
}
