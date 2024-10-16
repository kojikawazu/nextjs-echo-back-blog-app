package utils_cookie

import (
	"backend/config"
	"backend/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// CreateToken - JWTトークンを作成
func (u *CookieUtilsImpl) CreateToken(user *models.UserData) (string, error) {
	// トークンの有効期限を1時間に設定
	expirationTime := u.GetAuthCookieExpirationTime()
	// JWTトークンの作成
	claims := &models.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// トークンを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// トークンを文字列に変換
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		log.Println("Could not create JWT token: " + err.Error())
		return "", err
	}

	log.Println("JWT token created successfully")
	return tokenString, nil
}

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

// UpdateAuthCookie - 認証用のCookieを更新
func (u *CookieUtilsImpl) UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
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

// GetAuthCookieExpirationTime - 認証用のCookieの有効期限を取得
func (u *CookieUtilsImpl) GetAuthCookieExpirationTime() time.Time {
	return time.Now().Add(1 * time.Hour)
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
