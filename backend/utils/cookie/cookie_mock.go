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
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenName: 取得対象のCookie名
//
// 戻り値:
//   - *http.Cookie: モックの戻り値設定に従うCookie
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error) {
	args := m.Called(c, tokenName)

	// nil チェックを行い、安全にキャストする
	if args.Get(0) != nil {
		return args.Get(0).(*http.Cookie), args.Error(1)
	}

	return nil, args.Error(1)
}

// GetAuthCookieValue は GetAuthCookieValue のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenName: 取得対象のCookie名
//
// 戻り値:
//   - string: モックの戻り値設定に従うCookieの値
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) GetAuthCookieValue(c echo.Context, tokenName string) (string, error) {
	args := m.Called(c, tokenName)
	return args.String(0), args.Error(1)
}

// GetAuthCookieExpirationTime は GetAuthCookieExpirationTime のテスト用モック。
//
// 戻り値:
//   - time.Time: モックの戻り値設定に従う有効期限
func (m *MockCookieUtils) GetAuthCookieExpirationTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

// ExistsAuthCookie は ExistsAuthCookie のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenName: 確認対象のCookie名
//
// 戻り値:
//   - bool: モックの戻り値設定に従う存在有無
func (m *MockCookieUtils) ExistsAuthCookie(c echo.Context, tokenName string) bool {
	args := m.Called(c, tokenName)
	return args.Bool(0)
}

// VerifyToken は VerifyToken のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: 検証対象のJWTトークン文字列
//
// 戻り値:
//   - *models.Claims: モックの戻り値設定に従うクレーム情報
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) VerifyToken(c echo.Context, tokenString string) (*models.Claims, error) {
	args := m.Called(c, tokenString)
	return args.Get(0).(*models.Claims), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// 認証Token用
// ----------------------------------------------------------------------------------------------------------

// CreateToken は CreateToken のテスト用モック。
//
// 引数:
//   - user: トークンに埋め込むユーザー情報
//
// 戻り値:
//   - string: モックの戻り値設定に従うJWTトークン文字列
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) CreateToken(user *models.BlogUsersData) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// AddAuthCookie は AddAuthCookie のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存するJWTトークン文字列
//   - expirationTime: Cookieの有効期限
func (m *MockCookieUtils) AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// UpdateAuthCookie は UpdateAuthCookie のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存する更新後のJWTトークン文字列
//   - expirationTime: Cookieの有効期限
func (m *MockCookieUtils) UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// DelAuthCookie は DelAuthCookie のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
func (m *MockCookieUtils) DelAuthCookie(c echo.Context) {
	m.Called(c)
}

// GetUserIdFromToken は GetUserIdFromToken のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: 解析対象のJWTトークン文字列
//
// 戻り値:
//   - string: モックの戻り値設定に従うユーザーID
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) GetUserIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}

// ----------------------------------------------------------------------------------------------------------
// VisitId用
// ----------------------------------------------------------------------------------------------------------

// CreateVisitIdToken は CreateVisitIdToken のテスト用モック。
//
// 戻り値:
//   - string: モックの戻り値設定に従うvisitId用JWTトークン文字列
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) CreateVisitIdToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// AddVisitIdCoookie は AddVisitIdCoookie のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: Cookieに保存するvisitId用JWTトークン文字列
//   - expirationTime: Cookieの有効期限
func (m *MockCookieUtils) AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time) {
	m.Called(c, tokenString, expirationTime)
}

// GetVisitIdFromToken は GetVisitIdFromToken のテスト用モック。
//
// 引数:
//   - c: Echoのリクエストコンテキスト
//   - tokenString: 解析対象のvisitId用JWTトークン文字列
//
// 戻り値:
//   - string: モックの戻り値設定に従う訪問者ID
//   - error: モックの戻り値設定に従うエラー
func (m *MockCookieUtils) GetVisitIdFromToken(c echo.Context, tokenString string) (string, error) {
	args := m.Called(c, tokenString)
	return args.String(0), args.Error(1)
}
