package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Email represents an email value object
type Email struct {
	value string
}

// NewEmail creates and validates an email
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: email}, nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.value
}

// Value returns the email value
func (e Email) Value() string {
	return e.value
}
