package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"domain/user/entities"
	"domain/user/repositories"
	"domain/user/valueobjects"
	"todo-app/application/mappers"
	"todo-app/infrastructure/persistence"
	"todo-app/internal/dtos"
)

func setupUserRepositoryTest(t *testing.T) (*gorm.DB, repositories.UserRepository) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the User table
	err = db.AutoMigrate(&dtos.User{})
	require.NoError(t, err)

	// Create mapper and repository
	mapper := &mappers.UserMapper{}
	repo := persistence.NewGormUserRepository(db, mapper)

	return db, repo
}

func TestGormUserRepository_Save_ReturnsEntity(t *testing.T) {
	_, repo := setupUserRepositoryTest(t)

	// Create a valid user entity
	email, err := valueobjects.NewEmail("test@example.com")
	require.NoError(t, err)

	profile := valueobjects.NewUserProfile("John Doe")
	preferences := valueobjects.NewDefaultUserPreferences()

	user, err := entities.NewUser(
		valueobjects.NewUserID(1),
		email,
		profile,
		preferences,
	)
	require.NoError(t, err)

	// Save the entity
	err = repo.Save(user)
	require.NoError(t, err)

	// Verify that it was saved by retrieving it
	savedUser, err := repo.FindByID(valueobjects.NewUserID(1))
	require.NoError(t, err)
	require.NotNil(t, savedUser)

	// Verify it's an entity (not DTO)
	assert.Equal(t, uint(1), savedUser.ID().Value())
	assert.Equal(t, "test@example.com", savedUser.Email().Value())
	assert.Equal(t, "John Doe", savedUser.Profile().DisplayName())
}

func TestGormUserRepository_FindByID_ReturnsEntity(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert a DTO directly into database
	dto := &dtos.User{
		ID:        1,
		Email:     "test@example.com",
		Name:      "Jane Smith",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Retrieve using repository
	user, err := repo.FindByID(valueobjects.NewUserID(1))
	require.NoError(t, err)
	require.NotNil(t, user)

	// Verify it's an entity (not DTO)
	assert.Equal(t, uint(1), user.ID().Value())
	assert.Equal(t, "test@example.com", user.Email().Value())
	assert.Equal(t, "Jane Smith", user.Profile().DisplayName())
	assert.NotNil(t, user.Preferences())
}

func TestGormUserRepository_FindByEmail_ReturnsEntity(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert a DTO directly into database
	dto := &dtos.User{
		ID:        1,
		Email:     "email@test.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Retrieve by email
	email, err := valueobjects.NewEmail("email@test.com")
	require.NoError(t, err)

	user, err := repo.FindByEmail(email)
	require.NoError(t, err)
	require.NotNil(t, user)

	// Verify it's an entity
	assert.Equal(t, "email@test.com", user.Email().Value())
	assert.Equal(t, "Test User", user.Profile().DisplayName())
}

func TestGormUserRepository_Update_PersistsChanges(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert initial DTO
	dto := &dtos.User{
		ID:        1,
		Email:     "old@example.com",
		Name:      "Old Name",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Retrieve entity
	user, err := repo.FindByID(valueobjects.NewUserID(1))
	require.NoError(t, err)

	// Update the entity
	newEmail, err := valueobjects.NewEmail("new@example.com")
	require.NoError(t, err)
	err = user.ChangeEmail(newEmail)
	require.NoError(t, err)

	newProfile := valueobjects.NewUserProfile("New Name")
	err = user.UpdateProfile(newProfile)
	require.NoError(t, err)

	// Save changes
	err = repo.Update(user)
	require.NoError(t, err)

	// Retrieve again and verify changes
	updatedUser, err := repo.FindByID(valueobjects.NewUserID(1))
	require.NoError(t, err)
	assert.Equal(t, "new@example.com", updatedUser.Email().Value())
	assert.Equal(t, "New Name", updatedUser.Profile().DisplayName())
}

func TestGormUserRepository_Delete_RemovesUser(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert DTO
	dto := &dtos.User{
		ID:        1,
		Email:     "delete@example.com",
		Name:      "To Delete",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Delete
	err = repo.Delete(valueobjects.NewUserID(1))
	require.NoError(t, err)

	// Verify deleted
	user, err := repo.FindByID(valueobjects.NewUserID(1))
	require.NoError(t, err)
	assert.Nil(t, user)
}

func TestGormUserRepository_ExistsByEmail_ReturnsTrue(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert DTO
	dto := &dtos.User{
		ID:        1,
		Email:     "exists@example.com",
		Name:      "Exists",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Check existence
	email, err := valueobjects.NewEmail("exists@example.com")
	require.NoError(t, err)

	exists, err := repo.ExistsByEmail(email)
	require.NoError(t, err)
	assert.True(t, exists)

	// Check non-existent
	nonExistentEmail, err := valueobjects.NewEmail("notexists@example.com")
	require.NoError(t, err)

	exists, err = repo.ExistsByEmail(nonExistentEmail)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestGormUserRepository_FindAll_ReturnsEntities(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Insert multiple DTOs
	users := []dtos.User{
		{ID: 1, Email: "user1@example.com", Name: "User 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Email: "user2@example.com", Name: "User 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Email: "user3@example.com", Name: "User 3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, u := range users {
		err := db.Create(&u).Error
		require.NoError(t, err)
	}

	// Retrieve all
	allUsers, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, allUsers, 3)

	// Verify they're entities
	for _, user := range allUsers {
		assert.NotNil(t, user.ID())
		assert.NotNil(t, user.Email())
		assert.NotNil(t, user.Profile())
	}
}

func TestGormUserRepository_Count_ReturnsCorrectCount(t *testing.T) {
	db, repo := setupUserRepositoryTest(t)

	// Initially empty
	count, err := repo.Count()
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Insert DTOs
	for i := 1; i <= 5; i++ {
		dto := &dtos.User{
			ID:        uint(i),
			Email:     "user" + string(rune(i+'0')) + "@example.com",
			Name:      "User " + string(rune(i+'0')),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := db.Create(dto).Error
		require.NoError(t, err)
	}

	// Count
	count, err = repo.Count()
	require.NoError(t, err)
	assert.Equal(t, int64(5), count)
}
