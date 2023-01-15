package school

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ujianmu/go-skeleton/internal/entity"
	"github.com/ujianmu/go-skeleton/pkg/log"
	"time"
)

// Service encapsulates usecase logic for School.
type Service interface {
	Get(ctx context.Context, id string) (School, error)
	Query(ctx context.Context, offset, limit int) ([]School, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateSchoolRequest) (School, error)
	Update(ctx context.Context, id string, input UpdateSchoolRequest) (School, error)
	Delete(ctx context.Context, id string) (School, error)
}

// School represents the data about an School.
type School struct {
	entity.School
}

// CreateSchoolRequest represents an school creation request.
type CreateSchoolRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateSchoolRequest fields.
func (m CreateSchoolRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateSchoolRequest represents an School update request.
type UpdateSchoolRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateSchoolRequest fields.
func (m UpdateSchoolRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new School service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the School with the specified the School ID.
func (s service) Get(ctx context.Context, id string) (School, error) {
	school, err := s.repo.Get(ctx, id)
	if err != nil {
		return School{}, err
	}
	return School{school}, nil
}

// Create creates a new school.
func (s service) Create(ctx context.Context, req CreateSchoolRequest) (School, error) {
	if err := req.Validate(); err != nil {
		return School{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.School{
		ID:        id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return School{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the school with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateSchoolRequest) (School, error) {
	if err := req.Validate(); err != nil {
		return School{}, err
	}

	school, err := s.Get(ctx, id)
	if err != nil {
		return school, err
	}
	school.Name = req.Name
	school.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, school.School); err != nil {
		return school, err
	}
	return school, nil
}

// Delete deletes the school with the specified ID.
func (s service) Delete(ctx context.Context, id string) (School, error) {
	school, err := s.Get(ctx, id)
	if err != nil {
		return School{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return School{}, err
	}
	return school, nil
}

// Count returns the number of school.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the school with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]School, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []School{}
	for _, item := range items {
		result = append(result, School{item})
	}
	return result, nil
}
