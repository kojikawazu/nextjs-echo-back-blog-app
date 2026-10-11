package dto

import (
	"backend/models"
	"time"
)

// CommentResponse はコメントのレスポンス。
// キー名・型は DTO 化以前の JSON と同一に保つ。
type CommentResponse struct {
	ID        string    `json:"id"`
	BlogId    string    `json:"blog_id"`
	GuestUser string    `json:"guest_user"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// NewCommentResponse はコメントモデルをレスポンス DTO に変換する。
//
// 引数:
//   - comment: 変換元のコメントモデル
//
// 戻り値:
//   - CommentResponse: レスポンス DTO
func NewCommentResponse(comment *models.BlogCommentsData) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		BlogId:    comment.BlogId,
		GuestUser: comment.GuestUser,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
	}
}

// NewCommentResponses はコメントモデルの一覧をレスポンス DTO の一覧に変換する。
// 入力が nil の場合は nil を返す（JSON では null。DTO 化以前の挙動を維持する）。
//
// 引数:
//   - comments: 変換元のコメントモデルの一覧
//
// 戻り値:
//   - []CommentResponse: レスポンス DTO の一覧
func NewCommentResponses(comments []models.BlogCommentsData) []CommentResponse {
	if comments == nil {
		return nil
	}
	res := make([]CommentResponse, 0, len(comments))
	for i := range comments {
		res = append(res, NewCommentResponse(&comments[i]))
	}
	return res
}
