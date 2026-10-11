//go:build integration

package repositories_blog_users

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

// TestRepository_DoesNotLogCredentials は実 DB に対する取得・更新でログにパスワード・メールアドレスが残らないことを検証する。
// FetchBlogUsersById は password 列を読み出すため、構造体の丸ごと出力があれば DB 上の値がログに載る。
func TestRepository_DoesNotLogCredentials(t *testing.T) {
	testID := os.Getenv("TEST_USER_ID")
	testName := os.Getenv("TEST_USER_NAME")
	testEmail := os.Getenv("TEST_USER_EMAIL")
	testPasswd := os.Getenv("TEST_USER_PASSWD")

	tests := []struct {
		name    string
		run     func(r BlogUsersRepository) error
		wantErr bool
		wantLog string
	}{
		{
			// 正常系: email + password での取得
			name: "fetch_by_email_and_password",
			run: func(r BlogUsersRepository) error {
				_, err := r.FetchBlogUsersByEmailAndPassword(testEmail, testPasswd)
				return err
			},
			wantLog: "Fetched user successfully: id=" + testID,
		},
		{
			// 正常系: password 列を含む ID での取得
			name: "fetch_by_id",
			run: func(r BlogUsersRepository) error {
				_, err := r.FetchBlogUsersById(testID)
				return err
			},
			wantLog: "Fetched user successfully: id=" + testID,
		},
		{
			// 正常系: 更新
			name: "update",
			run: func(r BlogUsersRepository) error {
				_, err := r.UpdateBlogUsers(testID, testName, testEmail, testPasswd)
				return err
			},
			wantLog: "Updated user successfully: id=" + testID,
		},
		{
			// 準正常系: 認証失敗（不一致）時もログに資格情報が残らない
			name: "fetch_by_email_and_wrong_password",
			run: func(r BlogUsersRepository) error {
				_, err := r.FetchBlogUsersByEmailAndPassword(testEmail, "wrong-password-value")
				return err
			},
			wantErr: true,
			wantLog: "User not found or failed to fetch user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := captureStdLog(t)

			err := tt.run(NewBlogUsersRepository())
			assert.Equal(t, tt.wantErr, err != nil, "unexpected error: %v", err)

			logs := buf.String()
			assert.NotContains(t, logs, testPasswd)
			assert.NotContains(t, logs, testEmail)
			assert.NotContains(t, logs, "wrong-password-value")
			assert.Contains(t, logs, tt.wantLog)
		})
	}
}
