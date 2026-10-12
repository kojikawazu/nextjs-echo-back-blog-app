//go:build e2e

package e2e

import (
	"backend/config"
	"backend/models"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertErrorResponse はステータスと統一エラーレスポンス（dto.ErrorResponse）の code・error を検証し、ボディを Close する。
func assertErrorResponse(t *testing.T, resp *http.Response, wantStatus int, wantCode, wantMessage string) {
	t.Helper()
	assert.Equal(t, wantStatus, resp.StatusCode)
	var body map[string]any
	decodeBody(t, resp, &body)
	// code・error 以外のキー（内部情報）を含まない
	assert.Equal(t, map[string]any{"code": wantCode, "error": wantMessage}, body)
}

// doRawJSON は JSON として不正な本文をそのまま送る（doJSON は本文をエンコードするため不正な JSON を送れない）。
func doRawJSON(t *testing.T, client *http.Client, method, path, rawBody string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, baseURL+path, strings.NewReader(rawBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", nextClientIP())
	resp, err := client.Do(req)
	require.NoError(t, err)
	return resp
}

// 準正常系: 存在しないルートも統一形式の 404 を返す（Echo 既定の {"message": ...} にならない）。
func TestE2E_Error_UnknownRoute(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodGet, "/api/no-such-route", nil)

	assertErrorResponse(t, resp, http.StatusNotFound, "NOT_FOUND", "Not Found")
}

// 準正常系: 入力不正・認証エラー・不在が、それぞれのステータスと code で返る。
func TestE2E_Error_ClientErrors(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		body        any
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name:        "blog_detail_not_found",
			method:      http.MethodGet,
			path:        "/api/blogs/detail/00000000-0000-0000-0000-000000000000",
			wantStatus:  http.StatusNotFound,
			wantCode:    "BLOG_NOT_FOUND",
			wantMessage: "Blog not found",
		},
		{
			// UUID 形式でない ID は DB に渡さず 400（以前は 404）
			name:        "blog_detail_malformed_id",
			method:      http.MethodGet,
			path:        "/api/blogs/detail/not-a-uuid",
			wantStatus:  http.StatusBadRequest,
			wantCode:    "VALIDATION_ERROR",
			wantMessage: "Invalid id",
		},
		{
			// 0 件指定は 400（以前は 500）
			name:        "popular_zero_count",
			method:      http.MethodGet,
			path:        "/api/blogs/popular/0",
			wantStatus:  http.StatusBadRequest,
			wantCode:    "VALIDATION_ERROR",
			wantMessage: "Invalid count",
		},
		{
			name:        "login_missing_fields",
			method:      http.MethodPost,
			path:        "/api/users/login",
			body:        map[string]string{},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "VALIDATION_ERROR",
			wantMessage: "Email and password are required",
		},
		{
			name:        "login_wrong_password",
			method:      http.MethodPost,
			path:        "/api/users/login",
			body:        map[string]string{"email": os.Getenv("TEST_USER_EMAIL"), "password": "wrong-password"},
			wantStatus:  http.StatusUnauthorized,
			wantCode:    "INVALID_CREDENTIALS",
			wantMessage: "Invalid credentials",
		},
		{
			// 以前は {"message": "Token not found"} だった
			name:        "auth_check_without_cookie",
			method:      http.MethodGet,
			path:        "/api/users/auth-check",
			wantStatus:  http.StatusUnauthorized,
			wantCode:    "AUTH_TOKEN_MISSING",
			wantMessage: "Authentication token is missing",
		},
		{
			name:        "create_blog_without_cookie",
			method:      http.MethodPost,
			path:        "/api/blogs/create",
			body:        map[string]string{"title": "t"},
			wantStatus:  http.StatusUnauthorized,
			wantCode:    "AUTH_TOKEN_MISSING",
			wantMessage: "Authentication token is missing",
		},
		{
			name:        "create_comment_missing_comment",
			method:      http.MethodPost,
			path:        "/api/comments/create",
			body:        map[string]string{"blogId": os.Getenv("TEST_BLOG_ID"), "guestUser": "E2E Guest"},
			wantStatus:  http.StatusBadRequest,
			wantCode:    "VALIDATION_ERROR",
			wantMessage: "Invalid comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newClient(t)

			resp := doJSON(t, client, tt.method, tt.path, tt.body)

			assertErrorResponse(t, resp, tt.wantStatus, tt.wantCode, tt.wantMessage)
		})
	}
}

// 準正常系: JSON として不正な本文は 400 INVALID_REQUEST_BODY で、パースエラーの詳細を返さない。
func TestE2E_Error_InvalidRequestBody(t *testing.T) {
	client := newClient(t)

	resp := doRawJSON(t, client, http.MethodPost, "/api/comments/create", "{")

	assertErrorResponse(t, resp, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
}

// 準正常系: 存在しないブログの更新は 404（以前は 500）。
func TestE2E_Error_UpdateBlogNotFound(t *testing.T) {
	client := newClient(t)
	login(t, client)

	resp := doJSON(t, client, http.MethodPut, "/api/blogs/update/00000000-0000-0000-0000-000000000000", map[string]string{
		"title":       "title",
		"githubUrl":   "https://github.com/example/repo",
		"category":    "category",
		"description": "description",
		"tags":        "tag",
	})

	assertErrorResponse(t, resp, http.StatusNotFound, "BLOG_NOT_FOUND", "Blog not found")
}

// 準正常系: 現在のパスワードが誤っているユーザー更新は 400（以前は 500）で、ユーザーは更新されない。
func TestE2E_Error_UpdateUserWrongCurrentPassword(t *testing.T) {
	client := newClient(t)
	login(t, client)

	resp := doJSON(t, client, http.MethodPut, "/api/users/update", map[string]string{
		"name":        "Changed Name",
		"email":       os.Getenv("TEST_USER_EMAIL"),
		"password":    "wrong-password",
		"newPassword": "new-password",
	})
	assertErrorResponse(t, resp, http.StatusBadRequest, "INVALID_CURRENT_PASSWORD", "Invalid current password")

	// 元のパスワードで引き続きログインできる（更新されていない）
	login(t, client)
}

// 準正常系: 同じ訪問者が同じブログに 2 回いいねすると 409（以前は 400）。
func TestE2E_Error_AlreadyLiked(t *testing.T) {
	client := newClient(t)
	blogID := os.Getenv("TEST_BLOG_ID")

	genResp := doJSON(t, client, http.MethodGet, "/api/blog-likes/generate-visit-id", nil)
	genResp.Body.Close()
	require.Equal(t, http.StatusOK, genResp.StatusCode)

	first := doJSON(t, client, http.MethodPost, "/api/blog-likes/create/"+blogID, nil)
	first.Body.Close()
	require.Equal(t, http.StatusOK, first.StatusCode)
	// 他のテストにいいねを残さない
	t.Cleanup(func() {
		resp := doJSON(t, client, http.MethodDelete, "/api/blog-likes/delete/"+blogID, nil)
		resp.Body.Close()
	})

	second := doJSON(t, client, http.MethodPost, "/api/blog-likes/create/"+blogID, nil)

	assertErrorResponse(t, second, http.StatusConflict, "ALREADY_LIKED", "Blog is already liked")
}

// visitorTokenEndpoints は visit-id-token Cookie を要求するいいね系エンドポイント。
var visitorTokenEndpoints = []struct {
	name   string
	method string
	path   func(blogID string) string
}{
	{name: "fetch_likes", method: http.MethodGet, path: func(string) string { return "/api/blog-likes" }},
	{name: "is_liked", method: http.MethodGet, path: func(id string) string { return "/api/blog-likes/is-liked/" + id }},
	{name: "create_like", method: http.MethodPost, path: func(id string) string { return "/api/blog-likes/create/" + id }},
	{name: "delete_like", method: http.MethodDelete, path: func(id string) string { return "/api/blog-likes/delete/" + id }},
}

// newClientWithVisitorToken は visit-id-token Cookie に token を持つクライアントを生成する。
func newClientWithVisitorToken(t *testing.T, token string) *http.Client {
	t.Helper()
	client := newClient(t)
	u, err := url.Parse(baseURL)
	require.NoError(t, err)
	client.Jar.SetCookies(u, []*http.Cookie{{Name: "visit-id-token", Value: token, Path: "/"}})
	return client
}

// signVisitorToken は key で署名した訪問者トークンを返す。
func signVisitorToken(t *testing.T, key []byte, expiresAt time.Time) string {
	t.Helper()
	claims := &models.ClaimsVisitId{
		VisitId:        "e2e-visitor",
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiresAt.Unix()},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	require.NoError(t, err)
	return token
}

// 準正常系: 訪問者トークンが無いいいね操作は 401 VISITOR_TOKEN_MISSING（#148 で 500 から変更）。
func TestE2E_Error_VisitorTokenMissing(t *testing.T) {
	blogID := os.Getenv("TEST_BLOG_ID")

	for _, ep := range visitorTokenEndpoints {
		t.Run(ep.name, func(t *testing.T) {
			client := newClient(t)

			resp := doJSON(t, client, ep.method, ep.path(blogID), nil)

			assertErrorResponse(t, resp, http.StatusUnauthorized, "VISITOR_TOKEN_MISSING", "Visitor token is missing")
		})
	}
}

// 準正常系: 訪問者トークンが不正・改ざん・期限切れのいいね操作は 401 VISITOR_TOKEN_INVALID（#148 で 500 から変更）。
func TestE2E_Error_VisitorTokenInvalid(t *testing.T) {
	blogID := os.Getenv("TEST_BLOG_ID")
	tokens := []struct {
		name  string
		token string
	}{
		{name: "malformed", token: "not-a-jwt"},
		// 別の鍵で署名（改ざん）
		{name: "wrong_signature", token: signVisitorToken(t, []byte("e2e-attacker-key"), time.Now().Add(time.Hour))},
		// 正しい鍵だが期限切れ
		{name: "expired", token: signVisitorToken(t, config.JwtKey, time.Now().Add(-time.Hour))},
	}

	for _, tk := range tokens {
		for _, ep := range visitorTokenEndpoints {
			t.Run(tk.name+"/"+ep.name, func(t *testing.T) {
				client := newClientWithVisitorToken(t, tk.token)

				resp := doJSON(t, client, ep.method, ep.path(blogID), nil)

				assertErrorResponse(t, resp, http.StatusUnauthorized, "VISITOR_TOKEN_INVALID", "Visitor token is invalid")
			})
		}
	}
}
