package port

import (
	"context"

	"type-note/backend/internal/domain/account"
)

// AccountInputPort はアカウントユースケースの入力インターフェース。
type AccountInputPort interface {
	CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error
	GetByID(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) error
}

// AccountOutputPort はアカウントのプレゼンター用インターフェース。
type AccountOutputPort interface {
	PresentAccount(ctx context.Context, account *account.Account) error
}

// AccountRepository はアカウントの永続化を抽象化するインターフェース。
type AccountRepository interface {
	UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error)
	GetByID(ctx context.Context, id string) (*account.Account, error)
	GetByEmail(ctx context.Context, email string) (*account.Account, error)
}
