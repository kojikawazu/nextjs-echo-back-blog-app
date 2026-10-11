package services_blog_users

import (
	"backend/models"
	repositories_blog_users "backend/repositories/blog_users"
	"bytes"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	logTestEmail    = "secret-user@example.com"
	logTestPassword = "plain-secret-password"
)

// captureStdLog は標準 log の出力をバッファへ差し替え、テスト終了時に復元する。
func captureStdLog(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(prev) })
	return &buf
}

// TestService_FetchUserByEmailAndPassword_DoesNotLogCredentials はユーザー取得時のログに資格情報が残らないことを検証する。
func TestService_FetchUserByEmailAndPassword_DoesNotLogCredentials(t *testing.T) {
	tests := []struct {
		name    string
		repoRet *models.BlogUsersData
		repoErr error
		wantErr bool
		wantLog string
	}{
		{
			// 正常系
			name:    "found",
			repoRet: &models.BlogUsersData{ID: "1", Email: logTestEmail, Password: mustHash(t, logTestPassword)},
		},
		{
			// 準正常系: ユーザー不在のログにメールアドレスを含めない
			name:    "not_found",
			repoErr: repositories_blog_users.ErrUserNotFound,
			wantErr: true,
			wantLog: "User not found",
		},
		{
			// 異常系: 想定外の取得エラー
			name:    "repository_error",
			repoErr: errors.New("connection refused"),
			wantErr: true,
			wantLog: "Failed to fetch user: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := captureStdLog(t)

			mockRepo := new(repositories_blog_users.MockUserRepository)
			mockRepo.On("FetchBlogUsersByEmail", logTestEmail).Return(tt.repoRet, tt.repoErr)

			_, err := NewUserService(mockRepo).FetchUserByEmailAndPassword(logTestEmail, logTestPassword)

			assert.Equal(t, tt.wantErr, err != nil)
			logs := buf.String()
			assert.NotContains(t, logs, logTestEmail)
			assert.NotContains(t, logs, logTestPassword)
			if tt.wantLog != "" {
				assert.Contains(t, logs, tt.wantLog)
			}
		})
	}
}
