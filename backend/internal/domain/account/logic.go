package account

import (
	"errors"
	"strings"

	domainerr "type-note/backend/internal/domain/errors"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidName  = errors.New("first or last name is required")
)

func Validate(a Account) error {
	if strings.TrimSpace(a.FirstName) == "" && strings.TrimSpace(a.LastName) == "" {
		return ErrInvalidName
	}

	if strings.TrimSpace(a.Provider) == "" {
		return domainerr.ErrProviderRequired
	}

	if strings.TrimSpace(a.ProviderAccountID) == "" {
		return domainerr.ErrProviderAccountRequired
	}

	return nil
}

func UpdateProfile(current Account, input OAuthAccountInput) (Account, error) {
	email, err := ParseEmail(input.Email)
	if err != nil {
		return current, err
	}
	current.Email = email

	if input.Provider != "" {
		current.Provider = input.Provider
	}

	if input.ProviderAccountID != "" {
		current.ProviderAccountID = input.ProviderAccountID
	}

	if input.FirstName != "" {
		current.FirstName = input.FirstName
	}

	if input.LastName != "" {
		current.LastName = input.LastName
	}

	if input.Thumbnail != nil {
		current.Thumbnail = *input.Thumbnail
	}

	return current, Validate(current)
}
