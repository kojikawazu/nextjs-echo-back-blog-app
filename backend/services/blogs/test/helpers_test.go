package services_blogs_test

import (
	"backend/apperrors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// テストで使う UUID 形式の ID（サービス層の UUID 検証を通過する値）
const (
	testBlogID = "11111111-1111-1111-1111-111111111111"
	testUserID = "22222222-2222-2222-2222-222222222222"
)

// errRepositoryFailure はリポジトリ層の想定外エラー（DB 障害など）を表すテスト用エラー
var errRepositoryFailure = errors.New("repository failure")

// assertInternalError は err がクライアント向けの sentinel（*apperrors.Error）を含まず（= 500 として扱われる）、
// 原因の cause を errors.Is で辿れることを検証する。
//
// 引数:
//   - t: テスト構造体
//   - err: 検証対象のエラー
//   - cause: 原因として含まれているべきエラー
func assertInternalError(t *testing.T, err error, cause error) {
	t.Helper()
	assert.ErrorIs(t, err, cause)
	var appErr *apperrors.Error
	assert.False(t, errors.As(err, &appErr), "想定外エラーは sentinel に変換しない: %v", err)
}
