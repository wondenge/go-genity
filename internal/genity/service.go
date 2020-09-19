package genity

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/wondenge/go-genity/internal/entity"
	"github.com/wondenge/go-genity/pkg/log"
)

// Service encapsulates usecase logic for genitys.
type Service interface {
	Get(ctx context.Context, id string) (Genity, error)
	Query(ctx context.Context, offset, limit int) ([]Genity, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateGenityRequest) (Genity, error)
	Update(ctx context.Context, id string, input UpdateGenityRequest) (Genity, error)
	Delete(ctx context.Context, id string) (Genity, error)
}

// Genity represents the data about an genity.
type Genity struct {
	entity.Genity
}

// CreateGenityRequest represents an genity creation request.
type CreateGenityRequest struct {
	Title string `json:"Title"`
}

// Validate validates the CreateGenityRequest fields.
func (m CreateGenityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateGenityRequest represents an genity update request.
type UpdateGenityRequest struct {
	Title string `json:"Title"`
}

// Validate validates the CreateGenityRequest fields.
func (m UpdateGenityRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new genity service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the genity with the specified the genity ID.
func (s service) Get(ctx context.Context, id string) (Genity, error) {
	genity, err := s.repo.Get(ctx, id)
	if err != nil {
		return Genity{}, err
	}
	return Genity{genity}, nil
}

// Create creates a new genity.
func (s service) Create(ctx context.Context, req CreateGenityRequest) (Genity, error) {
	if err := req.Validate(); err != nil {
		return Genity{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.Genity{
		Title:     req.Title,
		ID:        id,
		Timestamp: now,
	})
	if err != nil {
		return Genity{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the genity with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateGenityRequest) (Genity, error) {
	if err := req.Validate(); err != nil {
		return Genity{}, err
	}

	genity, err := s.Get(ctx, id)
	if err != nil {
		return genity, err
	}
	genity.Title = req.Title
	genity.Timestamp = time.Now()

	if err := s.repo.Update(ctx, genity.Genity); err != nil {
		return genity, err
	}
	return genity, nil
}

// Delete deletes the genity with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Genity, error) {
	genity, err := s.Get(ctx, id)
	if err != nil {
		return Genity{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Genity{}, err
	}
	return genity, nil
}

// Count returns the number of genitys.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the genitys with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Genity, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Genity{}
	for _, item := range items {
		result = append(result, Genity{item})
	}
	return result, nil
}
