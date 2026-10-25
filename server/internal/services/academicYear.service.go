package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ReynoldArun09/school-procurement-management/internal/models"
	"github.com/ReynoldArun09/school-procurement-management/internal/repositories"
)

var (
	ErrInvalidInput  = errors.New("invalid academic year input")
	ErrCurrentDelete = errors.New("cannot delete the current academic year")
)

type Service struct {
	repo repositories.AcademicYearRepository
}

func NewService(repo repositories.AcademicYearRepository) *Service {
	return &Service{repo: repo}
}

type CreateInput struct {
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsCurrent bool      `json:"is_current"`
	Status    *bool     `json:"status"`
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*models.AcademicYear, error) {
	if err := validateInput(input.Name, input.StartDate, input.EndDate); err != nil {
		return nil, err
	}

	status := true
	if input.Status != nil {
		status = *input.Status
	}

	Newyear := &models.AcademicYear{
		Name:      strings.TrimSpace(input.Name),
		StartDate: &input.StartDate,
		EndDate:   &input.EndDate,
		IsCurrent: false,
		Status:    status,
	}

	if err := s.repo.Create(ctx, Newyear); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fmt.Errorf("%w: Academic year already exists", ErrInvalidInput)
		}
		return nil, err
	}

	return Newyear, nil

}

func validateInput(name string, startDate, endDate time.Time) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if len([]rune(strings.TrimSpace(name))) > 255 {
		return fmt.Errorf("%w: name must not exceed 255 characters", ErrInvalidInput)
	}
	if startDate.IsZero() {
		return fmt.Errorf("%w: start_date is required", ErrInvalidInput)
	}
	if endDate.IsZero() {
		return fmt.Errorf("%w: end_date is required", ErrInvalidInput)
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("%w: end_date must be on or after start_date", ErrInvalidInput)
	}
	return nil
}
