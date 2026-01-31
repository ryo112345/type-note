package account

import (
	"errors"
	"testing"
)

func TestParseEmail_VO(t *testing.T) {
	tests := []struct {
		name      string
		raw       string
		want      Email
		wantError error
	}{
		{
			name: "[Success] trims spaces",
			raw:  "  user@example.com  ",
			want: "user@example.com",
		},
		{
			name:      "[Fail] missing at",
			raw:       "example.com",
			wantError: ErrInvalidEmail,
		},
		{
			name:      "[Fail] empty string",
			raw:       "",
			wantError: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEmail(tt.raw)

			if tt.wantError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantError != nil && !errors.Is(err, tt.wantError) {
				t.Fatalf("want %v, got %v", tt.wantError, err)
			}

			if err == nil && got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	t.Run("[Success] returns underlying string", func(t *testing.T) {
		e := Email("user@example.com")
		if e.String() != "user@example.com" {
			t.Fatalf("unexpected string: %s", e.String())
		}
	})
}

func TestEmail_Domain(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{
			name:  "[Success] returns domain part",
			email: "user@example.com",
			want:  "example.com",
		},
		{
			name:  "[Success] subdomain",
			email: "user@mail.example.com",
			want:  "mail.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.email.Domain()
			if got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestEmail_Local(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{
			name:  "[Success] returns local part",
			email: "user@example.com",
			want:  "user",
		},
		{
			name:  "[Success] with dots",
			email: "first.last@example.com",
			want:  "first.last",
		},
		{
			name:  "[Success] with plus",
			email: "user+tag@example.com",
			want:  "user+tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.email.Local()
			if got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}

func TestEmail_Mask(t *testing.T) {
	tests := []struct {
		name  string
		email Email
		want  string
	}{
		{
			name:  "[Success] masks normal email",
			email: "user@example.com",
			want:  "use******@example.com",
		},
		{
			name:  "[Success] masks short local part (2 chars)",
			email: "ab@example.com",
			want:  "ab******@example.com",
		},
		{
			name:  "[Success] masks single char local",
			email: "a@example.com",
			want:  "a******@example.com",
		},
		{
			name:  "[Success] masks exactly 3 chars",
			email: "abc@example.com",
			want:  "abc******@example.com",
		},
		{
			name:  "[Success] masks long email",
			email: "verylongemail@example.com",
			want:  "ver******@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.email.Mask()
			if got != tt.want {
				t.Fatalf("want %s, got %s", tt.want, got)
			}
		})
	}
}
