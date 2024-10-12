package handlers_auth

import (
	"backend/config"
	"backend/models"
	utils_cookie "backend/utils/cookie"
	utils "backend/utils/log"

	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// ログインエンドポイント（JWTトークンの発行）
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
		utils.LogError(c, "Failed to bind request body: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// バリデーション
	err := h.AuthService.Login(reqBody.Email, reqBody.Password)
	if err != nil {
		switch err.Error() {
		case "email and password are required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Email and password are required",
			})
		case "invalid email format":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "An error occurred",
			})
		}
	}

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserByEmailAndPassword(reqBody.Email, reqBody.Password)
	if err != nil {
		utils.LogError(c, "Error fetching user: "+err.Error())
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found",
		})
	}

	// 認証成功
	utils.LogInfo(c, "User authenticated successfully:"+user.Email)

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
		utils.LogError(c, "Could not create JWT token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Could not create token",
		})
	}
	utils.LogInfo(c, "JWT token created successfully")

	// HTTPS-onlyクッキーにトークンをセット
	utils_cookie.AddAuthCookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "JWT token set in HTTPS-only cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

// 認証確認エンドポイント
func (h *AuthHandler) CheckAuth(c echo.Context) error {
	utils.LogInfo(c, "Checking authentication...")

	// クッキーからJWTトークンを取得
	cookie, err := c.Cookie("token")
	if err != nil {
		utils.LogError(c, "Token not found in cookies")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token not found"})
	}
	tokenString := cookie.Value

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		utils.LogError(c, "Failed to parse token: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	if !token.Valid {
		utils.LogError(c, "Invalid token")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	// 認証成功
	utils.LogInfo(c, "Authentication successful for user: "+claims.Email)
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "Authenticated",
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
	})
}

// ログアウトエンドポイント
func (h *AuthHandler) Logout(c echo.Context) error {
	utils.LogInfo(c, "Logging out...")

	// クッキーを削除するために、空のトークンと過去の有効期限を設定
	utils_cookie.DelAuthCookie(c)

	utils.LogInfo(c, "User logged out and token removed from cookie")
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}
