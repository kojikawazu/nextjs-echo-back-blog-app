package handlers_blogs

import (
	utils "backend/utils/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 全ブログデータを取得する
func (h *BlogHandler) FetchBlogs(c echo.Context) error {
	utils.LogInfo(c, "Fetching blogs...")

	// サービス層から全ブログデータを取得
	blogs, err := h.BlogService.FetchBlogs()
	if err != nil {
		utils.LogError(c, "Error fetching blogs: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error fetching blogs",
		})
	}

	utils.LogInfo(c, "Fetched blogs successfully")
	return c.JSON(http.StatusOK, blogs)
}

// ユーザーIDでブログデータを取得する
func (h *BlogHandler) FetchBlogsByUserId(c echo.Context) error {
	utils.LogInfo(c, "Fetching blogs by userId...")

	// パスパラメータからuserIdを取得
	userId := c.Param("userId")

	// サービス層からユーザーIDでブログデータを取得
	blogs, err := h.BlogService.FetchBlogsByUserId(userId)
	if err != nil {
		switch err.Error() {
		case "invalid userId":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid userId",
			})
		case "blog not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Blog not found",
			})
		default:
			utils.LogError(c, "Error fetching blogs: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error fetching blogs",
			})
		}
	}

	utils.LogInfo(c, "Fetched blogs successfully")
	return c.JSON(http.StatusOK, blogs)
}
