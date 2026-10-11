package models

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testPassword = "plain-secret-password"
	testEmail    = "secret-user@example.com"
)

func newTestUser() BlogUsersData {
	ts := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	return BlogUsersData{
		ID:        "user-id-123",
		Name:      "Test User",
		Email:     testEmail,
		Password:  testPassword,
		CreatedAt: ts,
		UpdatedAt: ts,
	}
}

// 正常系: String は識別に必要な ID / Name を含み、機密値はマスク表記になる。
func TestBlogUsersData_String(t *testing.T) {
	user := newTestUser()

	assert.Equal(t,
		"BlogUsersData{ID:user-id-123 Name:Test User Email:[REDACTED] Password:[REDACTED] CreatedAt:2026-10-01T12:00:00Z UpdatedAt:2026-10-01T12:00:00Z}",
		user.String())
	assert.Equal(t, "models."+user.String(), user.GoString())
}

// 準正常系: fmt のどの書式・値/ポインタ/スライス経由でもパスワード・メールアドレスが出力されない。
func TestBlogUsersData_FormatDoesNotLeakSecrets(t *testing.T) {
	user := newTestUser()

	tests := []struct {
		name   string
		format string
		arg    any
	}{
		{name: "value_%v", format: "%v", arg: user},
		{name: "value_%+v", format: "%+v", arg: user},
		{name: "value_%s", format: "%s", arg: user},
		{name: "value_%#v", format: "%#v", arg: user},
		{name: "pointer_%v", format: "%v", arg: &user},
		{name: "pointer_%+v", format: "%+v", arg: &user},
		{name: "pointer_%#v", format: "%#v", arg: &user},
		{name: "slice_%v", format: "%v", arg: []BlogUsersData{user}},
		{name: "slice_%+v", format: "%+v", arg: []BlogUsersData{user}},
		{name: "slice_%#v", format: "%#v", arg: []BlogUsersData{user}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := fmt.Sprintf(tt.format, tt.arg)

			assert.NotContains(t, out, testPassword)
			assert.NotContains(t, out, testEmail)
			assert.Contains(t, out, "user-id-123")
			assert.Contains(t, out, "[REDACTED]")
		})
	}
}

// 準正常系: log.Printf に構造体をそのまま渡してもログにパスワード・メールアドレスが残らない。
func TestBlogUsersData_LogDoesNotLeakSecrets(t *testing.T) {
	user := newTestUser()

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	logger.Printf("Fetched user successfully: %v", user)
	logger.Printf("Fetched user successfully: %+v", &user)

	assert.NotContains(t, buf.String(), testPassword)
	assert.NotContains(t, buf.String(), testEmail)
}

// 準正常系: パスワードが空でもマスク表記は変わらない（空かどうかを推測させない）。
func TestBlogUsersData_String_EmptyPassword(t *testing.T) {
	user := newTestUser()
	user.Password = ""

	assert.Contains(t, user.String(), "Password:[REDACTED]")
}
