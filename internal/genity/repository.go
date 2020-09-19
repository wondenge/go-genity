package genity

import (
	"context"

	"github.com/wondenge/go-genity/internal/entity"
	"github.com/wondenge/go-genity/pkg/dbcontext"
	"github.com/wondenge/go-genity/pkg/log"
)

// Repository encapsulates the logic to access genitys from the data source.
type Repository interface {
	// Get returns the genity with the specified genity ID.
	Get(ctx context.Context, id string) (entity.Genity, error)
	// Count returns the number of genitys.
	Count(ctx context.Context) (int, error)
	// Query returns the list of genitys with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Genity, error)
	// Create saves a new genity in the storage.
	Create(ctx context.Context, genity entity.Genity) error
	// Update updates the genity with given ID in the storage.
	Update(ctx context.Context, genity entity.Genity) error
	// Delete removes the genity with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists genitys in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new genity repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the genity with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Genity, error) {
	var genity entity.Genity
	err := r.db.With(ctx).Select().Model(id, &genity)
	return genity, err
}

// Create saves a new genity record in the database.
// It returns the ID of the newly inserted genity record.
func (r repository) Create(ctx context.Context, genity entity.Genity) error {
	return r.db.With(ctx).Model(&genity).Insert()
}

// Update saves the changes to an genity in the database.
func (r repository) Update(ctx context.Context, genity entity.Genity) error {
	return r.db.With(ctx).Model(&genity).Update()
}

// Delete deletes an genity with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	genity, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&genity).Delete()
}

// Count returns the number of the genity records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("genity").Row(&count)
	return count, err
}

// Query retrieves the genity records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Genity, error) {
	var genitys []entity.Genity
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&genitys)
	return genitys, err
}
