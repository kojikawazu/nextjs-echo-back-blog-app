package handlers_auth

import (
	"backend/apperrors"
	"backend/config"
	"backend/models"
	utils_cookie "backend/utils/cookie"
	utils "backend/utils/log"

	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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
	claims := &models.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
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
//   - error: トークン不在は apperrors.ErrAuthTokenMissing、解析失敗・無効は ErrAuthTokenInvalid。成功時はユーザー情報を 200 で返す（エラーのステータス・本文は middlewares.HTTPErrorHandler が決める）
func (h *AuthHandler) CheckAuth(c echo.Context) error {
	utils.LogInfo(c, "Checking authentication...")

	// クッキーからJWTトークンを取得
	cookie, err := c.Cookie("token")
	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenMissing, err)
	}
	tokenString := cookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		return apperrors.Wrap(apperrors.ErrAuthTokenInvalid, err)
	}

	if !token.Valid {
		return apperrors.ErrAuthTokenInvalid
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
