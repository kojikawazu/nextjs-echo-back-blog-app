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

// ブログIDでブログデータを取得する
func (h *BlogHandler) FetchBlogById(c echo.Context) error {
	utils.LogInfo(c, "Fetching blog by id...")

	// パスパラメータからidを取得
	id := c.Param("id")

	// サービス層からIDでブログデータを取得
	blog, err := h.BlogService.FetchBlogById(id)
	if err != nil {
		switch err.Error() {
		case "invalid id":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid id",
			})
		case "blog not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Blog not found",
			})
		default:
			utils.LogError(c, "Error fetching blog: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error fetching blog",
			})
		}
	}

	utils.LogInfo(c, "Fetched blog successfully")
	return c.JSON(http.StatusOK, blog)
}

// ブログデータを作成する
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	utils.LogInfo(c, "Creating blog...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c)
	if err != nil {
		utils.LogError(c, "Error getting cookie: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting cookie",
		})
	}

	// JWTトークンを解析してユーザーIDを取得
	userId, err := h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting userId from token: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting userId from token",
		})
	}

	// JSONボディのバインド
	type CreateBlogRequest struct {
		Title       string `json:"title"`
		GitHubURL   string `json:"githubUrl"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}

	var req CreateBlogRequest
	if err := c.Bind(&req); err != nil {
		utils.LogError(c, "Error binding request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// サービス層からブログデータを作成
	blog, err := h.BlogService.CreateBlog(userId, req.Title, req.GitHubURL, req.Category, req.Description, req.Tags)
	if err != nil {
		switch err.Error() {
		case "invalid userId":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid userId",
			})
		case "invalid title":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid title",
			})
		case "invalid githubUrl":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid githubUrl",
			})
		case "invalid category":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid category",
			})
		case "invalid description":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid description",
			})
		case "invalid tags":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid tags",
			})
		case "failed to create blog":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create blog",
			})
		default:
			utils.LogError(c, "server error: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Server error",
			})
		}
	}

	utils.LogInfo(c, "Created blog successfully")
	return c.JSON(http.StatusCreated, blog)
}

// ブログデータを更新する
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	utils.LogInfo(c, "Updating blog...")

	// パスパラメータからidを取得
	id := c.Param("id")

	// JSONボディのバインド
	type UpdateBlogRequest struct {
		Title       string `json:"title"`
		GitHubURL   string `json:"githubUrl"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Tags        string `json:"tags"`
	}

	var req UpdateBlogRequest
	if err := c.Bind(&req); err != nil {
		utils.LogError(c, "Error binding request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// サービス層からブログデータを更新
	blog, err := h.BlogService.UpdateBlog(id, req.Title, req.GitHubURL, req.Category, req.Description, req.Tags)
	if err != nil {
		switch err.Error() {
		case "invalid id":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid id",
			})
		case "invalid title":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid title",
			})
		case "invalid githubUrl":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid githubUrl",
			})
		case "invalid category":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid category",
			})
		case "invalid description":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid description",
			})
		case "invalid tags":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid tags",
			})
		case "failed to update blog":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update blog",
			})
		default:
			utils.LogError(c, "server error: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Server error",
			})
		}
	}

	utils.LogInfo(c, "Updated blog successfully")
	return c.JSON(http.StatusOK, blog)
}

// ブログデータを削除する
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	utils.LogInfo(c, "Deleting blog...")

	// パスパラメータからidを取得
	id := c.Param("id")

	// サービス層からブログデータを削除
	err := h.BlogService.DeleteBlog(id)
	if err != nil {
		switch err.Error() {
		case "invalid id":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid id",
			})
		case "failed to delete blog":
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to delete blog",
			})
		default:
			utils.LogError(c, "server error: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Server error",
			})
		}
	}

	utils.LogInfo(c, "Deleted blog successfully")
	return c.NoContent(http.StatusNoContent)
}
