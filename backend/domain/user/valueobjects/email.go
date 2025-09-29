package valueobjects

import (
	"errors"
	"net/mail"
	"strings"
)

// Email represents a validated email address value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object with validation
func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	// Normalize email by trimming spaces and converting to lowercase
	normalizedValue := strings.ToLower(strings.TrimSpace(value))

	// Validate email format using Go's mail package
	if _, err := mail.ParseAddress(normalizedValue); err != nil {
		return Email{}, errors.New("invalid email format")
	}

	// Check maximum length
	if len(normalizedValue) > 255 {
		return Email{}, errors.New("email exceeds maximum length of 255 characters")
	}

	return Email{value: normalizedValue}, nil
}

// Value returns the email address string
func (e Email) Value() string {
	return e.value
}

// String returns the string representation of the email
func (e Email) String() string {
	return e.value
}

// IsEmpty checks if the email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// Domain returns the domain part of the email address
func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// LocalPart returns the local part of the email address (before @)
func (e Email) LocalPart() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}