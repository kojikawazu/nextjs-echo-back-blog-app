package utils_cookie

import (
	"backend/config"
	"backend/models"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// GetAuthCookie - 認証用のCookieを取得
func (u *CookieUtilsImpl) GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error) {
	cookie, err := c.Cookie(tokenName)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

// GetAuthCookieValue - 認証用のCookieの値を取得
func (u *CookieUtilsImpl) GetAuthCookieValue(c echo.Context, tokenName string) (string, error) {
	cookie, err := u.GetAuthCookie(c, tokenName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetAuthCookieExpirationTime - 認証用のCookieの有効期限を取得
func (u *CookieUtilsImpl) GetAuthCookieExpirationTime() time.Time {
	return time.Now().Add(1 * time.Hour)
}

// ExistsAuthCookie - 認証用のCookieが存在するか確認
func (u *CookieUtilsImpl) ExistsAuthCookie(c echo.Context, tokenName string) bool {
	_, err := u.GetAuthCookie(c, tokenName)
	return err == nil
}

// VerifyToken - JWTトークンを検証
func (u *CookieUtilsImpl) VerifyToken(c echo.Context, tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// トークンの有効期限をチェック (Unixタイムスタンプとして比較)
	expirationTime := time.Unix(claims.ExpiresAt, 0)
	if !token.Valid || expirationTime.Before(time.Now()) {
		return nil, errors.New("token expired or invalid")
	}

	return claims, nil
}
