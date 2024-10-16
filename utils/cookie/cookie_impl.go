package utils_cookie

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// CookieUtils インターフェース
type CookieUtils interface {
	CreateToken(user *models.UserData) (string, error)
	AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
	UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
	DelAuthCookie(c echo.Context)
	GetAuthCookie(c echo.Context) (*http.Cookie, error)
	GetAuthCookieValue(c echo.Context) (string, error)
	GetAuthCookieExpirationTime() time.Time
	ExistsAuthCookie(c echo.Context) bool
	VerifyToken(c echo.Context, tokenString string) (*models.Claims, error)
	GetUserIdFromToken(c echo.Context, tokenString string) (string, error)
}

type CookieUtilsImpl struct{}

func NewCookieUtils() CookieUtils {
	return &CookieUtilsImpl{}
}
