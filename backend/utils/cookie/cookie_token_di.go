package utils_cookie

import (
	"backend/config"
	"backend/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// CreateToken - JWTトークンを作成
//
// 引数:
//   - user: トークンに埋め込むユーザー情報
//
// 戻り値:
//   - string: 生成したJWTトークン文字列
//   - error: トークンの署名に失敗した場合のエラー
func (u *CookieUtilsImpl) CreateToken(user *models.BlogUsersData) (string, error) {
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
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存するJWTトークン文字列
//   - expirationTime: Cookieの有効期限
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
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存する更新後のJWTトークン文字列
//   - expirationTime: Cookieの有効期限
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
//
// 引数:
//   - c: Echoのリクエストコンテキスト
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

// GetUserIdFromToken - JWTトークンを解析してユーザーIDを取得
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: 解析対象のJWTトークン文字列
//
// 戻り値:
//   - string: トークンから取得したユーザーID
//   - error: トークンが無効または有効期限切れの場合のエラー
func (u *CookieUtilsImpl) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	claims, err := u.VerifyToken(c, tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
