package repositories

import (
	"context"
	"errors"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type AdminAuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Admin, error)
	FindByActiveID(ctx context.Context, id uint) (*models.Admin, error)
}

type GormAdminAuthRepository struct {
	db *gorm.DB
}

func NewGormAdminAuthRepository(db *gorm.DB) *GormAdminAuthRepository {
	return &GormAdminAuthRepository{
		db: db,
	}
}

func (r *GormAdminAuthRepository) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admins models.Admin

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&admins).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &admins, nil
}

func (r *GormAdminAuthRepository) FindByActiveID(ctx context.Context, id uint) (*models.Admin, error) {
	var admin models.Admin

	err := r.db.WithContext(ctx).Where("id = ? AND is_active = ?", id, true).First(&admin).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, err
	}

	return &admin, nil
}
