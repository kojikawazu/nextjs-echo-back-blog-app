package handlers_blog_comments

import (
	"backend/apperrors"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errService はテストでサービスのモックが返す想定外のエラー。ハンドラの戻り値から errors.Is で辿れることを検証する。
var errService = errors.New("db error")

// assertHandlerError はハンドラがレスポンスを書かずに wantErr を返したことを検証する。
// ステータス・本文への変換は middlewares.HTTPErrorHandler の責務のため、ここではエラーの同一性と分類のみを見る。
// wantErr が apperrors の sentinel でない場合（= HTTPErrorHandler で 500 になる）、戻り値も apperrors.Error に該当しないことを検証する。
//
// 引数:
//   - t: テスト構造体
//   - err: ハンドラの戻り値
//   - rec: レスポンスレコーダ
//   - wantErr: 期待するエラー
func assertHandlerError(t *testing.T, err error, rec *httptest.ResponseRecorder, wantErr error) {
	t.Helper()
	assert.ErrorIs(t, err, wantErr)
	var wantAppErr, gotAppErr *apperrors.Error
	assert.Equal(t, errors.As(wantErr, &wantAppErr), errors.As(err, &gotAppErr))
	assert.Empty(t, rec.Body.String())
}
