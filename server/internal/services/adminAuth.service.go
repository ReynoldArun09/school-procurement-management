package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"github.com/ReynoldArun09/school-procurement-management/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveAdmin      = errors.New("admin account is inactive")
	ErrInvalidSession     = errors.New("invalid or expired session")
)

const sessionDuration = 8 * time.Hour

type AdminAuthService struct {
	adminAuthRepo repositories.AdminAuthRepository
	sessionRepo   repositories.AdminSessionRepository
}

type LoginResult struct {
	Admin     models.Admin
	SessionID string
	ExpiresAt time.Time
}

func NewAdminAuthService(adminAuthRepo repositories.AdminAuthRepository, sessionRepo repositories.AdminSessionRepository) *AdminAuthService {
	return &AdminAuthService{
		adminAuthRepo: adminAuthRepo,
		sessionRepo:   sessionRepo,
	}
}

func (s *AdminAuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	admin, err := s.adminAuthRepo.FindByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}

	sessionToken := hex.EncodeToString(tokenBytes)

	tokenHash := sha256.Sum256([]byte(sessionToken))

	hashedToken := hex.EncodeToString(tokenHash[:])

	expiresAt := time.Now().Add(sessionDuration)

	session := &models.Session{
		AdminID:   uint(admin.ID),
		TokenHash: hashedToken,
		ExpiresAt: expiresAt,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &LoginResult{
		Admin:     *admin,
		SessionID: sessionToken,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AdminAuthService) Authenticate(ctx context.Context, sessionToken string) (*models.Admin, error) {
	if sessionToken == "" {
		return nil, ErrInvalidCredentials
	}

	tokenHash := sha256.Sum256([]byte(sessionToken))
	hashedToken := hex.EncodeToString(tokenHash[:])

	session, err := s.sessionRepo.FindValidByTokenHash(ctx, hashedToken, time.Now())

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	admin, err := s.adminAuthRepo.FindByActiveID(ctx, session.AdminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	return admin, nil

}

func (s *AdminAuthService) Logout(ctx context.Context, sessionToken string) error {
	if sessionToken == "" {
		return nil
	}

	tokenHash := sha256.Sum256([]byte(sessionToken))
	hashedtoken := hex.EncodeToString(tokenHash[:])

	return s.sessionRepo.DeleteByTokenHash(ctx, hashedtoken)
}
