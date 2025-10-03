package services

import (
	"errors"
	"time"

	"domain/user/entities"
	"domain/user/repositories"
	"domain/user/valueobjects"
)

// ProfileUpdateData represents data for profile updates
type ProfileUpdateData struct {
	FirstName *string
	LastName  *string
	Timezone  *string
}

// UserProfileService provides domain profile management logic for users
type UserProfileService interface {
	// UpdateProfile updates user profile information with validation
	UpdateProfile(userID valueobjects.UserID, profile valueobjects.UserProfile) error

	// ValidateProfileData validates profile data before updates
	ValidateProfileData(profile valueobjects.UserProfile) error

	// UpdatePartialProfile updates only specified profile fields
	UpdatePartialProfile(userID valueobjects.UserID, updates ProfileUpdateData) error

	// ValidateTimezoneChange validates timezone changes for consistency
	ValidateTimezoneChange(userID valueobjects.UserID, newTimezone string) error
}

// userProfileService implements UserProfileService
type userProfileService struct {
	userRepo repositories.UserRepository
}

// NewUserProfileService creates a new user profile service
func NewUserProfileService(userRepo repositories.UserRepository) UserProfileService {
	return &userProfileService{
		userRepo: userRepo,
	}
}

// UpdateProfile updates the complete user profile
func (s *userProfileService) UpdateProfile(userID valueobjects.UserID, profile valueobjects.UserProfile) error {
	// Validate the profile data
	if err := s.ValidateProfileData(profile); err != nil {
		return err
	}

	// Retrieve the user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Update the profile
	if err := user.UpdateProfile(profile); err != nil {
		return err
	}

	// Save the updated user
	return s.userRepo.Update(user)
}

// ValidateProfileData validates profile data for business rules
func (s *userProfileService) ValidateProfileData(profile valueobjects.UserProfile) error {
	// Check for required fields (already validated in value object constructor, but ensure business rules)
	if profile.FirstName() == "" {
		return errors.New("first name is required")
	}

	if profile.LastName() == "" {
		return errors.New("last name is required")
	}

	if profile.Timezone() == "" {
		return errors.New("timezone is required")
	}

	// Additional business rules can be added here
	// For example, checking against a blacklist of names, validating timezone against business hours, etc.

	return nil
}

// UpdatePartialProfile updates only specified profile fields
func (s *userProfileService) UpdatePartialProfile(userID valueobjects.UserID, updates ProfileUpdateData) error {
	// Retrieve the current user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	currentProfile := user.Profile()

	// Build the new profile with updates
	firstName := currentProfile.FirstName()
	lastName := currentProfile.LastName()
	timezone := currentProfile.Timezone()

	if updates.FirstName != nil {
		firstName = *updates.FirstName
	}

	if updates.LastName != nil {
		lastName = *updates.LastName
	}

	if updates.Timezone != nil {
		timezone = *updates.Timezone
	}

	// Create new profile with validation
	newProfile, err := valueobjects.NewUserProfile(firstName, lastName, timezone)
	if err != nil {
		return err
	}

	// Update the user
	if err := user.UpdateProfile(newProfile); err != nil {
		return err
	}

	// Save the updated user
	return s.userRepo.Update(user)
}

// ValidateTimezoneChange validates timezone changes for business consistency
func (s *userProfileService) ValidateTimezoneChange(userID valueobjects.UserID, newTimezone string) error {
	// Validate that the timezone is valid
	_, err := time.LoadLocation(newTimezone)
	if err != nil {
		return errors.New("invalid timezone: must be a valid IANA timezone identifier")
	}

	// Additional business rules for timezone changes can be added here
	// For example:
	// - Checking if timezone change affects scheduled tasks
	// - Validating against company timezone policies
	// - Ensuring timezone is within acceptable business regions

	return nil
}