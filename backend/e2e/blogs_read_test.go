//go:build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: ブログ参照系エンドポイント（一覧・詳細・カテゴリ・タグ・人気）。
func TestE2E_BlogsRead(t *testing.T) {
	client := newClient(t)

	// 一覧 → 200、シード2件以上
	listResp := doJSON(t, client, http.MethodGet, "/api/blogs", nil)
	assert.Equal(t, http.StatusOK, listResp.StatusCode)
	var blogs []map[string]any
	decodeBody(t, listResp, &blogs)
	assert.GreaterOrEqual(t, len(blogs), 2)

	// 詳細 → 200、シードのタイトル
	detailResp := doJSON(t, client, http.MethodGet, "/api/blogs/detail/"+os.Getenv("TEST_BLOG_ID"), nil)
	assert.Equal(t, http.StatusOK, detailResp.StatusCode)
	var blog map[string]any
	decodeBody(t, detailResp, &blog)
	assert.Equal(t, "Test Blog Title", blog["title"])

	// カテゴリ → 200、"Tech" を含む
	catResp := doJSON(t, client, http.MethodGet, "/api/blogs/categories", nil)
	assert.Equal(t, http.StatusOK, catResp.StatusCode)
	var categories []string
	decodeBody(t, catResp, &categories)
	assert.Contains(t, categories, "Tech")

	// タグ → 200、非空
	tagResp := doJSON(t, client, http.MethodGet, "/api/blogs/tags", nil)
	assert.Equal(t, http.StatusOK, tagResp.StatusCode)
	var tags []string
	decodeBody(t, tagResp, &tags)
	assert.NotEmpty(t, tags)

	// 人気 → 200
	popResp := doJSON(t, client, http.MethodGet, "/api/blogs/popular/10", nil)
	defer popResp.Body.Close()
	assert.Equal(t, http.StatusOK, popResp.StatusCode)
}

// 準正常系: 存在しない（が UUID 形式は正しい）ID の詳細取得は 404。
func TestE2E_BlogDetail_NotFound(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodGet, "/api/blogs/detail/00000000-0000-0000-0000-000000000000", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
