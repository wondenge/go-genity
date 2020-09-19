package genity

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wondenge/go-genity/internal/entity"
	"github.com/wondenge/go-genity/pkg/log"
)

var errCRUD = errors.New("error crud")

func TestCreateGenityRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateGenityRequest
		wantError bool
	}{
		{"success", CreateGenityRequest{Name: "test"}, false},
		{"required", CreateGenityRequest{Name: ""}, true},
		{"too long", CreateGenityRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateGenityRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateGenityRequest
		wantError bool
	}{
		{"success", UpdateGenityRequest{Name: "test"}, false},
		{"required", UpdateGenityRequest{Name: ""}, true},
		{"too long", UpdateGenityRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	genity, err := s.Create(ctx, CreateGenityRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, genity.ID)
	id := genity.ID
	assert.Equal(t, "test", genity.Name)
	assert.NotEmpty(t, genity.CreatedAt)
	assert.NotEmpty(t, genity.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateGenityRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateGenityRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateGenityRequest{Name: "test2"})

	// update
	genity, err = s.Update(ctx, id, UpdateGenityRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", genity.Name)
	_, err = s.Update(ctx, "none", UpdateGenityRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateGenityRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateGenityRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	genity, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", genity.Name)
	assert.Equal(t, id, genity.ID)

	// query
	genitys, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(genitys))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	genity, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, genity.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Genity
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Genity, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Genity{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.Genity, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, genity entity.Genity) error {
	if genity.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, genity)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, genity entity.Genity) error {
	if genity.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == genity.ID {
			m.items[i] = genity
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
