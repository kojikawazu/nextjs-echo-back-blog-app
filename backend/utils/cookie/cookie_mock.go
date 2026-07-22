package utils_cookie

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

// MockCookieUtils は CookieUtils のテスト用モック。
type MockCookieUtils struct {
	mock.Mock
}

// ----------------------------------------------------------------------------------------------------------
// Common用
// ----------------------------------------------------------------------------------------------------------

// GetAuthCookie は GetAuthCookie のテスト用モック。
func (m *MockCookieUtils) GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error) {
	args := m.Called(c, tokenName)

	// nil チェックを行い、安全にキャストする
	if args.Get(0) != nil {
		return args.Get(0).(*http.Cookie), args.Error(1)
	}

	return nil, args.Error(1)
}

// GetAuthCookieValue は GetAuthCookieValue のテスト用モック。
func (m *MockCookieUtils) GetAuthCookieValue(c echo.Context, tokenName string) (string, error) {
	args := m.Called(c, tokenName)
	return args.String(0), args.Error(1)
}

// GetAuthCookieExpirationTime は GetAuthCookieExpirationTime のテスト用モック。
func (m *MockCookieUtils) GetAuthCookieExpirationTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

// ExistsAuthCookie は ExistsAuthCookie のテスト用モック。
func (m *MockCookieUtils) ExistsAuthCookie(c echo.Context, tokenName string) bool {
	args := m.Called(c, tokenName)
	return args.Bool(0)
}

// VerifyToken は VerifyToken のテスト用モック。
func (m *MockCookieUtils) VerifyToken(c echo.Context, tokenString string) (*models.Claims, error) {
	args := m.Called(c, tokenString)
	return args.Get(0).(*models.Claims), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// 認証Token用
// ----------------------------------------------------------------------------------------------------------

// CreateToken は CreateToken のテスト用モック。
func (m *MockCookieUtils) CreateToken(user *models.BlogUsersData) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// AddAuthCookie は AddAuthCookie のテスト用モック。
func (m *MockCookieUtils) AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// UpdateAuthCookie は UpdateAuthCookie のテスト用モック。
func (m *MockCookieUtils) UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// DelAuthCookie は DelAuthCookie のテスト用モック。
func (m *MockCookieUtils) DelAuthCookie(c echo.Context) {
	m.Called(c)
}

// GetUserIdFromToken は GetUserIdFromToken のテスト用モック。
func (m *MockCookieUtils) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// VisitId用
// ----------------------------------------------------------------------------------------------------------

// CreateVisitIdToken は CreateVisitIdToken のテスト用モック。
func (m *MockCookieUtils) CreateVisitIdToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// AddVisitIdCoookie は AddVisitIdCoookie のテスト用モック。
func (m *MockCookieUtils) AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// GetVisitIdFromToken は GetVisitIdFromToken のテスト用モック。
func (m *MockCookieUtils) GetVisitIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}
