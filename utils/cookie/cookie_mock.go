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

// ----------------------------------------------------------------------------------------------------------
// Common用
// ----------------------------------------------------------------------------------------------------------

func (m *MockCookieUtils) GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error) {
	args := m.Called(c, tokenName)

	// nil チェックを行い、安全にキャストする
	if args.Get(0) != nil {
		return args.Get(0).(*http.Cookie), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockCookieUtils) GetAuthCookieValue(c echo.Context, tokenName string) (string, error) {
	args := m.Called(c, tokenName)
	return args.String(0), args.Error(1)
}

func (m *MockCookieUtils) GetAuthCookieExpirationTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockCookieUtils) ExistsAuthCookie(c echo.Context, tokenName string) bool {
	args := m.Called(c, tokenName)
	return args.Bool(0)
}

func (m *MockCookieUtils) VerifyToken(c echo.Context, tokenString string) (*models.Claims, error) {
	args := m.Called(c, tokenString)
	return args.Get(0).(*models.Claims), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// 認証Token用
// ----------------------------------------------------------------------------------------------------------

func (m *MockCookieUtils) CreateToken(user *models.UserData) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockCookieUtils) AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

func (m *MockCookieUtils) UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

func (m *MockCookieUtils) DelAuthCookie(c echo.Context) {
	m.Called(c)
}

func (m *MockCookieUtils) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// VisitId用
// ----------------------------------------------------------------------------------------------------------

func (m *MockCookieUtils) CreateVisitIdToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockCookieUtils) AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

func (m *MockCookieUtils) GetVisitIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}
