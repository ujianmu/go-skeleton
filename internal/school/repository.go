package school

import (
	"context"
	"github.com/ujianmu/go-skeleton/internal/entity"
	"github.com/ujianmu/go-skeleton/pkg/dbcontext"
	"github.com/ujianmu/go-skeleton/pkg/log"
)

// Repository encapsulates the logic to access school from the data source.
type Repository interface {
	// Get returns the school with the specified school ID.
	Get(ctx context.Context, id string) (entity.School, error)
	// Count returns the number of school.
	Count(ctx context.Context) (int, error)
	// Query returns the list of school with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.School, error)
	// Create saves a new school in the storage.
	Create(ctx context.Context, school entity.School) error
	// Update updates the school with given ID in the storage.
	Update(ctx context.Context, school entity.School) error
	// Delete removes the school with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists school in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new school repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the school with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.School, error) {
	var school entity.School
	err := r.db.With(ctx).Select().Model(id, &school)
	return school, err
}

// Create saves a new school record in the database.
// It returns the ID of the newly inserted school record.
func (r repository) Create(ctx context.Context, school entity.School) error {
	return r.db.With(ctx).Model(&school).Insert()
}

// Update saves the changes to an school in the database.
func (r repository) Update(ctx context.Context, school entity.School) error {
	return r.db.With(ctx).Model(&school).Update()
}

// Delete deletes an school with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	school, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&school).Delete()
}

// Count returns the number of the school records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("school").Row(&count)
	return count, err
}

// Query retrieves the school records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.School, error) {
	var schools []entity.School
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&schools)
	return schools, err
}
