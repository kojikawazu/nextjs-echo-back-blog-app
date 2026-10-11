package dto

import (
	"backend/models"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 期待値は DTO 化以前（models.BlogData を直接 c.JSON していた時点）の出力をそのまま記録したもの。
// API の形（キー名・型・null/[] の扱い）を変えていないことの証明として使う。
const (
	goldenBlog     = `{"id":"b1","blog_user_id":"u1","title":"T","description":"D","github_url":"https://github.com/x/y","category":"Tech","tags":"go,test","likes":3,"comment_cnt":2,"created_at":"2026-10-01T12:00:00Z","updated_at":"2026-10-01T13:00:00Z"}`
	goldenBlogZero = `{"id":"","blog_user_id":"","title":"","description":"","github_url":"","category":"","tags":"","likes":0,"comment_cnt":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`
)

var goldenTime = time.Date(2026, 10, 1, 12, 0, 0, 0, time.UTC)

func newTestBlog() models.BlogData {
	return models.BlogData{
		ID: "b1", BlogUserId: "u1", Title: "T", Description: "D", GithubUrl: "https://github.com/x/y",
		Category: "Tech", Tags: "go,test", Likes: 3, CommentCnt: 2,
		CreatedAt: goldenTime, UpdatedAt: goldenTime.Add(time.Hour),
	}
}

func marshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

// 正常系: 1 件・一覧とも DTO 化以前の JSON と同一。
func TestBlogResponse_MatchesPreDTOJSON(t *testing.T) {
	blog := newTestBlog()

	assert.JSONEq(t, goldenBlog, marshal(t, NewBlogResponse(&blog)))
	assert.JSONEq(t, "["+goldenBlog+"]", marshal(t, NewBlogResponses([]models.BlogData{blog})))
}

// 準正常系: ゼロ値・nil・空の一覧も DTO 化以前と同一（nil は null、空は []）。
func TestBlogResponse_EdgeCases(t *testing.T) {
	assert.JSONEq(t, goldenBlogZero, marshal(t, NewBlogResponse(&models.BlogData{})))
	assert.Equal(t, "null", marshal(t, NewBlogResponses(nil)))
	assert.Equal(t, "[]", marshal(t, NewBlogResponses([]models.BlogData{})))
}

// 準正常系: 一覧の順序と件数を保ち、要素ごとに対応するモデルの値を写す。
func TestNewBlogResponses_PreservesOrder(t *testing.T) {
	a, b := newTestBlog(), newTestBlog()
	a.ID, b.ID = "first", "second"

	res := NewBlogResponses([]models.BlogData{a, b})

	require.Len(t, res, 2)
	assert.Equal(t, "first", res[0].ID)
	assert.Equal(t, "second", res[1].ID)
}
