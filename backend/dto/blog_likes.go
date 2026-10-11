package dto

import (
	"backend/models"
	"time"
)

// BlogLikeResponse はいいねのレスポンス。
// キー名・型は DTO 化以前の JSON と同一に保つ（フロントエンドの BlogLike 型が visit_id を参照する）。
type BlogLikeResponse struct {
	ID        string    `json:"id"`
	BlogId    string    `json:"blog_id"`
	VisitId   string    `json:"visit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBlogLikeResponse はいいねモデルをレスポンス DTO に変換する。
//
// 引数:
//   - like: 変換元のいいねモデル
//
// 戻り値:
//   - BlogLikeResponse: レスポンス DTO
func NewBlogLikeResponse(like *models.BlogLikesData) BlogLikeResponse {
	return BlogLikeResponse{
		ID:        like.ID,
		BlogId:    like.BlogId,
		VisitId:   like.VisitId,
		CreatedAt: like.CreatedAt,
		UpdatedAt: like.UpdatedAt,
	}
}

// NewBlogLikeResponses はいいねモデルの一覧をレスポンス DTO の一覧に変換する。
// 入力が nil の場合は nil を返す（JSON では null。DTO 化以前の挙動を維持する）。
//
// 引数:
//   - likes: 変換元のいいねモデルの一覧
//
// 戻り値:
//   - []BlogLikeResponse: レスポンス DTO の一覧
func NewBlogLikeResponses(likes []models.BlogLikesData) []BlogLikeResponse {
	if likes == nil {
		return nil
	}
	res := make([]BlogLikeResponse, 0, len(likes))
	for i := range likes {
		res = append(res, NewBlogLikeResponse(&likes[i]))
	}
	return res
}
