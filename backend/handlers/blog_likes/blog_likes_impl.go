package handlers_blog_likes

import (
	"backend/apperrors"
	"backend/dto"
	utils "backend/utils/log"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FetchBlogLikesByVisitId - 訪問IDに紐づくブログいいねデータを取得するハンドラ
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（visit-id-token クッキーを含む）
//
// 戻り値:
//   - error: 訪問者トークンの不在・不正は apperrors.ErrVisitorTokenMissing / ErrVisitorTokenInvalid、データ取得の失敗はエラー（500）。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogLikeHandler) FetchBlogLikesByVisitId(c echo.Context) error {
	utils.LogInfo(c, "Fetching blog likes by visit id...")

	// クッキーの訪問者トークンから訪問IDを取得
	visitId, err := h.visitIdFromCookie(c)
	if err != nil {
		return err
	}

	// VisitIDに紐づくいいねデータを取得
	blogLikesData, err := h.BlogLikeService.FetchBlogLikesByVisitId(visitId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Blog likes fetched successfully")
	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, dto.NewBlogLikeResponses(blogLikesData))
}

// GenerateVisitorId -　訪問者IDを生成するハンドラ
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（既存の visit-id-token クッキーがあれば利用）
//
// 戻り値:
//   - error: トークン生成の失敗はエラー（500）。生成成功または既存時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
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
		return fmt.Errorf("create visitor token: %w", err)
	}
	//  Visit用トークンをクッキーに保存
	h.CookieUtils.AddVisitIdCoookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "Visitor id generated successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Visitor id generated successfully",
	})
}

// IsBlogLiked は対象ブログにいいね済みかを確認するハンドラ。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（visit-id-token クッキーとパスパラメータ blogId を含む）
//
// 戻り値:
//   - error: 訪問者トークンの不在・不正は apperrors.ErrVisitorTokenMissing / ErrVisitorTokenInvalid、確認の失敗はエラー（500）、blogId が空は ErrBlogIDOrVisitIDEmpty。成功時は isLiked を 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogLikeHandler) IsBlogLiked(c echo.Context) error {
	utils.LogInfo(c, "Checking if blog is liked...")

	// クッキーの訪問者トークンから訪問IDを取得
	visitId, err := h.visitIdFromCookie(c)
	if err != nil {
		return err
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータが存在するか確認
	isLiked, err := h.BlogLikeService.IsBlogLiked(blogId, visitId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Checked blog like")
	return c.JSON(http.StatusOK, map[string]bool{
		"isLiked": isLiked,
	})
}

// CreateBlogLike はブログいいねを追加するハンドラ。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（visit-id-token クッキーとパスパラメータ blogId を含む）
//
// 戻り値:
//   - error: 訪問者トークンの不在・不正は apperrors.ErrVisitorTokenMissing / ErrVisitorTokenInvalid、空値は ErrBlogIDOrVisitIDEmpty、いいね済みは ErrAlreadyLiked。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogLikeHandler) CreateBlogLike(c echo.Context) error {
	utils.LogInfo(c, "Creating blog like...")

	// クッキーの訪問者トークンから訪問IDを取得
	visitId, err := h.visitIdFromCookie(c)
	if err != nil {
		return err
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータを作成
	createdBlogLikeData, err := h.BlogLikeService.CreateBlogLike(blogId, visitId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Blog like created successfully")
	return c.JSON(http.StatusOK, dto.NewBlogLikeResponse(createdBlogLikeData))
}

// DeleteBlogLike はブログいいねを削除するハンドラ。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（visit-id-token クッキーとパスパラメータ blogId を含む）
//
// 戻り値:
//   - error: 訪問者トークンの不在・不正は apperrors.ErrVisitorTokenMissing / ErrVisitorTokenInvalid、削除の失敗はエラー（500）。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogLikeHandler) DeleteBlogLike(c echo.Context) error {
	utils.LogInfo(c, "Deleting blog like...")

	// クッキーの訪問者トークンから訪問IDを取得
	visitId, err := h.visitIdFromCookie(c)
	if err != nil {
		return err
	}

	// パスパラメータからブログIDと訪問者IDを取得
	blogId := c.Param("blogId")

	// いいねデータを削除
	err = h.BlogLikeService.DeleteBlogLike(blogId, visitId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Blog like deleted successfully")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Blog like deleted successfully",
	})
}

// visitIdFromCookie は visit-id-token クッキーの JWT を解析して訪問IDを取り出す。
// 不在・不正はクライアント側の要因（未発行・期限切れ・改ざん）のため、500 ではなく 401 になる sentinel を返す。
//
// 引数:
//   - c: リクエストを保持する echo コンテキスト（visit-id-token クッキーを含む）
//
// 戻り値:
//   - string: 訪問ID
//   - error: クッキー不在は apperrors.ErrVisitorTokenMissing、解析失敗・期限切れは ErrVisitorTokenInvalid（いずれも原因を添える）
func (h *BlogLikeHandler) visitIdFromCookie(c echo.Context) (string, error) {
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "visit-id-token")
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrVisitorTokenMissing, err)
	}
	visitId, err := h.CookieUtils.GetVisitIdFromToken(c, cookieValue)
	if err != nil {
		return "", apperrors.Wrap(apperrors.ErrVisitorTokenInvalid, err)
	}
	return visitId, nil
}
