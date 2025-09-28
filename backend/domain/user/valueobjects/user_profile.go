package valueobjects

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// UserProfile represents user profile information value object
type UserProfile struct {
	firstName string
	lastName  string
	timezone  string
}

// NewUserProfile creates a new UserProfile value object with validation
func NewUserProfile(firstName, lastName, timezone string) (UserProfile, error) {
	if err := validateName(firstName, "first name"); err != nil {
		return UserProfile{}, err
	}

	if err := validateName(lastName, "last name"); err != nil {
		return UserProfile{}, err
	}

	if err := validateTimezone(timezone); err != nil {
		return UserProfile{}, err
	}

	return UserProfile{
		firstName: strings.TrimSpace(firstName),
		lastName:  strings.TrimSpace(lastName),
		timezone:  timezone,
	}, nil
}

// validateName validates first and last names
func validateName(name, fieldName string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New(fieldName + " cannot be empty")
	}

	if len(name) > 50 {
		return errors.New(fieldName + " exceeds maximum length of 50 characters")
	}

	// Allow letters, spaces, hyphens, and apostrophes (for names like O'Connor, Mary-Jane)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !nameRegex.MatchString(name) {
		return errors.New(fieldName + " can only contain letters, spaces, hyphens, and apostrophes")
	}

	return nil
}

// validateTimezone validates IANA timezone identifier
func validateTimezone(timezone string) error {
	if timezone == "" {
		return errors.New("timezone cannot be empty")
	}

	// Try to load the timezone to validate it
	_, err := time.LoadLocation(timezone)
	if err != nil {
		return errors.New("invalid timezone: must be a valid IANA timezone identifier")
	}

	return nil
}

// FirstName returns the first name
func (p UserProfile) FirstName() string {
	return p.firstName
}

// LastName returns the last name
func (p UserProfile) LastName() string {
	return p.lastName
}

// Timezone returns the timezone
func (p UserProfile) Timezone() string {
	return p.timezone
}

// FullName returns the full name (first + last)
func (p UserProfile) FullName() string {
	return p.firstName + " " + p.lastName
}

// DisplayName returns a formatted display name
func (p UserProfile) DisplayName() string {
	return p.FullName()
}

// Equals checks if two user profiles are equal
func (p UserProfile) Equals(other UserProfile) bool {
	return p.firstName == other.firstName &&
		p.lastName == other.lastName &&
		p.timezone == other.timezone
}

// WithFirstName returns a new UserProfile with updated first name
func (p UserProfile) WithFirstName(firstName string) (UserProfile, error) {
	return NewUserProfile(firstName, p.lastName, p.timezone)
}

// WithLastName returns a new UserProfile with updated last name
func (p UserProfile) WithLastName(lastName string) (UserProfile, error) {
	return NewUserProfile(p.firstName, lastName, p.timezone)
}

// WithTimezone returns a new UserProfile with updated timezone
func (p UserProfile) WithTimezone(timezone string) (UserProfile, error) {
	return NewUserProfile(p.firstName, p.lastName, timezone)
}