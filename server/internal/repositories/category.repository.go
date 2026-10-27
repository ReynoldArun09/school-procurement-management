package repositories

import (
	"context"
	"fmt"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, input *models.Category) error
	DeleteCategory(ctx context.Context, id uint) error
	UpdateCategory(ctx context.Context, input *models.Category) error
}

type GormCategoryRepository struct {
	db *gorm.DB
}

func NewGormCategoryRepository(db *gorm.DB) *GormCategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (r *GormCategoryRepository) CreateCategory(ctx context.Context, input *models.Category) error {
	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *GormCategoryRepository) DeleteCategory(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Category{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	return nil
}

func (r *GormCategoryRepository) UpdateCategory(ctx context.Context, input *models.Category) error {

	result := r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ?", input.ID).Select("name", "slug", "code", "status").Updates(input)

	if result.Error != nil {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}
	return nil
}
