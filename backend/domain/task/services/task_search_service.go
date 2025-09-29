package services

import (
	"todo-app/domain/task/entities"
	"todo-app/domain/task/repositories"
	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
)

// TaskSearchService provides domain search logic for tasks
type TaskSearchService interface {
	// FindTasksByStatus retrieves tasks by user and status
	FindTasksByStatus(userID uservo.UserID, status valueobjects.TaskStatus) ([]*entities.Task, error)

	// FindTasksByPriority retrieves tasks by user and priority
	FindTasksByPriority(userID uservo.UserID, priority valueobjects.TaskPriority) ([]*entities.Task, error)

	// FindActiveTasksForUser retrieves all non-archived tasks for a user
	FindActiveTasksForUser(userID uservo.UserID) ([]*entities.Task, error)

	// FindCompletedTasksForUser retrieves completed tasks for a user
	FindCompletedTasksForUser(userID uservo.UserID) ([]*entities.Task, error)
}

// taskSearchService implements TaskSearchService
type taskSearchService struct {
	taskRepo repositories.TaskRepository
}

// NewTaskSearchService creates a new task search service
func NewTaskSearchService(taskRepo repositories.TaskRepository) TaskSearchService {
	return &taskSearchService{
		taskRepo: taskRepo,
	}
}

// FindTasksByStatus retrieves tasks filtered by status for a specific user
func (s *taskSearchService) FindTasksByStatus(userID uservo.UserID, status valueobjects.TaskStatus) ([]*entities.Task, error) {
	return s.taskRepo.FindByUserIDAndStatus(userID, status)
}

// FindTasksByPriority retrieves tasks filtered by priority for a specific user
func (s *taskSearchService) FindTasksByPriority(userID uservo.UserID, priority valueobjects.TaskPriority) ([]*entities.Task, error) {
	return s.taskRepo.FindByUserIDAndPriority(userID, priority)
}

// FindActiveTasksForUser retrieves all tasks that are not archived
func (s *taskSearchService) FindActiveTasksForUser(userID uservo.UserID) ([]*entities.Task, error) {
	allTasks, err := s.taskRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Filter out archived tasks
	var activeTasks []*entities.Task
	for _, task := range allTasks {
		if !task.Status().IsArchived() {
			activeTasks = append(activeTasks, task)
		}
	}

	return activeTasks, nil
}

// FindCompletedTasksForUser retrieves only completed tasks
func (s *taskSearchService) FindCompletedTasksForUser(userID uservo.UserID) ([]*entities.Task, error) {
	completedStatus := valueobjects.NewCompletedStatus()
	return s.taskRepo.FindByUserIDAndStatus(userID, completedStatus)
}