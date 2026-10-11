package dto

import (
	"backend/models"
	"time"
)

// BlogResponse はブログ記事のレスポンス。
// キー名・型は DTO 化以前（モデルを直接返していた時点）の JSON と同一に保つ。
type BlogResponse struct {
	ID          string    `json:"id"`
	BlogUserId  string    `json:"blog_user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GithubUrl   string    `json:"github_url"`
	Category    string    `json:"category"`
	Tags        string    `json:"tags"`
	Likes       int8      `json:"likes"`
	CommentCnt  int8      `json:"comment_cnt"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewBlogResponse はブログモデルをレスポンス DTO に変換する。
//
// 引数:
//   - blog: 変換元のブログモデル
//
// 戻り値:
//   - BlogResponse: レスポンス DTO
func NewBlogResponse(blog *models.BlogData) BlogResponse {
	return BlogResponse{
		ID:          blog.ID,
		BlogUserId:  blog.BlogUserId,
		Title:       blog.Title,
		Description: blog.Description,
		GithubUrl:   blog.GithubUrl,
		Category:    blog.Category,
		Tags:        blog.Tags,
		Likes:       blog.Likes,
		CommentCnt:  blog.CommentCnt,
		CreatedAt:   blog.CreatedAt,
		UpdatedAt:   blog.UpdatedAt,
	}
}

// NewBlogResponses はブログモデルの一覧をレスポンス DTO の一覧に変換する。
// 入力が nil の場合は nil を返す（JSON では null。DTO 化以前の挙動を維持する）。
//
// 引数:
//   - blogs: 変換元のブログモデルの一覧
//
// 戻り値:
//   - []BlogResponse: レスポンス DTO の一覧
func NewBlogResponses(blogs []models.BlogData) []BlogResponse {
	if blogs == nil {
		return nil
	}
	res := make([]BlogResponse, 0, len(blogs))
	for i := range blogs {
		res = append(res, NewBlogResponse(&blogs[i]))
	}
	return res
}
