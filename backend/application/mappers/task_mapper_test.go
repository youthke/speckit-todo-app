package mappers

import (
	"strings"
	"testing"
	"time"

	"todo-app/internal/dtos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskMapper_ToEntity_ValidDTO(t *testing.T) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "Test task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, err := mapper.ToEntity(dto)

	require.NoError(t, err)
	require.NotNil(t, entity)
	assert.Equal(t, uint(1), entity.ID().Value())
	assert.Equal(t, "Test task", entity.Title().Value())
	assert.True(t, entity.Status().IsPending())
	assert.Equal(t, uint(1), entity.UserID().Value())
}

func TestTaskMapper_ToEntity_CompletedTask(t *testing.T) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "Completed task",
		Completed: true,
		UserID:    1,
	}

	entity, err := mapper.ToEntity(dto)

	require.NoError(t, err)
	require.NotNil(t, entity)
	assert.True(t, entity.Status().IsCompleted())
}

func TestTaskMapper_ToEntity_EmptyTitle(t *testing.T) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "",
		Completed: false,
		UserID:    1,
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
	assert.Contains(t, err.Error(), "title")
}

func TestTaskMapper_ToEntity_LongTitle(t *testing.T) {
	mapper := &TaskMapper{}

	// Create a title longer than 500 characters
	longTitle := strings.Repeat("a", 501)

	dto := &dtos.Task{
		ID:        1,
		Title:     longTitle,
		Completed: false,
		UserID:    1,
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
	assert.Contains(t, err.Error(), "title")
}

func TestTaskMapper_ToEntity_ZeroTaskID(t *testing.T) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        0,
		Title:     "Test task",
		Completed: false,
		UserID:    1,
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
}

func TestTaskMapper_ToEntity_ZeroUserID(t *testing.T) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "Test task",
		Completed: false,
		UserID:    0, // Zero userID should fail
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
}

func TestTaskMapper_ToDTO_ValidEntity(t *testing.T) {
	mapper := &TaskMapper{}

	// First create a valid entity via ToEntity
	dto := &dtos.Task{
		ID:        1,
		Title:     "Test task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, err := mapper.ToEntity(dto)
	require.NoError(t, err)

	// Now convert back to DTO
	resultDTO := mapper.ToDTO(entity)

	require.NotNil(t, resultDTO)
	assert.Equal(t, uint(1), resultDTO.ID)
	assert.Equal(t, "Test task", resultDTO.Title)
	assert.False(t, resultDTO.Completed)
}

func TestTaskMapper_ToDTO_StatusConversion(t *testing.T) {
	mapper := &TaskMapper{}

	tests := []struct {
		name           string
		completed      bool
		expectedStatus bool
	}{
		{"Pending to false", false, false},
		{"Completed to true", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &dtos.Task{
				ID:        1,
				Title:     "Test task",
				Completed: tt.completed,
				UserID:    1,
			}

			entity, err := mapper.ToEntity(dto)
			require.NoError(t, err)

			resultDTO := mapper.ToDTO(entity)
			assert.Equal(t, tt.expectedStatus, resultDTO.Completed)
		})
	}
}

func TestTaskMapper_ToEntity_ToDTO_Roundtrip(t *testing.T) {
	mapper := &TaskMapper{}

	originalDTO := &dtos.Task{
		ID:        1,
		Title:     "Roundtrip test task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// DTO → Entity
	entity, err := mapper.ToEntity(originalDTO)
	require.NoError(t, err)

	// Entity → DTO
	resultDTO := mapper.ToDTO(entity)

	// Core fields should match
	assert.Equal(t, originalDTO.ID, resultDTO.ID)
	assert.Equal(t, originalDTO.Title, resultDTO.Title)
	assert.Equal(t, originalDTO.Completed, resultDTO.Completed)
}

// Benchmarks

func BenchmarkTaskMapper_ToEntity(b *testing.B) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "Benchmark test task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mapper.ToEntity(dto)
	}
}

func BenchmarkTaskMapper_ToDTO(b *testing.B) {
	mapper := &TaskMapper{}

	dto := &dtos.Task{
		ID:        1,
		Title:     "Benchmark test task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, _ := mapper.ToEntity(dto)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.ToDTO(entity)
	}
}
