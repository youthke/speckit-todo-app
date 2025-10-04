package mappers

import (
	"testing"
	"time"

	"todo-app/internal/dtos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserMapper_ToEntity_ValidDTO(t *testing.T) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, err := mapper.ToEntity(dto)

	require.NoError(t, err)
	require.NotNil(t, entity)
	assert.Equal(t, uint(1), entity.ID().Value())
	assert.Equal(t, "test@example.com", entity.Email().Value())
	assert.Equal(t, "John Doe", entity.Profile().DisplayName())
}

func TestUserMapper_ToEntity_InvalidEmail(t *testing.T) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:    1,
		Email: "invalid-email",
		Name:  "John Doe",
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
	assert.Contains(t, err.Error(), "email")
}

func TestUserMapper_ToEntity_EmptyEmail(t *testing.T) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:    1,
		Email: "",
		Name:  "John Doe",
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
	assert.Contains(t, err.Error(), "email")
}

func TestUserMapper_ToEntity_ZeroID(t *testing.T) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:    0,
		Email: "test@example.com",
		Name:  "John Doe",
	}

	entity, err := mapper.ToEntity(dto)

	assert.Error(t, err)
	assert.Nil(t, entity)
}

func TestUserMapper_ToDTO_ValidEntity(t *testing.T) {
	mapper := &UserMapper{}

	// First create a valid entity via ToEntity to ensure we have proper structure
	dto := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, err := mapper.ToEntity(dto)
	require.NoError(t, err)

	// Now convert back to DTO
	resultDTO := mapper.ToDTO(entity)

	require.NotNil(t, resultDTO)
	assert.Equal(t, uint(1), resultDTO.ID)
	assert.Equal(t, "test@example.com", resultDTO.Email)
	assert.Equal(t, "John Doe", resultDTO.Name)
	assert.Equal(t, "password", resultDTO.AuthMethod) // Default value
	assert.True(t, resultDTO.IsActive)                // Default value
}

func TestUserMapper_ToEntity_ToDTO_Roundtrip(t *testing.T) {
	mapper := &UserMapper{}

	originalDTO := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Jane Smith",
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
	assert.Equal(t, originalDTO.Email, resultDTO.Email)
	assert.Equal(t, originalDTO.Name, resultDTO.Name)
}

// Benchmarks

func BenchmarkUserMapper_ToEntity(b *testing.B) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mapper.ToEntity(dto)
	}
}

func BenchmarkUserMapper_ToDTO(b *testing.B) {
	mapper := &UserMapper{}

	dto := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entity, _ := mapper.ToEntity(dto)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mapper.ToDTO(entity)
	}
}
