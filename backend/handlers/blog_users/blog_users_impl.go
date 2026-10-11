package handlers_blog_users

import (
	"backend/apperrors"
	"backend/dto"
	utils "backend/utils/log"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FetchBlogUsers はクッキーの JWT トークンから解決したユーザー ID でユーザーデータを取得するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキーを含む）
//
// 戻り値:
//   - error: 認証トークンの不在・不正は apperrors.ErrAuthTokenMissing / ErrAuthTokenInvalid、ユーザー不在は ErrUserNotFound、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogUsersHandler) FetchBlogUsers(c echo.Context) error {
	utils.LogInfo(c, "Fetching user...")

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

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserById(userId)
	if err != nil {
		return err
	}

	utils.LogInfo(c, "Fetched user successfully")
	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

// UpdateBlogUsers はユーザーデータを更新し、認証クッキーを再発行するエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキーと name/email/password/newPassword を含む JSON ボディ）
//
// 戻り値:
//   - error: 認証トークンの不在・不正は apperrors.ErrAuthTokenMissing / ErrAuthTokenInvalid、不正入力は ErrXxx、現在のパスワード不一致は ErrInvalidCurrentPassword、ユーザー不在は ErrUserNotFound、その他の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *BlogUsersHandler) UpdateBlogUsers(c echo.Context) error {
	utils.LogInfo(c, "Updating user...")

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

	// JSONのリクエストボディからname, email, passwordを取得
	type RequestBody struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		return apperrors.Wrap(apperrors.ErrInvalidRequestBody, err)
	}

	// サービス層からユーザーデータを更新
	user, err := h.UserService.UpdateUser(userId, reqBody.Name, reqBody.Email, reqBody.Password, reqBody.NewPassword)
	if err != nil {
		return err
	}

	// トークンの有効期限を1時間に設定
	expirationTime := h.CookieUtils.GetAuthCookieExpirationTime()
	// トークンの再作成
	tokenString, err := h.CookieUtils.CreateToken(user)
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	// トークンをクッキーに更新
	h.CookieUtils.UpdateAuthCookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "Updated user successfully")
	return c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
