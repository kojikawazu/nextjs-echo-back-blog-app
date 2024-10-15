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

// AddAuthCookie - 認証用のCookieを追加
func (u *CookieUtilsImpl) AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime
	cookie.HttpOnly = true
	cookie.Path = "/"
	if config.IsProduction {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}
	c.SetCookie(cookie)
}

// DelAuthCookie - 認証用のCookieを削除
func (u *CookieUtilsImpl) DelAuthCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0) // 有効期限を過去に設定して削除
	cookie.HttpOnly = true
	cookie.Path = "/"
	if config.IsProduction {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}
	c.SetCookie(cookie)
}

// GetAuthCookie - 認証用のCookieを取得
func (u *CookieUtilsImpl) GetAuthCookie(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

// GetAuthCookieValue - 認証用のCookieの値を取得
func (u *CookieUtilsImpl) GetAuthCookieValue(c echo.Context) (string, error) {
	cookie, err := u.GetAuthCookie(c)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ExistsAuthCookie - 認証用のCookieが存在するか確認
func (u *CookieUtilsImpl) ExistsAuthCookie(c echo.Context) bool {
	_, err := u.GetAuthCookie(c)
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

// GetUserIdFromToken - JWTトークンを解析してユーザーIDを取得
func (u *CookieUtilsImpl) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	claims, err := u.VerifyToken(c, tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
