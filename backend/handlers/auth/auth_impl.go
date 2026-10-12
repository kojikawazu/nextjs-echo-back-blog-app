package handlers_auth

import (
	"backend/apperrors"
	utils_cookie "backend/utils/cookie"
	utils "backend/utils/log"

	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Login はメール/パスワードで認証し JWT トークンを発行するログインエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（email/password を含む JSON ボディ）
//
// 戻り値:
//   - error: ボディ不正は apperrors.ErrInvalidRequestBody、入力不正は ErrXxx、認証失敗は ErrInvalidCredentials、トークン生成等の失敗はエラー。成功時は 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *AuthHandler) Login(c echo.Context) error {
	utils.LogInfo(c, "Logging in...")

	// JSONのリクエストボディからemailとpasswordを取得
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		return apperrors.Wrap(apperrors.ErrInvalidRequestBody, err)
	}

	// バリデーション
	err := h.AuthService.Login(reqBody.Email, reqBody.Password)
	if err != nil {
		return err
	}

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserByEmailAndPassword(reqBody.Email, reqBody.Password)
	if err != nil {
		return err
	}

	// 認証成功
	utils.LogInfo(c, "User authenticated successfully: id="+user.ID)

	// JWTトークンの作成
	expirationTime := time.Now().Add(1 * time.Hour) // トークンの有効期限を1時間に設定
	tokenString, err := utils_cookie.NewAuthToken(user, expirationTime)
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	utils.LogInfo(c, "JWT token created successfully")

	// HTTPS-onlyクッキーにトークンをセット
	utils_cookie.AddAuthCookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "JWT token set in HTTPS-only cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

// CheckAuth はクッキーの JWT トークンを検証する認証確認エンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト（token クッキーを含む）
//
// 戻り値:
//   - error: トークン不在は apperrors.ErrAuthTokenMissing、形式不正・署名不正・期限切れ・種類違い（訪問者トークン等）・user_id 空は ErrAuthTokenInvalid。成功時はユーザー情報を 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *AuthHandler) CheckAuth(c echo.Context) error {
	utils.LogInfo(c, "Checking authentication...")

	// クッキーからJWTトークンを取得
	cookie, err := c.Cookie("token")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenMissing, err)
	}
	tokenString := cookie.Value

	// 署名・種類（aud）・期限・user_id を検証する
	claims, err := utils_cookie.ParseAuthToken(tokenString)
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenInvalid, err)
	}

	// 認証成功
	utils.LogInfo(c, "Authentication successful for user: id="+claims.UserID)
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Authenticated",
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}

// Logout は認証クッキーを削除するログアウトエンドポイント。
//
// 引数:
//   - c: リクエスト/レスポンスを保持する echo コンテキスト
//
// 戻り値:
//   - error: レスポンス書き込みエラー。成功時は 200 を返す
func (h *AuthHandler) Logout(c echo.Context) error {
	utils.LogInfo(c, "Logging out...")

	// クッキーを削除するために、空のトークンと過去の有効期限を設定
	utils_cookie.DelAuthCookie(c)

	utils.LogInfo(c, "User logged out and token removed from cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
