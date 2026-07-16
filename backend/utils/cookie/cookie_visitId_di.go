package utils_cookie

import (
	"backend/config"
	"backend/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CreateVisitIdToken - visitId用JWTトークンを作成
func (u *CookieUtilsImpl) CreateVisitIdToken() (string, error) {
	// 一意の訪問者IDを生成
	visitorID := uuid.New().String()

	// トークンの有効期限を1時間に設定
	expirationTime := u.GetAuthCookieExpirationTime()
	// JWTトークンの作成
	claims := &models.ClaimsVisitId{
		VisitId: visitorID,
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

// AddVisitIdCoookie - 訪問者IDを生成
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
func (u *CookieUtilsImpl) GetVisitIdFromToken(c echo.Context, tokenString string) (string, error) {
	claims := &models.ClaimsVisitId{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})

	if err != nil {
		return "", err
	}

	// トークンの有効期限をチェック (Unixタイムスタンプとして比較)
	expirationTime := time.Unix(claims.ExpiresAt, 0)
	if !token.Valid || expirationTime.Before(time.Now()) {
		return "", errors.New("token expired or invalid")
	}

	return claims.VisitId, nil
}
