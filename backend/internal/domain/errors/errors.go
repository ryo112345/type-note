package errors

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrProviderRequired        = errors.New("provider is required")
	ErrProviderAccountRequired = errors.New("provider account id is required")
)
