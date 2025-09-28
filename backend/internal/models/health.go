package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// HealthStatus represents the overall health status of the service
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// DatabaseStatus represents the database connectivity status
type DatabaseStatus string

const (
	DatabaseStatusConnected    DatabaseStatus = "connected"
	DatabaseStatusDisconnected DatabaseStatus = "disconnected"
	DatabaseStatusError        DatabaseStatus = "error"
)

// HealthResponse represents the response structure for the health endpoint
type HealthResponse struct {
	Status    HealthStatus    `json:"status" validate:"required"`
	Database  DatabaseStatus  `json:"database" validate:"required"`
	Timestamp string          `json:"timestamp" validate:"required"`
	Version   string          `json:"version,omitempty"`
	Uptime    int64           `json:"uptime,omitempty"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error   string `json:"error" validate:"required"`
	Message string `json:"message" validate:"required"`
}

// Validate validates the HealthResponse fields
func (h *HealthResponse) Validate() error {
	// Validate status enum
	if !h.Status.IsValid() {
		return fmt.Errorf("invalid status: %s, must be one of: healthy, degraded, unhealthy", h.Status)
	}

	// Validate database enum
	if !h.Database.IsValid() {
		return fmt.Errorf("invalid database status: %s, must be one of: connected, disconnected, error", h.Database)
	}

	// Validate timestamp format (should be ISO 8601 / RFC3339)
	if h.Timestamp == "" {
		return fmt.Errorf("timestamp cannot be empty")
	}
	if _, err := time.Parse(time.RFC3339, h.Timestamp); err != nil {
		return fmt.Errorf("invalid timestamp format: %s, must be ISO 8601 / RFC3339", h.Timestamp)
	}

	// Validate uptime (must be non-negative if provided)
	if h.Uptime < 0 {
		return fmt.Errorf("uptime must be non-negative, got: %d", h.Uptime)
	}

	// Validate version (must be non-empty if provided)
	if h.Version != "" && strings.TrimSpace(h.Version) == "" {
		return fmt.Errorf("version cannot be empty or whitespace-only")
	}

	return nil
}

// IsValid checks if the HealthStatus is a valid enum value
func (h HealthStatus) IsValid() bool {
	switch h {
	case HealthStatusHealthy, HealthStatusDegraded, HealthStatusUnhealthy:
		return true
	default:
		return false
	}
}

// String returns the string representation of HealthStatus
func (h HealthStatus) String() string {
	return string(h)
}

// MarshalJSON implements custom JSON marshaling for HealthStatus
func (h HealthStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(h))
}

// UnmarshalJSON implements custom JSON unmarshaling for HealthStatus
func (h *HealthStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	status := HealthStatus(str)
	if !status.IsValid() {
		return fmt.Errorf("invalid health status: %s", str)
	}

	*h = status
	return nil
}

// IsValid checks if the DatabaseStatus is a valid enum value
func (d DatabaseStatus) IsValid() bool {
	switch d {
	case DatabaseStatusConnected, DatabaseStatusDisconnected, DatabaseStatusError:
		return true
	default:
		return false
	}
}

// String returns the string representation of DatabaseStatus
func (d DatabaseStatus) String() string {
	return string(d)
}

// MarshalJSON implements custom JSON marshaling for DatabaseStatus
func (d DatabaseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(d))
}

// UnmarshalJSON implements custom JSON unmarshaling for DatabaseStatus
func (d *DatabaseStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	status := DatabaseStatus(str)
	if !status.IsValid() {
		return fmt.Errorf("invalid database status: %s", str)
	}

	*d = status
	return nil
}

// NewHealthResponse creates a new HealthResponse with current timestamp
func NewHealthResponse(status HealthStatus, database DatabaseStatus, version string, uptime int64) *HealthResponse {
	return &HealthResponse{
		Status:    status,
		Database:  database,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
		Uptime:    uptime,
	}
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(errorCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   errorCode,
		Message: message,
	}
}

// DetermineOverallHealth determines the overall health status based on database status
func DetermineOverallHealth(dbStatus DatabaseStatus) HealthStatus {
	switch dbStatus {
	case DatabaseStatusConnected:
		return HealthStatusHealthy
	case DatabaseStatusDisconnected:
		return HealthStatusDegraded
	case DatabaseStatusError:
		return HealthStatusUnhealthy
	default:
		return HealthStatusUnhealthy
	}
}