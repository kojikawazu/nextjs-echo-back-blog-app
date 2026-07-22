package handlers_blog_users

import (
	utils_cookie "backend/utils/cookie"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SetMockBlogCookies は、ブログ関連のクッキーを設定します
//
// 引数:
//   - c: モックの期待値照合に用いる echo コンテキスト
//   - req: クッキーを追加する対象の HTTP リクエスト
//   - mockCookieUtils: 振る舞いを設定するクッキーユーティリティのモック
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
	mockCookieUtils.On("GetAuthCookieValue", c, "token").Return(token, nil)
	mockCookieUtils.On("GetUserIdFromToken", c, token).Return(validUserId, nil)
}
