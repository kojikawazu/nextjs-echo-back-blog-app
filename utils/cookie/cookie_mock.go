package utils_cookie

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockCookieUtils struct {
	mock.Mock
}

func (m *MockCookieUtils) AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

func (m *MockCookieUtils) DelAuthCookie(c echo.Context) {
	m.Called(c)
}

func (m *MockCookieUtils) GetAuthCookie(c echo.Context) (*http.Cookie, error) {
	args := m.Called(c)

	// nil チェックを行い、安全にキャストする
	if args.Get(0) != nil {
		return args.Get(0).(*http.Cookie), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockCookieUtils) GetAuthCookieValue(c echo.Context) (string, error) {
	args := m.Called(c)
	return args.String(0), args.Error(1)
}

func (m *MockCookieUtils) ExistsAuthCookie(c echo.Context) bool {
	args := m.Called(c)
	return args.Bool(0)
}

func (m *MockCookieUtils) VerifyToken(c echo.Context, tokenString string) (*models.Claims, error) {
	args := m.Called(c, tokenString)
	return args.Get(0).(*models.Claims), args.Error(1)
}

func (m *MockCookieUtils) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}
