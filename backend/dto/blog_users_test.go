package dto

import (
	"backend/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 正常系: モデルの公開フィールドがそのまま DTO に写る。
func TestNewUserResponse(t *testing.T) {
	ts := time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)
	user := &models.BlogUsersData{
		ID:        "user-id-123",
		Name:      "Test User",
		Email:     "user@example.com",
		Password:  "plain-secret-password",
		CreatedAt: ts,
		UpdatedAt: ts.Add(time.Hour),
	}

	got := NewUserResponse(user)

	assert.Equal(t, UserResponse{
		ID:        "user-id-123",
		Name:      "Test User",
		Email:     "user@example.com",
		CreatedAt: ts,
		UpdatedAt: ts.Add(time.Hour),
	}, got)
}

// 準正常系: JSON のキーは公開フィールドのみで、モデルにパスワードがあっても出力されない。
func TestUserResponse_JSONHasOnlyPublicKeys(t *testing.T) {
	user := &models.BlogUsersData{ID: "user-id-123", Name: "Test User", Email: "user@example.com", Password: "plain-secret-password"}

	b, err := json.Marshal(NewUserResponse(user))
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(b, &m))
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, []string{"id", "name", "email", "created_at", "updated_at"}, keys)
	assert.NotContains(t, string(b), "plain-secret-password")
}

// 準正常系: パスワードが空のモデルでも password キーは出力されない（空文字で露出しない）。
func TestUserResponse_JSONOmitsPasswordKeyWhenEmpty(t *testing.T) {
	b, err := json.Marshal(NewUserResponse(&models.BlogUsersData{ID: "user-id-123"}))
	require.NoError(t, err)

	assert.NotContains(t, string(b), "password")
}
