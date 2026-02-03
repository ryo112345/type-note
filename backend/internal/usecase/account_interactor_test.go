package usecase_test

import (
	"context"
	"errors"
	"testing"

	"type-note/backend/internal/domain/account"
	domainerr "type-note/backend/internal/domain/errors"
	"type-note/backend/internal/usecase"
)

// モック: AccountRepository
type mockAccountRepo struct {
	upsertFn     func(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error)
	getByIDFn    func(ctx context.Context, id string) (*account.Account, error)
	getByEmailFn func(ctx context.Context, email string) (*account.Account, error)

	// 呼び出し回数を記録
	upsertCalled     int
	getByIDCalled    int
	getByEmailCalled int
}

func (m *mockAccountRepo) UpsertOAuthAccount(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
	m.upsertCalled++
	if m.upsertFn != nil {
		return m.upsertFn(ctx, input)
	}
	return nil, nil
}

func (m *mockAccountRepo) GetByID(ctx context.Context, id string) (*account.Account, error) {
	m.getByIDCalled++
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *mockAccountRepo) GetByEmail(ctx context.Context, email string) (*account.Account, error) {
	m.getByEmailCalled++
	if m.getByEmailFn != nil {
		return m.getByEmailFn(ctx, email)
	}
	return nil, nil
}

// モック: AccountOutputPort
type mockAccountOutput struct {
	presentFn func(ctx context.Context, a *account.Account) error

	// 渡された値と呼び出し回数を記録
	presented     *account.Account
	presentCalled int
}

func (m *mockAccountOutput) PresentAccount(ctx context.Context, a *account.Account) error {
	m.presentCalled++
	m.presented = a
	if m.presentFn != nil {
		return m.presentFn(ctx, a)
	}
	return nil
}

func TestAccountInteractor_CreateOrGet(t *testing.T) {
	tests := []struct {
		name           string
		input          account.OAuthAccountInput
		repoAcc        *account.Account
		repoErr        error
		wantError      error
		wantRepoCalled bool
		wantPresented  bool
	}{
		{
			name: "成功: アカウント作成",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoAcc:        &account.Account{ID: "acc-1"},
			wantRepoCalled: true,
			wantPresented:  true,
		},
		{
			name: "成功: LastName のみ",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoAcc:        &account.Account{ID: "acc-2"},
			wantRepoCalled: true,
			wantPresented:  true,
		},
		{
			name: "失敗: 無効なメール",
			input: account.OAuthAccountInput{
				Email:             "invalid",
				FirstName:         "Taro",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			wantError:      account.ErrInvalidEmail,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name: "失敗: Provider なし",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				Provider:          "",
				ProviderAccountID: "pid",
			},
			wantError:      domainerr.ErrProviderRequired,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name: "失敗: ProviderAccountID なし",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				Provider:          "google",
				ProviderAccountID: "",
			},
			wantError:      domainerr.ErrProviderAccountRequired,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name: "失敗: Repository エラー",
			input: account.OAuthAccountInput{
				Email:             "user@example.com",
				FirstName:         "Taro",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			repoErr:        errors.New("db connection error"),
			wantError:      errors.New("db connection error"),
			wantRepoCalled: true,
			wantPresented:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAccountRepo{
				upsertFn: func(ctx context.Context, input account.OAuthAccountInput) (*account.Account, error) {
					return tt.repoAcc, tt.repoErr
				},
			}
			out := &mockAccountOutput{}

			interactor := usecase.NewAccountInteractor(repo, out)
			err := interactor.CreateOrGet(context.Background(), tt.input)

			// エラーの検証
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil {
				if err == nil {
					t.Fatalf("want error %v, got nil", tt.wantError)
				}
				if !errors.Is(err, tt.wantError) && err.Error() != tt.wantError.Error() {
					t.Fatalf("want error %v, got %v", tt.wantError, err)
				}
			}

			// Repository が呼ばれたかの検証
			if tt.wantRepoCalled && repo.upsertCalled != 1 {
				t.Fatalf("want repo called once, got %d", repo.upsertCalled)
			}
			if !tt.wantRepoCalled && repo.upsertCalled != 0 {
				t.Fatalf("want repo not called, got %d", repo.upsertCalled)
			}

			// OutputPort が呼ばれたかの検証
			if tt.wantPresented && out.presentCalled != 1 {
				t.Fatalf("want presenter called once, got %d", out.presentCalled)
			}
			if !tt.wantPresented && out.presentCalled != 0 {
				t.Fatalf("want presenter not called, got %d", out.presentCalled)
			}

			// 正しい Account が渡されたかの検証
			if tt.wantPresented && out.presented.ID != tt.repoAcc.ID {
				t.Fatalf("want presented account ID %s, got %s", tt.repoAcc.ID, out.presented.ID)
			}
		})
	}
}

func TestAccountInteractor_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		repoAcc       *account.Account
		repoErr       error
		wantError     error
		wantPresented bool
	}{
		{
			name:          "成功: ID で取得",
			id:            "acc-1",
			repoAcc:       &account.Account{ID: "acc-1", FirstName: "Taro"},
			wantPresented: true,
		},
		{
			name:          "失敗: 見つからない",
			id:            "missing",
			repoErr:       domainerr.ErrNotFound,
			wantError:     domainerr.ErrNotFound,
			wantPresented: false,
		},
		{
			name:          "失敗: DB エラー",
			id:            "acc-1",
			repoErr:       errors.New("db error"),
			wantError:     errors.New("db error"),
			wantPresented: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAccountRepo{
				getByIDFn: func(ctx context.Context, id string) (*account.Account, error) {
					if id != tt.id {
						t.Fatalf("want id %s, got %s", tt.id, id)
					}
					return tt.repoAcc, tt.repoErr
				},
			}
			out := &mockAccountOutput{}

			interactor := usecase.NewAccountInteractor(repo, out)
			err := interactor.GetByID(context.Background(), tt.id)

			// エラーの検証
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil {
				if err == nil {
					t.Fatalf("want error %v, got nil", tt.wantError)
				}
				if !errors.Is(err, tt.wantError) && err.Error() != tt.wantError.Error() {
					t.Fatalf("want error %v, got %v", tt.wantError, err)
				}
			}

			// Repository が呼ばれたかの検証
			if repo.getByIDCalled != 1 {
				t.Fatalf("want repo called once, got %d", repo.getByIDCalled)
			}

			// OutputPort の検証
			if tt.wantPresented && out.presentCalled != 1 {
				t.Fatalf("want presenter called once, got %d", out.presentCalled)
			}
			if !tt.wantPresented && out.presentCalled != 0 {
				t.Fatalf("want presenter not called, got %d", out.presentCalled)
			}

			// 正しい Account が渡されたかの検証
			if tt.wantPresented {
				if out.presented == nil {
					t.Fatal("want presented account, got nil")
				}
				if out.presented.ID != tt.repoAcc.ID {
					t.Fatalf("want presented account ID %s, got %s", tt.repoAcc.ID, out.presented.ID)
				}
			}
		})
	}
}

func TestAccountInteractor_GetByEmail(t *testing.T) {
	email, _ := account.ParseEmail("user@example.com")

	tests := []struct {
		name           string
		email          string
		repoAcc        *account.Account
		repoErr        error
		wantError      error
		wantRepoCalled bool
		wantPresented  bool
	}{
		{
			name:           "成功: メールで取得",
			email:          "user@example.com",
			repoAcc:        &account.Account{ID: "acc-1", Email: email},
			wantRepoCalled: true,
			wantPresented:  true,
		},
		{
			name:           "失敗: 無効なメール形式",
			email:          "invalid-email",
			wantError:      account.ErrInvalidEmail,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name:           "失敗: @ なし",
			email:          "userexample.com",
			wantError:      account.ErrInvalidEmail,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name:           "失敗: 空文字",
			email:          "",
			wantError:      account.ErrInvalidEmail,
			wantRepoCalled: false,
			wantPresented:  false,
		},
		{
			name:           "失敗: 見つからない",
			email:          "missing@example.com",
			repoErr:        domainerr.ErrNotFound,
			wantError:      domainerr.ErrNotFound,
			wantRepoCalled: true,
			wantPresented:  false,
		},
		{
			name:           "失敗: DB エラー",
			email:          "user@example.com",
			repoErr:        errors.New("db error"),
			wantError:      errors.New("db error"),
			wantRepoCalled: true,
			wantPresented:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAccountRepo{
				getByEmailFn: func(ctx context.Context, email string) (*account.Account, error) {
					return tt.repoAcc, tt.repoErr
				},
			}
			out := &mockAccountOutput{}

			interactor := usecase.NewAccountInteractor(repo, out)
			err := interactor.GetByEmail(context.Background(), tt.email)

			// エラーの検証
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil {
				if err == nil {
					t.Fatalf("want error %v, got nil", tt.wantError)
				}
				if !errors.Is(err, tt.wantError) && err.Error() != tt.wantError.Error() {
					t.Fatalf("want error %v, got %v", tt.wantError, err)
				}
			}

			// Repository が呼ばれたかの検証
			if tt.wantRepoCalled && repo.getByEmailCalled != 1 {
				t.Fatalf("want repo called once, got %d", repo.getByEmailCalled)
			}
			if !tt.wantRepoCalled && repo.getByEmailCalled != 0 {
				t.Fatalf("want repo not called, got %d", repo.getByEmailCalled)
			}

			// OutputPort の検証
			if tt.wantPresented && out.presentCalled != 1 {
				t.Fatalf("want presenter called once, got %d", out.presentCalled)
			}
			if !tt.wantPresented && out.presentCalled != 0 {
				t.Fatalf("want presenter not called, got %d", out.presentCalled)
			}

			// 正しい Account が渡されたかの検証
			if tt.wantPresented {
				if out.presented == nil {
					t.Fatal("want presented account, got nil")
				}
				if out.presented.ID != tt.repoAcc.ID {
					t.Fatalf("want presented account ID %s, got %s", tt.repoAcc.ID, out.presented.ID)
				}
			}
		})
	}
}
