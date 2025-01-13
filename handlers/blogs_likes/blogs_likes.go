package handlers_blogs_likes

import (
	utils "backend/utils/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FetchBlogLikesByVisitId - 訪問IDに紐づくブログいいねデータを取得するハンドラ
func (h *BlogLikeHandler) FetchBlogLikesByVisitId(c echo.Context) error {
	utils.LogInfo(c, "Fetching blog likes by visit id...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err != nil {
		utils.LogError(c, "Error getting visit id token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id token",
		})
	}
	// JWTトークンを解析して訪問IDを取得
	visitId, err := h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting visit id: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id",
		})
	}

	// VisitIDに紐づくいいねデータを取得
	blogLikesData, err := h.BlogLikeService.FetchBlogLikesByVisitId(visitId)
	if err != nil {
		utils.LogError(c, "Error fetching blog likes by visit id: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error fetching blog likes by visit id",
		})
	}

	utils.LogInfo(c, "Blog likes fetched successfully")
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, blogLikesData)
}

// GenerateVisitorId -　訪問者IDを生成するハンドラ
func (h *BlogLikeHandler) GenerateVisitorId(c echo.Context) error {
	utils.LogInfo(c, "Generating visitor id...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err == nil {
		// JWTトークンを解析して訪問IDを取得
		_, err = h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
		// エラーがない場合は、すでに訪問者IDが存在する(スキップ)
		if err == nil {
			utils.LogInfo(c, "Visitor id already exists")
			return c.JSON(http.StatusOK, map[string]string{
				"message": "Visitor id already exists",
			})
		}
	}

	// Visit用トークンの有効期限を1時間に設定
	expirationTime := h.CookieUtils.GetAuthCookieExpirationTime()
	// Visit用トークンの作成
	tokenString, err := h.CookieUtils.CreateVisitIdToken()
	if err != nil {
		utils.LogError(c, "Error creating visitor token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create visitor token",
		})
	}
	//  Visit用トークンをクッキーに保存
	h.CookieUtils.AddVisitIdCoookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "Visitor id generated successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Visitor id generated successfully",
	})
}

// ブログいいねの取得ハンドラ
func (h *BlogLikeHandler) IsBlogLiked(c echo.Context) error {
	utils.LogInfo(c, "Checking if blog is liked...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err != nil {
		utils.LogError(c, "Error getting visit id token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id token",
		})
	}
	// JWTトークンを解析して訪問IDを取得
	visitId, err := h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting visit id: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id",
		})
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータが存在するか確認
	isLiked, err := h.BlogLikeService.IsBlogLiked(blogId, visitId)
	if err != nil {
		utils.LogInfo(c, "Blog like not found")
		return c.JSON(http.StatusOK, map[string]bool{
			"isLiked": isLiked,
		})
	}

	utils.LogInfo(c, "Blog like found")
	return c.JSON(http.StatusOK, map[string]bool{
		"isLiked": isLiked,
	})
}

// ブログいいねの追加ハンドラ
func (h *BlogLikeHandler) CreateBlogLike(c echo.Context) error {
	utils.LogInfo(c, "Creating blog like...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err != nil {
		utils.LogError(c, "Error getting visit id token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id token",
		})
	}
	// JWTトークンを解析して訪問IDを取得
	visitId, err := h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting visit id: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id",
		})
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータを作成
	createdBlogLikeData, err := h.BlogLikeService.CreateBlogLike(blogId, visitId)
	if err != nil {
		switch err.Error() {
		case "BlogId or VisitId is empty":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "BlogId or VisitId is empty",
			})
		case "blog is already liked":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Blog is already liked",
			})
		default:
			utils.LogError(c, "Error creating blog like: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Error creating blog like",
			})
		}
	}

	utils.LogInfo(c, "Blog like created successfully")
	return c.JSON(http.StatusOK, createdBlogLikeData)
}

// ブログいいねの削除ハンドラ
func (h *BlogLikeHandler) DeleteBlogLike(c echo.Context) error {
	utils.LogInfo(c, "Deleting blog like...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err != nil {
		utils.LogError(c, "Error getting visit id token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id token",
		})
	}
	// JWTトークンを解析して訪問IDを取得
	visitId, err := h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting visit id: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get visit id",
		})
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータを削除
	err = h.BlogLikeService.DeleteBlogLike(blogId, visitId)
	if err != nil {
		utils.LogError(c, "Error deleting blog like: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error deleting blog like",
		})
	}

	utils.LogInfo(c, "Blog like deleted successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Blog like deleted successfully",
	})
}
