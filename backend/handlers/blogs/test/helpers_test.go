package handlers_blogs_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errServiceFailure はサービス層の想定外エラー（DB 障害など）を表すテスト用エラー
var errServiceFailure = errors.New("service failure")

// assertHandlerError はハンドラがエラーを返し、レスポンスを書き込んでいないことを検証する。
// ステータス・本文への変換は middlewares.HTTPErrorHandler が担うため、ハンドラ単体ではレスポンスが空であること、
// 戻り値が want を errors.Is で辿れることを確認する。
//
// 引数:
//   - t: テスト構造体
//   - err: ハンドラの戻り値
//   - rec: レスポンスレコーダー
//   - want: 戻り値に含まれているべきエラー
func assertHandlerError(t *testing.T, err error, rec *httptest.ResponseRecorder, want error) {
	t.Helper()
	assert.ErrorIs(t, err, want)
	assert.Empty(t, rec.Body.String(), "エラー時にハンドラはレスポンスを書き込まない")
}
