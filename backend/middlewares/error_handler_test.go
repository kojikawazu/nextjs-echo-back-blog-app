package middlewares

import (
	"backend/apperrors"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestHTTPErrorHandler はエラーが HTTP ステータスと統一エラーレスポンスに変換されることを検証する。
func TestHTTPErrorHandler(t *testing.T) {
	// サービス層の sentinel を模したエラー（apperrors の各 Kind の変換を検証する）
	errNotFound := apperrors.New(apperrors.KindNotFound, apperrors.CodeBlogNotFound, "Blog not found")
	errConflict := apperrors.New(apperrors.KindConflict, apperrors.CodeAlreadyLiked, "Blog is already liked")
	errInvalid := apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid title")
	errInternal := apperrors.New(apperrors.KindInternal, apperrors.CodeInternal, "Internal server error")

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			// 準正常系: 入力不正は 400
			name:       "invalid_returns_400",
			err:        errInvalid,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"VALIDATION_ERROR","error":"Invalid title"}`,
		},
		{
			// 準正常系: 未認証は 401
			name:       "unauthorized_returns_401",
			err:        apperrors.ErrAuthTokenMissing,
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":"AUTH_TOKEN_MISSING","error":"Authentication token is missing"}`,
		},
		{
			// 準正常系: 不在は 404
			name:       "not_found_returns_404",
			err:        errNotFound,
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"BLOG_NOT_FOUND","error":"Blog not found"}`,
		},
		{
			// 準正常系: 競合は 409
			name:       "conflict_returns_409",
			err:        errConflict,
			wantStatus: http.StatusConflict,
			wantBody:   `{"code":"ALREADY_LIKED","error":"Blog is already liked"}`,
		},
		{
			// 準正常系: レートリミット超過は 429
			name:       "too_many_requests_returns_429",
			err:        apperrors.ErrTooManyRequests,
			wantStatus: http.StatusTooManyRequests,
			wantBody:   `{"code":"TOO_MANY_REQUESTS","error":"Too many requests"}`,
		},
		{
			// 異常系: KindInternal は 500
			name:       "internal_kind_returns_500",
			err:        errInternal,
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":"INTERNAL_ERROR","error":"Internal server error"}`,
		},
		{
			// 準正常系: 原因を添えた sentinel もレスポンスには sentinel のメッセージだけを返す（原因を漏らさない）
			name:       "wrapped_sentinel_hides_cause",
			err:        apperrors.Wrap(apperrors.ErrInvalidRequestBody, errors.New("syntax error at offset 1")),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"INVALID_REQUEST_BODY","error":"Invalid request body"}`,
		},
		{
			// 準正常系: fmt.Errorf で包まれた sentinel も判定できる
			name:       "fmt_wrapped_sentinel",
			err:        fmt.Errorf("fetch blog: %w", errNotFound),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"BLOG_NOT_FOUND","error":"Blog not found"}`,
		},
		{
			// 準正常系: Echo のルート不在（404）はステータス名のコードで返す
			name:       "echo_not_found",
			err:        echo.ErrNotFound,
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"NOT_FOUND","error":"Not Found"}`,
		},
		{
			// 準正常系: Echo の 405 はステータス名のコードで返す
			name:       "echo_method_not_allowed",
			err:        echo.ErrMethodNotAllowed,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   `{"code":"METHOD_NOT_ALLOWED","error":"Method Not Allowed"}`,
		},
		{
			// 準正常系: echo.HTTPError の独自メッセージは返さず、ステータス名を返す（内部情報を漏らさない）
			name:       "echo_http_error_hides_message",
			err:        echo.NewHTTPError(http.StatusBadRequest, "Syntax error: offset=1, error=unexpected EOF"),
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"code":"BAD_REQUEST","error":"Bad Request"}`,
		},
		{
			// 異常系: echo.HTTPError の 5xx は汎用の 500 で返す
			name:       "echo_5xx_returns_internal",
			err:        echo.NewHTTPError(http.StatusServiceUnavailable, "db down"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":"INTERNAL_ERROR","error":"Internal server error"}`,
		},
		{
			// 異常系: 想定外のエラーは 500 で、エラー内容を本文に含めない
			name:       "unexpected_error_hides_detail",
			err:        errors.New("pq: password authentication failed for user postgres"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":"INTERNAL_ERROR","error":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			HTTPErrorHandler(tt.err, c)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}

// TestHTTPErrorHandler_Head は HEAD リクエストでは本文を返さずステータスのみを返すことを検証する（準正常系）。
func TestHTTPErrorHandler_Head(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodHead, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	HTTPErrorHandler(apperrors.ErrAuthTokenMissing, c)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, rec.Body.String())
}

// TestHTTPErrorHandler_Committed は書き込み済みのレスポンスを上書きしないことを検証する（準正常系）。
// アクセスログ（HandleError: true）で処理済みのエラーが、Echo 本体から再度渡されるケースに相当する。
func TestHTTPErrorHandler_Committed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, c.String(http.StatusOK, "ok"))

	HTTPErrorHandler(apperrors.ErrAuthTokenMissing, c)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

// TestSetupMiddlewares_ErrorResponse は SetupMiddlewares を通した本番同等の構成で、
// エラーが統一形式で 1 回だけ書き込まれることを検証する（アクセスログ経由の二重処理がないこと）。
func TestSetupMiddlewares_ErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			// 準正常系: ハンドラが返した sentinel が統一形式になる
			name:       "handler_sentinel",
			path:       "/unauthorized",
			wantStatus: http.StatusUnauthorized,
			wantBody:   `{"code":"AUTH_TOKEN_MISSING","error":"Authentication token is missing"}`,
		},
		{
			// 準正常系: 存在しないルートも統一形式になる
			name:       "unknown_route",
			path:       "/no-such-route",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"code":"NOT_FOUND","error":"Not Found"}`,
		},
		{
			// 異常系: panic は Recover で 500 の統一形式になる
			name:       "panic_recovered",
			path:       "/panic",
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"code":"INTERNAL_ERROR","error":"Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			SetupMiddlewares(e)
			e.GET("/unauthorized", func(c echo.Context) error {
				return apperrors.ErrAuthTokenMissing
			})
			e.GET("/panic", func(c echo.Context) error {
				panic("boom")
			})

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}
