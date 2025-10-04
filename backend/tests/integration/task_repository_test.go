package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"domain/task/entities"
	"domain/task/repositories"
	"domain/task/valueobjects"
	uservo "domain/user/valueobjects"
	"todo-app/application/mappers"
	"todo-app/infrastructure/persistence"
	"todo-app/internal/dtos"
)

func setupTaskRepositoryTest(t *testing.T) (*gorm.DB, repositories.TaskRepository) {
	// Create in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the Task table
	err = db.AutoMigrate(&dtos.Task{})
	require.NoError(t, err)

	// Create mapper and repository
	mapper := &mappers.TaskMapper{}
	repo := persistence.NewGormTaskRepository(db, mapper)

	return db, repo
}

func TestGormTaskRepository_Save_ReturnsEntity(t *testing.T) {
	_, repo := setupTaskRepositoryTest(t)

	// Create a valid task entity
	title := valueobjects.NewTaskTitle("Test Task")
	description := valueobjects.NewTaskDescription("Test description")
	status := valueobjects.NewPendingStatus()
	priority := valueobjects.NewMediumPriority()
	userID := uservo.NewUserID(1)

	task, err := entities.NewTask(
		valueobjects.NewTaskID(1),
		title,
		description,
		status,
		priority,
		userID,
	)
	require.NoError(t, err)

	// Save the entity
	err = repo.Save(task)
	require.NoError(t, err)

	// Verify that it was saved by retrieving it
	savedTask, err := repo.FindByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)
	require.NotNil(t, savedTask)

	// Verify it's an entity (not DTO)
	assert.Equal(t, uint(1), savedTask.ID().Value())
	assert.Equal(t, "Test Task", savedTask.Title().Value())
	assert.Equal(t, uint(1), savedTask.UserID().Value())
}

func TestGormTaskRepository_FindByID_ReturnsEntity(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert a DTO directly into database
	dto := &dtos.Task{
		ID:        1,
		Title:     "Sample Task",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Retrieve using repository
	task, err := repo.FindByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)
	require.NotNil(t, task)

	// Verify it's an entity (not DTO)
	assert.Equal(t, uint(1), task.ID().Value())
	assert.Equal(t, "Sample Task", task.Title().Value())
	assert.False(t, task.Status().IsCompleted())
	assert.Equal(t, uint(1), task.UserID().Value())
}

func TestGormTaskRepository_FindByUserID_ReturnsEntities(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert multiple DTOs for same user
	tasks := []dtos.Task{
		{ID: 1, Title: "Task 1", Completed: false, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Task 2", Completed: true, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Title: "Task 3", Completed: false, UserID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, task := range tasks {
		err := db.Create(&task).Error
		require.NoError(t, err)
	}

	// Retrieve tasks for user 1
	userID := uservo.NewUserID(1)
	userTasks, err := repo.FindByUserID(userID)
	require.NoError(t, err)
	assert.Len(t, userTasks, 2)

	// Verify they're entities
	for _, task := range userTasks {
		assert.Equal(t, uint(1), task.UserID().Value())
		assert.NotNil(t, task.Title())
		assert.NotNil(t, task.Status())
	}
}

func TestGormTaskRepository_FindByUserIDAndStatus_ReturnsFilteredEntities(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert tasks with different statuses
	tasks := []dtos.Task{
		{ID: 1, Title: "Pending Task", Completed: false, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Completed Task", Completed: true, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Title: "Another Pending", Completed: false, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, task := range tasks {
		err := db.Create(&task).Error
		require.NoError(t, err)
	}

	// Find pending tasks
	userID := uservo.NewUserID(1)
	pendingStatus := valueobjects.NewPendingStatus()
	pendingTasks, err := repo.FindByUserIDAndStatus(userID, pendingStatus)
	require.NoError(t, err)
	assert.Len(t, pendingTasks, 2)

	// Verify status
	for _, task := range pendingTasks {
		assert.False(t, task.Status().IsCompleted())
	}

	// Find completed tasks
	completedStatus := valueobjects.NewCompletedStatus()
	completedTasks, err := repo.FindByUserIDAndStatus(userID, completedStatus)
	require.NoError(t, err)
	assert.Len(t, completedTasks, 1)
	assert.True(t, completedTasks[0].Status().IsCompleted())
}

func TestGormTaskRepository_Update_PersistsChanges(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert initial DTO
	dto := &dtos.Task{
		ID:        1,
		Title:     "Old Title",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Retrieve entity
	task, err := repo.FindByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)

	// Update the entity
	newTitle := valueobjects.NewTaskTitle("Updated Title")
	err = task.UpdateTitle(newTitle)
	require.NoError(t, err)

	err = task.MarkAsCompleted()
	require.NoError(t, err)

	// Save changes
	err = repo.Update(task)
	require.NoError(t, err)

	// Retrieve again and verify changes
	updatedTask, err := repo.FindByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedTask.Title().Value())
	assert.True(t, updatedTask.Status().IsCompleted())
}

func TestGormTaskRepository_Delete_RemovesTask(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert DTO
	dto := &dtos.Task{
		ID:        1,
		Title:     "To Delete",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Delete
	err = repo.Delete(valueobjects.NewTaskID(1))
	require.NoError(t, err)

	// Verify deleted
	task, err := repo.FindByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)
	assert.Nil(t, task)
}

func TestGormTaskRepository_ExistsByID_ReturnsTrue(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert DTO
	dto := &dtos.Task{
		ID:        1,
		Title:     "Exists",
		Completed: false,
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.Create(dto).Error
	require.NoError(t, err)

	// Check existence
	exists, err := repo.ExistsByID(valueobjects.NewTaskID(1))
	require.NoError(t, err)
	assert.True(t, exists)

	// Check non-existent
	exists, err = repo.ExistsByID(valueobjects.NewTaskID(999))
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestGormTaskRepository_FindByUserIDAndPriority_FiltersCorrectly(t *testing.T) {
	db, repo := setupTaskRepositoryTest(t)

	// Insert tasks (all will have medium priority from mapper)
	tasks := []dtos.Task{
		{ID: 1, Title: "Task 1", Completed: false, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Task 2", Completed: false, UserID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Title: "Task 3", Completed: false, UserID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	for _, task := range tasks {
		err := db.Create(&task).Error
		require.NoError(t, err)
	}

	// Find medium priority tasks for user 1
	userID := uservo.NewUserID(1)
	mediumPriority := valueobjects.NewMediumPriority()
	priorityTasks, err := repo.FindByUserIDAndPriority(userID, mediumPriority)
	require.NoError(t, err)

	// All tasks from DTO have medium priority by default from mapper
	assert.Len(t, priorityTasks, 2)

	for _, task := range priorityTasks {
		assert.Equal(t, "medium", task.Priority().Value())
	}
}
