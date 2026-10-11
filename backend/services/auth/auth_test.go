package services_auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestService_Login はログイン入力のバリデーションを検証する。
func TestService_Login(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  error
	}{
		{
			// 正常系: メール形式が正しく、パスワードがあればエラーなし
			name:     "valid_input",
			email:    "test@example.com",
			password: "password123",
			wantErr:  nil,
		},
		{
			// 準正常系: メールアドレスが空
			name:     "empty_email",
			email:    "",
			password: "password123",
			wantErr:  ErrEmailAndPasswordRequired,
		},
		{
			// 準正常系: パスワードが空
			name:     "empty_password",
			email:    "test@example.com",
			password: "",
			wantErr:  ErrEmailAndPasswordRequired,
		},
		{
			// 準正常系: メール形式の不正
			name:     "invalid_email_format",
			email:    "invalid-email",
			password: "password123",
			wantErr:  ErrInvalidEmailFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewAuthService().Login(tt.email, tt.password)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
