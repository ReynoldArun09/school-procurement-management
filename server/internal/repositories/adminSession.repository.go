package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"gorm.io/gorm"
)

type AdminSessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	FindValidByTokenHash(ctx context.Context, tokenHash string, now time.Time) (*models.Session, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}

type GormSessionRepository struct {
	db *gorm.DB
}

func NewGormAdminSessionRepository(db *gorm.DB) *GormSessionRepository {
	return &GormSessionRepository{
		db: db,
	}
}

func (r *GormSessionRepository) Create(ctx context.Context, session *models.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *GormSessionRepository) FindValidByTokenHash(ctx context.Context, tokenHash string, now time.Time) (*models.Session, error) {
	var session models.Session

	err := r.db.WithContext(ctx).Where("tokenHash = ? AND expires_at > ?", tokenHash, now).First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &session, nil
}

func (r *GormSessionRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).Where("tokenHash = ?", tokenHash).Delete(&models.Session{}).Error
}
