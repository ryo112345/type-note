package account

import (
	"errors"
	"testing"

	domainerr "type-note/backend/internal/domain/errors"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		acc       Account
		wantError error
	}{
		{
			name: "[Success] valid account",
			acc: Account{
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
		},
		{
			name: "[Success] only first name",
			acc: Account{
				FirstName:         "Taro",
				LastName:          "",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
		},
		{
			name: "[Success] only last name",
			acc: Account{
				FirstName:         "",
				LastName:          "Yamada",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
		},
		{
			name: "[Fail] missing name",
			acc: Account{
				FirstName:         "",
				LastName:          "",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			wantError: ErrInvalidName,
		},
		{
			name: "[Fail] missing provider",
			acc: Account{
				FirstName:         "Taro",
				LastName:          "Yamada",
				Provider:          "",
				ProviderAccountID: "pid",
			},
			wantError: domainerr.ErrProviderRequired,
		},
		{
			name: "[Fail] missing provider account",
			acc: Account{
				FirstName: "Taro",
				LastName:  "Yamada",
				Provider:  "google",
			},
			wantError: domainerr.ErrProviderAccountRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.acc)
			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	tests := []struct {
		name          string
		current       Account
		input         OAuthAccountInput
		wantEmail     Email
		wantFirstName string
		wantLastName  string
		wantThumbnail string
		wantError     error
	}{
		{
			name: "[Success] merge fields",
			current: Account{
				Email:             "old@example.com",
				FirstName:         "Old",
				LastName:          "Name",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			input: OAuthAccountInput{
				Email:     "new@example.com",
				FirstName: "New",
				LastName:  "Updated",
				Thumbnail: ptr("http://thumb"),
			},
			wantEmail:     "new@example.com",
			wantFirstName: "New",
			wantLastName:  "Updated",
			wantThumbnail: "http://thumb",
		},
		{
			name: "[Success] empty input keeps existing",
			current: Account{
				Email:             "old@example.com",
				FirstName:         "Existing",
				LastName:          "User",
				Provider:          "google",
				ProviderAccountID: "pid",
				Thumbnail:         "old-thumb.jpg",
			},
			input: OAuthAccountInput{
				Email:     "old@example.com",
				FirstName: "",
				LastName:  "",
				Thumbnail: nil,
			},
			wantEmail:     "old@example.com",
			wantFirstName: "Existing",
			wantLastName:  "User",
			wantThumbnail: "old-thumb.jpg",
		},
		{
			name: "[Success] partial update",
			current: Account{
				Email:             "old@example.com",
				FirstName:         "Old",
				LastName:          "Name",
				Provider:          "google",
				ProviderAccountID: "pid",
			},
			input: OAuthAccountInput{
				Email:     "old@example.com",
				FirstName: "New",
				LastName:  "",
			},
			wantEmail:     "old@example.com",
			wantFirstName: "New",
			wantLastName:  "Name",
		},
		{
			name: "[Fail] invalid email",
			current: Account{
				Email: "old@example.com",
			},
			input: OAuthAccountInput{
				Email: "invalid",
			},
			wantError: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := UpdateProfile(tt.current, tt.input)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}

			if err == nil {
				if updated.Email != tt.wantEmail {
					t.Fatalf("email: want %s, got %s", tt.wantEmail, updated.Email)
				}
				if updated.FirstName != tt.wantFirstName {
					t.Fatalf("first name: want %s, got %s", tt.wantFirstName, updated.FirstName)
				}
				if updated.LastName != tt.wantLastName {
					t.Fatalf("last name: want %s, got %s", tt.wantLastName, updated.LastName)
				}
				if tt.wantThumbnail != "" && updated.Thumbnail != tt.wantThumbnail {
					t.Fatalf("thumbnail: want %s, got %s", tt.wantThumbnail, updated.Thumbnail)
				}
			}
		})
	}
}

func ptr[T any](v T) *T { return &v }
