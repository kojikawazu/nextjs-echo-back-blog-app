package handlers_comments

import (
	utils "backend/utils/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ブログIDでコメントデータを取得する
func (h *CommentHandler) FetchCommentsByBlogId(c echo.Context) error {
	utils.LogInfo(c, "Fetching comments by blogId...")

	// パスパラメータからblogIdを取得
	blogId := c.Param("blogId")

	// サービス層からブログIDでコメントデータを取得
	comments, err := h.CommentService.FetchCommentsByBlogId(blogId)
	if err != nil {
		switch err.Error() {
		case "invalid blogId":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid blogId",
			})
		case "comments not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Comments not found",
			})
		default:
			utils.LogError(c, "Error fetching comments: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error fetching comments",
			})
		}
	}

	utils.LogInfo(c, "Fetched comments successfully")
	return c.JSON(http.StatusOK, comments)
}

// コメントデータを新規作成する
func (h *CommentHandler) CreateComment(c echo.Context) error {
	utils.LogInfo(c, "Creating comment...")

	// リクエストボディからコメントデータを取得
	req := new(struct {
		BlogId    string `json:"blogId"`
		GuestUser string `json:"guestUser"`
		Comment   string `json:"comment"`
	})
	if err := c.Bind(req); err != nil {
		utils.LogError(c, "Error binding request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Error binding request",
		})
	}

	// サービス層からコメントデータを新規作成
	newComment, err := h.CommentService.CreateComment(req.BlogId, req.GuestUser, req.Comment)
	if err != nil {
		switch err.Error() {
		case "invalid blogId":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid blogId",
			})
		case "invalid guestUser":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid guestUser",
			})
		case "invalid comment":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid comment",
			})
		case "failed to create comment":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create comment",
			})
		default:
			utils.LogError(c, "Error creating comment: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error creating comment",
			})
		}
	}

	utils.LogInfo(c, "Created comment successfully")
	return c.JSON(http.StatusCreated, newComment)
}
