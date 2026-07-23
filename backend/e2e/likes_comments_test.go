//go:build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系（いいね）: 訪問者ID発行 → いいね作成 → 状態確認 → 削除 の一連フロー。
func TestE2E_BlogLikeFlow(t *testing.T) {
	client := newClient(t)
	blogID := os.Getenv("TEST_BLOG_ID")

	// 訪問者ID発行 → 200（visit-id-token Cookie を保存）
	genResp := doJSON(t, client, http.MethodGet, "/api/blog-likes/generate-visit-id", nil)
	defer genResp.Body.Close()
	assert.Equal(t, http.StatusOK, genResp.StatusCode)

	// いいね作成 → 200
	createResp := doJSON(t, client, http.MethodPost, "/api/blog-likes/create/"+blogID, nil)
	defer createResp.Body.Close()
	assert.Equal(t, http.StatusOK, createResp.StatusCode)

	// いいね状態確認 → 200 かつ isLiked=true
	isLikedResp := doJSON(t, client, http.MethodGet, "/api/blog-likes/is-liked/"+blogID, nil)
	assert.Equal(t, http.StatusOK, isLikedResp.StatusCode)
	var liked map[string]bool
	decodeBody(t, isLikedResp, &liked)
	assert.True(t, liked["isLiked"])

	// いいね削除 → 200
	deleteResp := doJSON(t, client, http.MethodDelete, "/api/blog-likes/delete/"+blogID, nil)
	defer deleteResp.Body.Close()
	assert.Equal(t, http.StatusOK, deleteResp.StatusCode)
}

// 正常系（コメント）: コメント一覧取得 → コメント作成。
func TestE2E_CommentFlow(t *testing.T) {
	client := newClient(t)
	blogID := os.Getenv("TEST_BLOG_ID")

	// 一覧取得 → 200（シード済みコメントを含む）
	listResp := doJSON(t, client, http.MethodGet, "/api/comments/blog/"+blogID, nil)
	assert.Equal(t, http.StatusOK, listResp.StatusCode)
	var comments []map[string]any
	decodeBody(t, listResp, &comments)
	assert.GreaterOrEqual(t, len(comments), 1)

	// 作成 → 201
	createResp := doJSON(t, client, http.MethodPost, "/api/comments/create", map[string]string{
		"blogId":    blogID,
		"guestUser": "E2E Guest",
		"comment":   "e2e comment",
	})
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var created struct {
		ID      string `json:"id"`
		Comment string `json:"comment"`
	}
	decodeBody(t, createResp, &created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "e2e comment", created.Comment)
}
