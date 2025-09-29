package user

import (
	"errors"

	"todo-app/domain/user/entities"
	"todo-app/domain/user/repositories"
	"todo-app/domain/user/services"
	"todo-app/domain/user/valueobjects"
	taskvo "todo-app/domain/task/valueobjects"
)

// RegisterUserCommand represents a command to register a new user
type RegisterUserCommand struct {
	Email     string
	FirstName string
	LastName  string
	Timezone  string
	// Optional preferences
	DefaultTaskPriority *string
	EmailNotifications  *bool
	ThemePreference     *string
}

// UpdateUserProfileCommand represents a command to update user profile
type UpdateUserProfileCommand struct {
	UserID    uint
	FirstName *string
	LastName  *string
	Timezone  *string
}

// UpdateUserPreferencesCommand represents a command to update user preferences
type UpdateUserPreferencesCommand struct {
	UserID              uint
	DefaultTaskPriority *string
	EmailNotifications  *bool
	ThemePreference     *string
}

// UserApplicationService orchestrates user-related use cases
type UserApplicationService interface {
	// RegisterUser registers a new user with validation
	RegisterUser(cmd RegisterUserCommand) (*entities.User, error)

	// GetUserProfile retrieves a user's profile
	GetUserProfile(userID uint) (*entities.User, error)

	// UpdateUserProfile updates user profile information
	UpdateUserProfile(cmd UpdateUserProfileCommand) (*entities.User, error)

	// GetUserPreferences retrieves user preferences
	GetUserPreferences(userID uint) (valueobjects.UserPreferences, error)

	// UpdateUserPreferences updates user preferences
	UpdateUserPreferences(cmd UpdateUserPreferencesCommand) (valueobjects.UserPreferences, error)

	// GetUserByEmail retrieves a user by email address
	GetUserByEmail(email string) (*entities.User, error)

	// ChangeUserEmail changes a user's email address
	ChangeUserEmail(userID uint, newEmail string) (*entities.User, error)
}

// userApplicationService implements UserApplicationService
type userApplicationService struct {
	userRepo         repositories.UserRepository
	authService      services.UserAuthenticationService
	profileService   services.UserProfileService
}

// NewUserApplicationService creates a new user application service
func NewUserApplicationService(
	userRepo repositories.UserRepository,
	authService services.UserAuthenticationService,
	profileService services.UserProfileService,
) UserApplicationService {
	return &userApplicationService{
		userRepo:       userRepo,
		authService:    authService,
		profileService: profileService,
	}
}

// RegisterUser registers a new user with complete validation
func (s *userApplicationService) RegisterUser(cmd RegisterUserCommand) (*entities.User, error) {
	// Create email value object
	email, err := valueobjects.NewEmail(cmd.Email)
	if err != nil {
		return nil, err
	}

	// Create profile value object
	profile, err := valueobjects.NewUserProfile(cmd.FirstName, cmd.LastName, cmd.Timezone)
	if err != nil {
		return nil, err
	}

	// Validate registration data using domain service
	if err := s.authService.ValidateRegistrationData(email, profile); err != nil {
		return nil, err
	}

	// Create preferences (use defaults or provided values)
	preferences, err := s.createUserPreferences(cmd)
	if err != nil {
		return nil, err
	}

	// Generate user ID (will be set by repository)
	userID, err := valueobjects.NewUserID(0) // Repository will assign actual ID
	if err != nil {
		return nil, err
	}

	// Create user entity
	user, err := entities.NewUser(userID, email, profile, preferences)
	if err != nil {
		return nil, err
	}

	// Save the user
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// createUserPreferences creates user preferences from command, using defaults for missing values
func (s *userApplicationService) createUserPreferences(cmd RegisterUserCommand) (valueobjects.UserPreferences, error) {
	// Set default task priority
	var defaultPriority taskvo.TaskPriority
	if cmd.DefaultTaskPriority != nil {
		priority, err := taskvo.NewTaskPriority(*cmd.DefaultTaskPriority)
		if err != nil {
			return valueobjects.UserPreferences{}, err
		}
		defaultPriority = priority
	} else {
		defaultPriority = taskvo.NewMediumPriority()
	}

	// Set email notifications (default true)
	emailNotifications := true
	if cmd.EmailNotifications != nil {
		emailNotifications = *cmd.EmailNotifications
	}

	// Set theme preference (default auto)
	themePreference := valueobjects.ThemeAuto
	if cmd.ThemePreference != nil {
		themePreference = *cmd.ThemePreference
	}

	return valueobjects.NewUserPreferences(defaultPriority, emailNotifications, themePreference)
}

// GetUserProfile retrieves a user's complete profile
func (s *userApplicationService) GetUserProfile(userID uint) (*entities.User, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userIDVO)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// UpdateUserProfile updates user profile information with validation
func (s *userApplicationService) UpdateUserProfile(cmd UpdateUserProfileCommand) (*entities.User, error) {
	userIDVO, err := valueobjects.NewUserID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Retrieve current user
	user, err := s.userRepo.FindByID(userIDVO)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Create partial update data
	updates := services.ProfileUpdateData{
		FirstName: cmd.FirstName,
		LastName:  cmd.LastName,
		Timezone:  cmd.Timezone,
	}

	// Use domain service to update profile
	if err := s.profileService.UpdatePartialProfile(userIDVO, updates); err != nil {
		return nil, err
	}

	// Return updated user
	return s.userRepo.FindByID(userIDVO)
}

// GetUserPreferences retrieves user preferences
func (s *userApplicationService) GetUserPreferences(userID uint) (valueobjects.UserPreferences, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return valueobjects.UserPreferences{}, err
	}

	user, err := s.userRepo.FindByID(userIDVO)
	if err != nil {
		return valueobjects.UserPreferences{}, err
	}

	if user == nil {
		return valueobjects.UserPreferences{}, errors.New("user not found")
	}

	return user.Preferences(), nil
}

// UpdateUserPreferences updates user preferences
func (s *userApplicationService) UpdateUserPreferences(cmd UpdateUserPreferencesCommand) (valueobjects.UserPreferences, error) {
	userIDVO, err := valueobjects.NewUserID(cmd.UserID)
	if err != nil {
		return valueobjects.UserPreferences{}, err
	}

	// Retrieve current user
	user, err := s.userRepo.FindByID(userIDVO)
	if err != nil {
		return valueobjects.UserPreferences{}, err
	}

	if user == nil {
		return valueobjects.UserPreferences{}, errors.New("user not found")
	}

	// Get current preferences
	currentPrefs := user.Preferences()

	// Build new preferences with updates
	var defaultPriority taskvo.TaskPriority
	if cmd.DefaultTaskPriority != nil {
		priority, err := taskvo.NewTaskPriority(*cmd.DefaultTaskPriority)
		if err != nil {
			return valueobjects.UserPreferences{}, err
		}
		defaultPriority = priority
	} else {
		defaultPriority = currentPrefs.DefaultTaskPriority()
	}

	emailNotifications := currentPrefs.EmailNotifications()
	if cmd.EmailNotifications != nil {
		emailNotifications = *cmd.EmailNotifications
	}

	themePreference := currentPrefs.ThemePreference()
	if cmd.ThemePreference != nil {
		themePreference = *cmd.ThemePreference
	}

	// Create new preferences
	newPrefs, err := valueobjects.NewUserPreferences(defaultPriority, emailNotifications, themePreference)
	if err != nil {
		return valueobjects.UserPreferences{}, err
	}

	// Update user
	if err := user.UpdatePreferences(newPrefs); err != nil {
		return valueobjects.UserPreferences{}, err
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return valueobjects.UserPreferences{}, err
	}

	return newPrefs, nil
}

// GetUserByEmail retrieves a user by email address
func (s *userApplicationService) GetUserByEmail(email string) (*entities.User, error) {
	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(emailVO)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// ChangeUserEmail changes a user's email address with validation
func (s *userApplicationService) ChangeUserEmail(userID uint, newEmail string) (*entities.User, error) {
	userIDVO, err := valueobjects.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	emailVO, err := valueobjects.NewEmail(newEmail)
	if err != nil {
		return nil, err
	}

	// Validate email uniqueness
	if err := s.authService.ValidateEmailUniqueness(emailVO); err != nil {
		return nil, err
	}

	// Retrieve user
	user, err := s.userRepo.FindByID(userIDVO)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// Change email
	if err := user.ChangeEmail(emailVO); err != nil {
		return nil, err
	}

	// Save updated user
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}