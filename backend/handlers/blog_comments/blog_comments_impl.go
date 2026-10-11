package handlers_blog_comments

import (
	"backend/apperrors"
	"backend/dto"
	utils "backend/utils/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FetchCommentsByBlogId はブログ ID でコメントデータを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（パスパラメータ blogId を含む）
//
// 戻り値:
//   - error: 不正な blogId は ErrInvalidBlogID、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *CommentHandler) FetchCommentsByBlogId(c echo.Context) error {
	utils.LogInfo(c, "Fetching comments by blogId...")

	// パスパラメータからblogIdを取得
	blogId := c.Param("blogId")

	// サービス層からブログIDでコメントデータを取得
	comments, err := h.CommentService.FetchCommentsByBlogId(blogId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched comments successfully")
	return c.JSON(http.StatusOK, dto.NewCommentResponses(comments))
}

// CreateComment はコメントデータを新規作成するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（blogId/guestUser/comment を含む JSON ボディ）
//
// 戻り値:
//   - error: ボディ不正は apperrors.ErrInvalidRequestBody、不正入力は ErrInvalidXxx、その他の失敗はエラー。成功時は 201 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *CommentHandler) CreateComment(c echo.Context) error {
	utils.LogInfo(c, "Creating comment...")

	// リクエストボディからコメントデータを取得
	req := new(struct {
		BlogId    string `json:"blogId"`
		GuestUser string `json:"guestUser"`
		Comment   string `json:"comment"`
	})
	if err := c.Bind(req); err != nil {
		return apperrors.Wrap(apperrors.ErrInvalidRequestBody, err)
	}

	// サービス層からコメントデータを新規作成
	newComment, err := h.CommentService.CreateComment(req.BlogId, req.GuestUser, req.Comment)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Created comment successfully")
	return c.JSON(http.StatusCreated, dto.NewCommentResponse(newComment))
}
