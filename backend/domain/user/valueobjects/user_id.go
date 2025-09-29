package valueobjects

import "fmt"

// UserID represents a unique identifier for a user
type UserID struct {
	value uint
}

// NewUserID creates a new UserID
func NewUserID(id uint) UserID {
	return UserID{value: id}
}

// Value returns the underlying ID value
func (u UserID) Value() uint {
	return u.value
}

// Equals checks if two UserIDs are equal
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// String returns the string representation of the UserID
func (u UserID) String() string {
	return fmt.Sprintf("UserID(%d)", u.value)
}

// IsZero checks if the UserID is zero (uninitialized)
func (u UserID) IsZero() bool {
	return u.value == 0
}