package utils_cookie

import (
	"backend/config"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CreateVisitIdToken - visitId用JWTトークンを作成
//
// 戻り値:
//   - string: 生成したvisitId用JWTトークン文字列
//   - error: トークンの署名に失敗した場合のエラー
func (u *CookieUtilsImpl) CreateVisitIdToken() (string, error) {
	// 一意の訪問者IDを生成
	visitorID := uuid.New().String()

	// トークンの有効期限を1時間に設定
	expirationTime := u.GetAuthCookieExpirationTime()
	// JWTトークンの作成
	tokenString, err := NewVisitorToken(visitorID, expirationTime)
	if err != nil {
		log.Println("Could not create JWT token: " + err.Error())
		return "", err
	}

	log.Println("JWT token created successfully")
	return tokenString, nil
}

// AddVisitIdCoookie - 訪問者IDを生成
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存するvisitId用JWTトークン文字列
//   - expirationTime: Cookieの有効期限
func (u *CookieUtilsImpl) AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = "visit-id-token"
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

// GetVisitIdFromToken - 訪問者IDを取得
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: 解析対象のvisitId用JWTトークン文字列
//
// 戻り値:
//   - string: トークンから取得した訪問者ID
//   - error: 形式不正・署名不正・期限切れ・種類違い（認証トークン等）・visit_id 空の場合のエラー
func (u *CookieUtilsImpl) GetVisitIdFromToken(c echo.Context, tokenString string) (string, error) {
	claims, err := ParseVisitorToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.VisitId, nil
}
