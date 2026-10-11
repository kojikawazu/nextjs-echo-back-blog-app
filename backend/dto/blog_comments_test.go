package dto

import (
	"backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 期待値は DTO 化以前（models.BlogCommentsData を直接 c.JSON していた時点）の出力。
const (
	goldenComment     = `{"id":"c1","blog_id":"b1","guest_user":"guest","comment":"hello","created_at":"2026-10-01T12:00:00Z"}`
	goldenCommentZero = `{"id":"","blog_id":"","guest_user":"","comment":"","created_at":"0001-01-01T00:00:00Z"}`
)

func newTestComment() models.BlogCommentsData {
	return models.BlogCommentsData{ID: "c1", BlogId: "b1", GuestUser: "guest", Comment: "hello", CreatedAt: goldenTime}
}

// 正常系: 1 件・一覧とも DTO 化以前の JSON と同一。
func TestCommentResponse_MatchesPreDTOJSON(t *testing.T) {
	comment := newTestComment()

	assert.JSONEq(t, goldenComment, marshal(t, NewCommentResponse(&comment)))
	assert.JSONEq(t, "["+goldenComment+"]", marshal(t, NewCommentResponses([]models.BlogCommentsData{comment})))
}

// 準正常系: ゼロ値・nil・空の一覧も DTO 化以前と同一。
func TestCommentResponse_EdgeCases(t *testing.T) {
	assert.JSONEq(t, goldenCommentZero, marshal(t, NewCommentResponse(&models.BlogCommentsData{})))
	assert.Equal(t, "null", marshal(t, NewCommentResponses(nil)))
	assert.Equal(t, "[]", marshal(t, NewCommentResponses([]models.BlogCommentsData{})))
}

// 準正常系: 一覧の順序と件数を保つ。
func TestNewCommentResponses_PreservesOrder(t *testing.T) {
	a, b := newTestComment(), newTestComment()
	a.ID, b.ID = "first", "second"

	res := NewCommentResponses([]models.BlogCommentsData{a, b})

	require.Len(t, res, 2)
	assert.Equal(t, "first", res[0].ID)
	assert.Equal(t, "second", res[1].ID)
}
