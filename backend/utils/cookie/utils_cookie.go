package utils_cookie

import (
	"backend/config"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// 認証用のCookieを追加
func AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
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

// 認証用のCookieを削除
func DelAuthCookie(c echo.Context) {
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
