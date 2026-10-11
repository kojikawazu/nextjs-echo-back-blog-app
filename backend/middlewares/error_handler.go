package middlewares

import (
	"backend/apperrors"
	"backend/dto"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// HTTPErrorHandler はハンドラ・ミドルウェアが返したエラーを、HTTP ステータスと統一エラーレスポンス（dto.ErrorResponse）に変換して返す。
// Echo の e.HTTPErrorHandler に設定する。エラーの詳細はアクセスログ（newRequestLogger）に記録され、レスポンスには含めない。
//
// 引数:
//   - err: ハンドラ・ミドルウェアが返したエラー
//   - c: リクエスト/レスポンスを保持する echo コンテキスト
func HTTPErrorHandler(err error, c echo.Context) {
	// アクセスログ（HandleError: true）経由で処理済みのエラーが再度渡されるため、書き込み済みなら何もしない
	if c.Response().Committed {
		return
	}

	status, body := toErrorResponse(err)

	var writeErr error
	if c.Request().Method == http.MethodHead {
		writeErr = c.NoContent(status)
	} else {
		writeErr = c.JSON(status, body)
	}
	if writeErr != nil {
		c.Logger().Error(writeErr)
	}
}

// toErrorResponse はエラーを HTTP ステータスとレスポンス本文に変換する。
//
// 引数:
//   - err: 変換するエラー
//
// 戻り値:
//   - int: HTTP ステータス
//   - dto.ErrorResponse: レスポンス本文
func toErrorResponse(err error) (int, dto.ErrorResponse) {
	var appErr *apperrors.Error
	if errors.As(err, &appErr) {
		return statusOf(appErr.Kind()), dto.ErrorResponse{Code: appErr.Code(), Message: appErr.Message()}
	}

	// Echo 自身が返すエラー（ルート不在 404・405・ミドルウェアの拒否など）
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) && httpErr.Code < http.StatusInternalServerError {
		return httpErr.Code, dto.ErrorResponse{
			Code:    codeOfStatus(httpErr.Code),
			Message: http.StatusText(httpErr.Code),
		}
	}

	// 想定外のエラーは内部情報を漏らさないよう汎用メッセージで返す
	return http.StatusInternalServerError, dto.ErrorResponse{
		Code:    apperrors.CodeInternal,
		Message: "Internal server error",
	}
}

// statusOf はエラーの分類を HTTP ステータスに対応付ける。
//
// 引数:
//   - kind: エラーの分類
//
// 戻り値:
//   - int: HTTP ステータス
func statusOf(kind apperrors.Kind) int {
	switch kind {
	case apperrors.KindInvalid:
		return http.StatusBadRequest
	case apperrors.KindUnauthorized:
		return http.StatusUnauthorized
	case apperrors.KindNotFound:
		return http.StatusNotFound
	case apperrors.KindConflict:
		return http.StatusConflict
	case apperrors.KindTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// codeOfStatus は HTTP ステータスからエラーコードを作る（例: 404 → "NOT_FOUND"、405 → "METHOD_NOT_ALLOWED"）。
//
// 引数:
//   - status: HTTP ステータス
//
// 戻り値:
//   - string: エラーコード
func codeOfStatus(status int) string {
	return strings.ToUpper(strings.ReplaceAll(http.StatusText(status), " ", "_"))
}
