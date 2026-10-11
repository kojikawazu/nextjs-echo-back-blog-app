package dto

import (
	"backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 期待値は DTO 化以前（models.BlogLikesData を直接 c.JSON していた時点）の出力。
const goldenLike = `{"id":"l1","blog_id":"b1","visit_id":"v1","created_at":"2026-10-01T12:00:00Z","updated_at":"2026-10-01T12:00:00Z"}`

func newTestLike() models.BlogLikesData {
	return models.BlogLikesData{ID: "l1", BlogId: "b1", VisitId: "v1", CreatedAt: goldenTime, UpdatedAt: goldenTime}
}

// 正常系: 1 件・一覧とも DTO 化以前の JSON と同一（visit_id を含む）。
func TestBlogLikeResponse_MatchesPreDTOJSON(t *testing.T) {
	like := newTestLike()

	assert.JSONEq(t, goldenLike, marshal(t, NewBlogLikeResponse(&like)))
	assert.JSONEq(t, "["+goldenLike+"]", marshal(t, NewBlogLikeResponses([]models.BlogLikesData{like})))
}

// 準正常系: nil は null、空は []（DTO 化以前と同一）。
func TestBlogLikeResponse_EdgeCases(t *testing.T) {
	assert.Equal(t, "null", marshal(t, NewBlogLikeResponses(nil)))
	assert.Equal(t, "[]", marshal(t, NewBlogLikeResponses([]models.BlogLikesData{})))
}

// 準正常系: 一覧の順序と件数を保つ。
func TestNewBlogLikeResponses_PreservesOrder(t *testing.T) {
	a, b := newTestLike(), newTestLike()
	a.ID, b.ID = "first", "second"

	res := NewBlogLikeResponses([]models.BlogLikesData{a, b})

	require.Len(t, res, 2)
	assert.Equal(t, "first", res[0].ID)
	assert.Equal(t, "second", res[1].ID)
}
