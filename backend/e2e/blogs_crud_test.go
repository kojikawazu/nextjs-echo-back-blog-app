//go:build e2e

package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: 認証あり CRUD（作成 → 詳細取得 → 更新 → 削除）の一連フロー。
func TestE2E_BlogCRUD(t *testing.T) {
	client := newClient(t)
	login(t, client)

	// 作成 → 201
	createResp := doJSON(t, client, http.MethodPost, "/api/blogs/create", map[string]string{
		"title":       "E2E Blog",
		"githubUrl":   "https://github.com/test/e2e",
		"category":    "E2E",
		"description": "created by e2e test",
		"tags":        "e2e,go",
	})
	assert.Equal(t, http.StatusCreated, createResp.StatusCode)
	var created struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	decodeBody(t, createResp, &created)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "E2E Blog", created.Title)

	// 詳細取得 → 200
	detailResp := doJSON(t, client, http.MethodGet, "/api/blogs/detail/"+created.ID, nil)
	assert.Equal(t, http.StatusOK, detailResp.StatusCode)
	var fetched struct {
		Title string `json:"title"`
	}
	decodeBody(t, detailResp, &fetched)
	assert.Equal(t, "E2E Blog", fetched.Title)

	// 更新 → 200
	updateResp := doJSON(t, client, http.MethodPut, "/api/blogs/update/"+created.ID, map[string]string{
		"title":       "E2E Blog Updated",
		"githubUrl":   "https://github.com/test/e2e",
		"category":    "E2E",
		"description": "updated by e2e test",
		"tags":        "e2e,go",
	})
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)
	var updated struct {
		Title string `json:"title"`
	}
	decodeBody(t, updateResp, &updated)
	assert.Equal(t, "E2E Blog Updated", updated.Title)

	// 削除 → 204
	deleteResp := doJSON(t, client, http.MethodDelete, "/api/blogs/delete/"+created.ID, nil)
	defer deleteResp.Body.Close()
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	// 削除後の詳細取得 → 404
	afterResp := doJSON(t, client, http.MethodGet, "/api/blogs/detail/"+created.ID, nil)
	defer afterResp.Body.Close()
	assert.Equal(t, http.StatusNotFound, afterResp.StatusCode)
}

// 準正常系: 未認証でのブログ作成は 401。
func TestE2E_CreateBlog_Unauthorized(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodPost, "/api/blogs/create", map[string]string{
		"title":       "no auth",
		"githubUrl":   "https://github.com/test/x",
		"category":    "x",
		"description": "x",
		"tags":        "x",
	})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
