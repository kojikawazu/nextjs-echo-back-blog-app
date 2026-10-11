package handlers_blogs

import (
	"backend/apperrors"
	"backend/dto"
	services_blogs "backend/services/blogs"
	utils "backend/utils/log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// FetchBlogs は全ブログデータを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト
//
// 戻り値:
//   - error: 取得失敗時はエラー。成功時は 200 でブログ一覧を返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogs(c echo.Context) error {
	utils.LogInfo(c, "Fetching blogs...")

	// サービス層から全ブログデータを取得
	blogs, err := h.BlogService.FetchBlogs()
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched blogs successfully")
	return c.JSON(http.StatusOK, dto.NewBlogResponses(blogs))
}

// FetchBlogsByUserId はユーザー ID でブログデータを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（パスパラメータ userId を含む）
//
// 戻り値:
//   - error: 不正な userId は ErrInvalidUserID、その他の失敗はエラー。成功時は 200 でブログ一覧を返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogsByUserId(c echo.Context) error {
	utils.LogInfo(c, "Fetching blogs by userId...")

	// パスパラメータからuserIdを取得
	userId := c.Param("userId")

	// サービス層からユーザーIDでブログデータを取得
	blogs, err := h.BlogService.FetchBlogsByUserId(userId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched blogs successfully")
	return c.JSON(http.StatusOK, dto.NewBlogResponses(blogs))
}

// FetchBlogById はブログ ID でブログデータを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（パスパラメータ id を含む）
//
// 戻り値:
//   - error: 不正な id は ErrInvalidID、不在は ErrBlogNotFound、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogById(c echo.Context) error {
	utils.LogInfo(c, "Fetching blog by id...")

	// パスパラメータからidを取得
	id := c.Param("id")

	// サービス層からIDでブログデータを取得
	blog, err := h.BlogService.FetchBlogById(id)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched blog successfully")
	return c.JSON(http.StatusOK, dto.NewBlogResponse(blog))
}

// CreateBlog はブログデータを作成するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキーと title/githubUrl/category/description/tags を含む JSON ボディ）
//
// 戻り値:
//   - error: 認証トークンの不在・不正は apperrors.ErrAuthTokenMissing / ErrAuthTokenInvalid、不正入力は ErrInvalidXxx、その他の失敗はエラー。成功時は 201 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) CreateBlog(c echo.Context) error {
	utils.LogInfo(c, "Creating blog...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "token")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenMissing, err)
	}

	// JWTトークンを解析してユーザーIDを取得
	userId, err := h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenInvalid, err)
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
		return apperrors.Wrap(apperrors.ErrInvalidRequestBody, err)
	}

	// サービス層からブログデータを作成
	blog, err := h.BlogService.CreateBlog(userId, req.Title, req.GitHubURL, req.Category, req.Description, req.Tags)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Created blog successfully")
	return c.JSON(http.StatusCreated, dto.NewBlogResponse(blog))
}

// UpdateBlog はブログデータを更新するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキー、パスパラメータ id、更新内容を含む JSON ボディ）
//
// 戻り値:
//   - error: 認証トークンの不在・不正は apperrors.ErrAuthTokenMissing / ErrAuthTokenInvalid、不正入力は ErrInvalidXxx、不在は ErrBlogNotFound、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) UpdateBlog(c echo.Context) error {
	utils.LogInfo(c, "Updating blog...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "token")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenMissing, err)
	}

	// JWTトークンを解析してユーザーIDを取得
	_, err = h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenInvalid, err)
	}

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
		return apperrors.Wrap(apperrors.ErrInvalidRequestBody, err)
	}

	// サービス層からブログデータを更新
	blog, err := h.BlogService.UpdateBlog(id, req.Title, req.GitHubURL, req.Category, req.Description, req.Tags)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Updated blog successfully")
	return c.JSON(http.StatusOK, dto.NewBlogResponse(blog))
}

// DeleteBlog はブログデータを削除するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキーとパスパラメータ id を含む）
//
// 戻り値:
//   - error: 認証トークンの不在・不正は apperrors.ErrAuthTokenMissing / ErrAuthTokenInvalid、不正な id は ErrInvalidID、その他の失敗はエラー。成功時は 204 を返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) DeleteBlog(c echo.Context) error {
	utils.LogInfo(c, "Deleting blog...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "token")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenMissing, err)
	}

	// JWTトークンを解析してユーザーIDを取得
	_, err = h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenInvalid, err)
	}

	// パスパラメータからidを取得
	id := c.Param("id")

	// サービス層からブログデータを削除
	err = h.BlogService.DeleteBlog(id)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Deleted blog successfully")
	return c.NoContent(http.StatusNoContent)
}

// FetchBlogCategories はブログカテゴリーを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト
//
// 戻り値:
//   - error: 取得失敗時はエラー。成功時は 200 でカテゴリ一覧を返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogCategories(c echo.Context) error {
	utils.LogInfo(c, "Fetching categories...")

	// サービス層からカテゴリーを取得
	categories, err := h.BlogService.FetchBlogCategories()
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched categories successfully")
	return c.JSON(http.StatusOK, categories)
}

// FetchBlogTags はブログタグを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト
//
// 戻り値:
//   - error: 取得失敗時はエラー。成功時は 200 でタグ一覧を返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogTags(c echo.Context) error {
	utils.LogInfo(c, "Fetching tags...")

	// サービス層からタグを取得
	tags, err := h.BlogService.FetchBlogTags()
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched tags successfully")
	return c.JSON(http.StatusOK, tags)
}

// FetchBlogPopular は人気のあるブログを指定件数取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（パスパラメータ count を含む）
//
// 戻り値:
//   - error: count が数値でない・0 以下の場合は ErrInvalidCount、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogHandler) FetchBlogPopular(c echo.Context) error {
	utils.LogInfo(c, "Fetching popular blogs...")

	// パスパラメータからcountを取得
	count := c.Param("count")
	countInt, err := strconv.Atoi(count)
	if err != nil {
		return apperrors.Wrap(services_blogs.ErrInvalidCount, err)
	}

	// サービス層から人気のあるブログを取得
	blogs, err := h.BlogService.FetchBlogPopular(countInt)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched popular blogs successfully")
	return c.JSON(http.StatusOK, dto.NewBlogResponses(blogs))
}
