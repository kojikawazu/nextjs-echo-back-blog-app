package utils_cookie

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// CookieUtils インターフェース
type CookieUtils interface {
	// Token用
	GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error)
	GetAuthCookieValue(c echo.Context, tokenName string) (string, error)
	GetAuthCookieExpirationTime() time.Time
	ExistsAuthCookie(c echo.Context, tokenName string) bool
	VerifyToken(c echo.Context, tokenString string) (*models.Claims, error)

	// 認証Token用
	CreateToken(user *models.UserData) (string, error)
	AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
	UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
	DelAuthCookie(c echo.Context)
	GetUserIdFromToken(c echo.Context, tokenString string) (string, error)

	// VisitId用
	CreateVisitIdToken() (string, error)
	AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time)
	GetVisitIdFromToken(c echo.Context, tokenString string) (string, error)
}

type CookieUtilsImpl struct{}

func NewCookieUtils() CookieUtils {
	return &CookieUtilsImpl{}
}
