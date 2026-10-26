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

type UpdateInput struct {
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

func (s *Service) GetByID(ctx context.Context, id uint) (*models.AcademicYear, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: Id must be postivie", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: Id must be postivie", ErrInvalidInput)
	}
	year, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return err
	}
	if year.IsCurrent {
		return ErrCurrentDelete
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) SetCurrent(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: Id must be postivie", ErrInvalidInput)
	}

	return s.repo.SetCurrent(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint, input UpdateInput) (*models.AcademicYear, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: ID must be positive", ErrInvalidInput)
	}
	if err := validateInput(input.Name, input.StartDate, input.EndDate); err != nil {
		return nil, err
	}

	year, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	year.Name = input.Name
	year.StartDate = &input.StartDate
	year.EndDate = &input.EndDate

	if input.Status != nil {
		year.Status = *input.Status
	}

	year.IsCurrent = false
	if err := s.repo.Update(ctx, year); err != nil {
		return nil, err
	}

	if input.IsCurrent {
		if err := s.repo.SetCurrent(ctx, id); err != nil {
			return nil, err
		}
		year.IsCurrent = true
	}
	return year, nil
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
