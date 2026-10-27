package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type BrandRepository interface {
	DeleteBrand(ctx context.Context, id uint) error
	UpdateBrand(ctx context.Context, input *models.Brand) error
	CreateBrand(ctx context.Context, input *models.Brand) error
}

var BrandErrNotFound = errors.New("academic year not found")

type GormBrandRepository struct {
	db *gorm.DB
}

func NewGormBrandRepository(db *gorm.DB) *GormBrandRepository {
	return &GormBrandRepository{
		db: db,
	}
}

func (r *GormAcademicYearRepository) CreateBrand(ctx context.Context, input *models.Brand) error {
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return fmt.Errorf("failed to Create Brand: %w", err)
	}
	return nil
}

func (r *GormAcademicYearRepository) DeleteBrand(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Brand{}, id)

	if result.Error != nil {
		return fmt.Errorf("Failed to delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return BrandErrNotFound
	}
	return nil
}

func (r *GormAcademicYearRepository) UpdateBrand(ctx context.Context, input *models.Brand) error {
	result := r.db.WithContext(ctx).Model(&models.Brand{}).Where("id =?", input.ID).Select("name", "slug", "status").Updates(input)

	if result.Error != nil {
		return fmt.Errorf("delete academic year: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
