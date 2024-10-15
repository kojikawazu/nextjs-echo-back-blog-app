package handlers_users

import (
	utils_cookie "backend/utils/cookie"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SetMockBlogCookies は、ブログ関連のクッキーを設定します
func SetMockBlogCookies(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils) {
	// JWT の署名キーを設定し、正しいトークンを生成
	token := "mocked-token"
	validUserId := "valid-user-id"

	// リクエストにクッキーを追加
	cookie := &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)

	// モックの振る舞いを設定
	mockCookieUtils.On("GetAuthCookieValue", c).Return(token, nil)
	mockCookieUtils.On("GetUserIdFromToken", c, token).Return(validUserId, nil)
}
