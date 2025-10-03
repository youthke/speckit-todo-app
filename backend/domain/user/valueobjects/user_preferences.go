package valueobjects

import (
	"errors"

	"domain/task/valueobjects"
)

// UserPreferences represents user preference settings value object
type UserPreferences struct {
	defaultTaskPriority valueobjects.TaskPriority
	emailNotifications  bool
	themePreference     string
}

// Valid theme preferences
const (
	ThemeLight = "light"
	ThemeDark  = "dark"
	ThemeAuto  = "auto"
)

// NewUserPreferences creates a new UserPreferences value object with validation
func NewUserPreferences(
	defaultTaskPriority valueobjects.TaskPriority,
	emailNotifications bool,
	themePreference string,
) (UserPreferences, error) {
	if err := validateThemePreference(themePreference); err != nil {
		return UserPreferences{}, err
	}

	return UserPreferences{
		defaultTaskPriority: defaultTaskPriority,
		emailNotifications:  emailNotifications,
		themePreference:     themePreference,
	}, nil
}

// NewDefaultUserPreferences creates UserPreferences with sensible defaults
func NewDefaultUserPreferences() UserPreferences {
	defaultPriority := valueobjects.NewMediumPriority()
	prefs, _ := NewUserPreferences(defaultPriority, true, ThemeAuto)
	return prefs
}

// validateThemePreference validates that the theme preference is valid
func validateThemePreference(theme string) error {
	switch theme {
	case ThemeLight, ThemeDark, ThemeAuto:
		return nil
	default:
		return errors.New("theme preference must be 'light', 'dark', or 'auto'")
	}
}

// DefaultTaskPriority returns the default task priority
func (p UserPreferences) DefaultTaskPriority() valueobjects.TaskPriority {
	return p.defaultTaskPriority
}

// EmailNotifications returns whether email notifications are enabled
func (p UserPreferences) EmailNotifications() bool {
	return p.emailNotifications
}

// ThemePreference returns the theme preference
func (p UserPreferences) ThemePreference() string {
	return p.themePreference
}

// Equals checks if two user preferences are equal
func (p UserPreferences) Equals(other UserPreferences) bool {
	return p.defaultTaskPriority.Equals(other.defaultTaskPriority) &&
		p.emailNotifications == other.emailNotifications &&
		p.themePreference == other.themePreference
}

// WithDefaultTaskPriority returns new UserPreferences with updated default task priority
func (p UserPreferences) WithDefaultTaskPriority(priority valueobjects.TaskPriority) UserPreferences {
	prefs, _ := NewUserPreferences(priority, p.emailNotifications, p.themePreference)
	return prefs
}

// WithEmailNotifications returns new UserPreferences with updated email notification setting
func (p UserPreferences) WithEmailNotifications(enabled bool) UserPreferences {
	prefs, _ := NewUserPreferences(p.defaultTaskPriority, enabled, p.themePreference)
	return prefs
}

// WithThemePreference returns new UserPreferences with updated theme preference
func (p UserPreferences) WithThemePreference(theme string) (UserPreferences, error) {
	return NewUserPreferences(p.defaultTaskPriority, p.emailNotifications, theme)
}

// IsLightTheme returns true if the theme preference is light
func (p UserPreferences) IsLightTheme() bool {
	return p.themePreference == ThemeLight
}

// IsDarkTheme returns true if the theme preference is dark
func (p UserPreferences) IsDarkTheme() bool {
	return p.themePreference == ThemeDark
}

// IsAutoTheme returns true if the theme preference is auto
func (p UserPreferences) IsAutoTheme() bool {
	return p.themePreference == ThemeAuto
}