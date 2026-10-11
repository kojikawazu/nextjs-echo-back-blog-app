package handlers_blog_users

import (
	"backend/models"
	services_users "backend/services/blog_users"
	utils_cookie "backend/utils/cookie"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const responseTestPassword = "plain-secret-password"

var publicUserKeys = []string{"id", "name", "email", "created_at", "updated_at"}

// assertUserResponse はレスポンスが公開フィールドのみのユーザー情報であることを検証する。
func assertUserResponse(t *testing.T, body string) {
	t.Helper()
	var m map[string]any
	require.NoError(t, json.Unmarshal([]byte(body), &m))

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, publicUserKeys, keys)
	assert.Equal(t, "valid-user-id", m["id"])
	assert.Equal(t, "John Doe", m["name"])
	assert.NotContains(t, body, responseTestPassword)
}

func newUserWithPassword() *models.BlogUsersData {
	return &models.BlogUsersData{
		ID:       "valid-user-id",
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: responseTestPassword,
	}
}

// 正常系: ユーザー取得のレスポンスは DTO の公開フィールドのみで、password キーを含まない。
func TestHandler_FetchUser_ResponseExcludesPassword(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/users/detail", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	mockUserService.On("FetchUserById", "valid-user-id").Return(newUserWithPassword(), nil)
	SetMockBlogCookies(c, req, mockCookieUtils)

	assert.NoError(t, NewBlogUsersHandler(mockUserService, mockCookieUtils).FetchBlogUsers(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assertUserResponse(t, rec.Body.String())
}

// 正常系: ユーザー更新のレスポンスは DTO の公開フィールドのみで、password キーを含まない。
func TestHandler_UpdateUser_ResponseExcludesPassword(t *testing.T) {
	e := echo.New()
	body := `{"name":"John Doe","email":"john@example.com","password":"current","newPassword":"` + responseTestPassword + `"}`
	req := httptest.NewRequest(http.MethodPut, "/api/users/update", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := newUserWithPassword()
	mockCookieUtils := new(utils_cookie.MockCookieUtils)
	mockUserService := new(services_users.MockUserService)
	mockUserService.On("UpdateUser", "valid-user-id", "John Doe", "john@example.com", "current", responseTestPassword).Return(user, nil)
	SetMockBlogCookies(c, req, mockCookieUtils)
	mockCookieUtils.On("GetAuthCookieExpirationTime").Return(time.Now().Add(time.Hour))
	mockCookieUtils.On("CreateToken", user).Return("new-mocked-token", nil)
	mockCookieUtils.On("UpdateAuthCookie", c, "new-mocked-token", mock.Anything).Return(nil)

	assert.NoError(t, NewBlogUsersHandler(mockUserService, mockCookieUtils).UpdateBlogUsers(c))

	assert.Equal(t, http.StatusOK, rec.Code)
	assertUserResponse(t, rec.Body.String())
}
