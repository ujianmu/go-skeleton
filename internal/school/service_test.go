package school

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/ujianmu/go-skeleton/internal/entity"
	"github.com/ujianmu/go-skeleton/pkg/log"
	"testing"
)

var errCRUD = errors.New("error crud")

func TestCreateSchoolRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateSchoolRequest
		wantError bool
	}{
		{"success", CreateSchoolRequest{Name: "test"}, false},
		{"required", CreateSchoolRequest{Name: ""}, true},
		{"too long", CreateSchoolRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateSchoolRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateSchoolRequest
		wantError bool
	}{
		{"success", UpdateSchoolRequest{Name: "test"}, false},
		{"required", UpdateSchoolRequest{Name: ""}, true},
		{"too long", UpdateSchoolRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	school, err := s.Create(ctx, CreateSchoolRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, school.ID)
	id := school.ID
	assert.Equal(t, "test", school.Name)
	assert.NotEmpty(t, school.CreatedAt)
	assert.NotEmpty(t, school.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateSchoolRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateSchoolRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateSchoolRequest{Name: "test2"})

	// update
	school, err = s.Update(ctx, id, UpdateSchoolRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", school.Name)
	_, err = s.Update(ctx, "none", UpdateSchoolRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateSchoolRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateSchoolRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	school, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", school.Name)
	assert.Equal(t, id, school.ID)

	// query
	schools, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(schools))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	school, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, school.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.School
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.School, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.School{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.School, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, school entity.School) error {
	if school.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, school)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, school entity.School) error {
	if school.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == school.ID {
			m.items[i] = school
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
