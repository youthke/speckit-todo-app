package valueobjects

import (
	"time"

	userentities "domain/user/entities"
)

// GoogleIdentity represents the link between a User and their Google account
type GoogleIdentity struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	GoogleUserID   string    `json:"google_user_id" gorm:"uniqueIndex;size:255;not null"`
	Email          string    `json:"email" gorm:"size:255;not null"`
	EmailVerified  bool      `json:"email_verified" gorm:"not null;default:false"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationship
	User           userentities.User      `json:"-" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the GoogleIdentity model
func (GoogleIdentity) TableName() string {
	return "google_identities"
}
