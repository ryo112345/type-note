package account

import "time"

type Account struct {
	ID                string
	Email             Email
	FirstName         string
	LastName          string
	IsActive          bool
	Provider          string
	ProviderAccountID string
	Thumbnail         string
	LastLoginAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type OAuthAccountInput struct {
	Email             string
	FirstName         string
	LastName          string
	Provider          string
	ProviderAccountID string
	Thumbnail         *string
}
