package usecase

import (
	"context"

	"type-note/backend/internal/domain/account"
	"type-note/backend/internal/port"
)

// AccountInteractor はアカウントのユースケースを実装する。
type AccountInteractor struct {
	repo   port.AccountRepository
	output port.AccountOutputPort
}

// AccountInteractor が AccountInputPort を満たすことを確認
var _ port.AccountInputPort = (*AccountInteractor)(nil)

// NewAccountInteractor は AccountInteractor を生成する。
func NewAccountInteractor(repo port.AccountRepository, output port.AccountOutputPort) *AccountInteractor {
	return &AccountInteractor{
		repo:   repo,
		output: output,
	}
}

// CreateOrGet は OAuth ログインを処理し、アカウントを作成または取得する。
func (u *AccountInteractor) CreateOrGet(ctx context.Context, input account.OAuthAccountInput) error {
	email, err := account.ParseEmail(input.Email)
	if err != nil {
		return err
	}

	acc := account.Account{
		Email:             email,
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Provider:          input.Provider,
		ProviderAccountID: input.ProviderAccountID,
		Thumbnail:         valueOrEmpty(input.Thumbnail),
	}
	if err := account.Validate(acc); err != nil {
		return err
	}

	a, err := u.repo.UpsertOAuthAccount(ctx, input)
	if err != nil {
		return err
	}

	return u.output.PresentAccount(ctx, a)
}

// GetByID は ID でアカウントを取得する。
func (u *AccountInteractor) GetByID(ctx context.Context, id string) error {
	a, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return u.output.PresentAccount(ctx, a)
}

// GetByEmail はメールアドレスでアカウントを取得する。
func (u *AccountInteractor) GetByEmail(ctx context.Context, email string) error {
	if _, err := account.ParseEmail(email); err != nil {
		return err
	}

	a, err := u.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	return u.output.PresentAccount(ctx, a)
}

func valueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
