package task

import (
	"errors"

	"todo-app/domain/task/entities"
	"todo-app/domain/task/repositories"
	"todo-app/domain/task/services"
	"todo-app/domain/task/valueobjects"
	uservo "todo-app/domain/user/valueobjects"
)

// CreateTaskCommand represents a command to create a new task
type CreateTaskCommand struct {
	Title       string
	Description string
	Priority    string
	UserID      uint
}

// UpdateTaskCommand represents a command to update an existing task
type UpdateTaskCommand struct {
	TaskID      uint
	Title       *string
	Description *string
	Status      *string
	Priority    *string
	UserID      uint
}

// TaskQuery represents a query for tasks
type TaskQuery struct {
	UserID   uint
	Status   *string
	Priority *string
}

// TaskApplicationService orchestrates task-related use cases
type TaskApplicationService interface {
	// CreateTask creates a new task
	CreateTask(cmd CreateTaskCommand) (*entities.Task, error)

	// UpdateTask updates an existing task
	UpdateTask(cmd UpdateTaskCommand) (*entities.Task, error)

	// GetTask retrieves a specific task
	GetTask(taskID uint, userID uint) (*entities.Task, error)

	// GetUserTasks retrieves tasks for a user with optional filtering
	GetUserTasks(query TaskQuery) ([]*entities.Task, error)

	// DeleteTask deletes a task
	DeleteTask(taskID uint, userID uint) error

	// CompleteTask marks a task as completed
	CompleteTask(taskID uint, userID uint) (*entities.Task, error)

	// ArchiveTask archives a task
	ArchiveTask(taskID uint, userID uint) (*entities.Task, error)
}

// taskApplicationService implements TaskApplicationService
type taskApplicationService struct {
	taskRepo           repositories.TaskRepository
	validationService  services.TaskValidationService
	searchService      services.TaskSearchService
}

// NewTaskApplicationService creates a new task application service
func NewTaskApplicationService(
	taskRepo repositories.TaskRepository,
	validationService services.TaskValidationService,
	searchService services.TaskSearchService,
) TaskApplicationService {
	return &taskApplicationService{
		taskRepo:          taskRepo,
		validationService: validationService,
		searchService:     searchService,
	}
}

// CreateTask creates a new task with validation
func (s *taskApplicationService) CreateTask(cmd CreateTaskCommand) (*entities.Task, error) {
	// Create value objects
	title, err := valueobjects.NewTaskTitle(cmd.Title)
	if err != nil {
		return nil, err
	}

	description, err := valueobjects.NewTaskDescription(cmd.Description)
	if err != nil {
		return nil, err
	}

	priority, err := valueobjects.NewTaskPriority(cmd.Priority)
	if err != nil {
		return nil, err
	}

	userID, err := uservo.NewUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Validate task creation
	if err := s.validationService.ValidateTaskCreation(title, userID); err != nil {
		return nil, err
	}

	// Create pending status for new tasks
	status := valueobjects.NewPendingStatus()

	// Generate new task ID (in real implementation, this would come from repository)
	taskID, err := valueobjects.NewTaskID(0) // Repository will assign actual ID
	if err != nil {
		return nil, err
	}

	// Create the task entity
	task, err := entities.NewTask(taskID, title, description, status, priority, userID)
	if err != nil {
		return nil, err
	}

	// Save the task
	if err := s.taskRepo.Save(task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTask updates an existing task with validation
func (s *taskApplicationService) UpdateTask(cmd UpdateTaskCommand) (*entities.Task, error) {
	// Create task ID value object
	taskID, err := valueobjects.NewTaskID(cmd.TaskID)
	if err != nil {
		return nil, err
	}

	userID, err := uservo.NewUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Retrieve the existing task
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	// Check ownership
	if !task.IsOwnedBy(userID) {
		return nil, errors.New("access denied: task does not belong to user")
	}

	// Build updates for validation
	updates := services.TaskUpdates{}

	if cmd.Title != nil {
		title, err := valueobjects.NewTaskTitle(*cmd.Title)
		if err != nil {
			return nil, err
		}
		updates.Title = &title
	}

	if cmd.Description != nil {
		description, err := valueobjects.NewTaskDescription(*cmd.Description)
		if err != nil {
			return nil, err
		}
		updates.Description = &description
	}

	if cmd.Status != nil {
		status, err := valueobjects.NewTaskStatus(*cmd.Status)
		if err != nil {
			return nil, err
		}
		updates.Status = &status
	}

	if cmd.Priority != nil {
		priority, err := valueobjects.NewTaskPriority(*cmd.Priority)
		if err != nil {
			return nil, err
		}
		updates.Priority = &priority
	}

	// Validate the updates
	if err := s.validationService.ValidateTaskUpdate(task.Status(), updates); err != nil {
		return nil, err
	}

	// Apply the updates
	if updates.Title != nil {
		if err := task.UpdateTitle(*updates.Title); err != nil {
			return nil, err
		}
	}

	if updates.Description != nil {
		if err := task.UpdateDescription(*updates.Description); err != nil {
			return nil, err
		}
	}

	if updates.Status != nil {
		if updates.Status.IsCompleted() {
			if err := task.MarkAsCompleted(); err != nil {
				return nil, err
			}
		} else if updates.Status.IsArchived() {
			if err := task.Archive(); err != nil {
				return nil, err
			}
		}
	}

	if updates.Priority != nil {
		if err := task.ChangePriority(*updates.Priority); err != nil {
			return nil, err
		}
	}

	// Save the updated task
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask retrieves a specific task with ownership validation
func (s *taskApplicationService) GetTask(taskID uint, userID uint) (*entities.Task, error) {
	taskIDVO, err := valueobjects.NewTaskID(taskID)
	if err != nil {
		return nil, err
	}

	userIDVO, err := uservo.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	task, err := s.taskRepo.FindByID(taskIDVO)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	// Check ownership
	if !task.IsOwnedBy(userIDVO) {
		return nil, errors.New("access denied: task does not belong to user")
	}

	return task, nil
}

// GetUserTasks retrieves tasks for a user with optional filtering
func (s *taskApplicationService) GetUserTasks(query TaskQuery) ([]*entities.Task, error) {
	userID, err := uservo.NewUserID(query.UserID)
	if err != nil {
		return nil, err
	}

	// If status filter is provided
	if query.Status != nil {
		status, err := valueobjects.NewTaskStatus(*query.Status)
		if err != nil {
			return nil, err
		}
		return s.searchService.FindTasksByStatus(userID, status)
	}

	// If priority filter is provided
	if query.Priority != nil {
		priority, err := valueobjects.NewTaskPriority(*query.Priority)
		if err != nil {
			return nil, err
		}
		return s.searchService.FindTasksByPriority(userID, priority)
	}

	// No filters, return all tasks for user
	return s.taskRepo.FindByUserID(userID)
}

// DeleteTask deletes a task with ownership validation
func (s *taskApplicationService) DeleteTask(taskID uint, userID uint) error {
	taskIDVO, err := valueobjects.NewTaskID(taskID)
	if err != nil {
		return err
	}

	userIDVO, err := uservo.NewUserID(userID)
	if err != nil {
		return err
	}

	// Retrieve task to check ownership
	task, err := s.taskRepo.FindByID(taskIDVO)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("task not found")
	}

	// Check ownership
	if !task.IsOwnedBy(userIDVO) {
		return errors.New("access denied: task does not belong to user")
	}

	// Delete the task
	return s.taskRepo.Delete(taskIDVO)
}

// CompleteTask marks a task as completed
func (s *taskApplicationService) CompleteTask(taskID uint, userID uint) (*entities.Task, error) {
	cmd := UpdateTaskCommand{
		TaskID: taskID,
		Status: func() *string { s := "completed"; return &s }(),
		UserID: userID,
	}
	return s.UpdateTask(cmd)
}

// ArchiveTask archives a task
func (s *taskApplicationService) ArchiveTask(taskID uint, userID uint) (*entities.Task, error) {
	cmd := UpdateTaskCommand{
		TaskID: taskID,
		Status: func() *string { s := "archived"; return &s }(),
		UserID: userID,
	}
	return s.UpdateTask(cmd)
}