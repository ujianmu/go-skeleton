package school

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/ujianmu/go-skeleton/internal/entity"
	"github.com/ujianmu/go-skeleton/internal/test"
	"github.com/ujianmu/go-skeleton/pkg/log"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	test.ResetTables(t, db, "school")
	repo := NewRepository(db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	err = repo.Create(ctx, entity.School{
		ID:        "test1",
		Name:      "school1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	school, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "school1", school.Name)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, entity.School{
		ID:        "test1",
		Name:      "school1 updated",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	school, _ = repo.Get(ctx, "test1")
	assert.Equal(t, "school1 updated", school.Name)

	// query
	schools, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(schools))

	// delete
	err = repo.Delete(ctx, "test1")
	assert.Nil(t, err)
	_, err = repo.Get(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
}
