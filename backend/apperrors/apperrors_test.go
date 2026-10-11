package apperrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew は生成したエラーが分類・コード・メッセージを保持することを検証する（正常系）。
func TestNew(t *testing.T) {
	err := New(KindNotFound, CodeBlogNotFound, "Blog not found")

	assert.Equal(t, KindNotFound, err.Kind())
	assert.Equal(t, "BLOG_NOT_FOUND", err.Code())
	assert.Equal(t, "Blog not found", err.Message())
	assert.Equal(t, "Blog not found", err.Error())
}

// TestErrorsIs は sentinel の同一性で判定され、メッセージが同じ別のエラーとは一致しないことを検証する。
func TestErrorsIs(t *testing.T) {
	sentinel := New(KindInvalid, CodeValidation, "Invalid id")

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			// 正常系: sentinel そのもの
			name: "same_sentinel",
			err:  sentinel,
			want: true,
		},
		{
			// 準正常系: fmt.Errorf で包まれていても一致する
			name: "wrapped_by_fmt",
			err:  fmt.Errorf("fetch blog: %w", sentinel),
			want: true,
		},
		{
			// 準正常系: メッセージが同じでも別の sentinel とは一致しない（文字列比較に依存しない）
			name: "same_message_other_sentinel",
			err:  New(KindInvalid, CodeValidation, "Invalid id"),
			want: false,
		},
		{
			// 準正常系: メッセージが同じ標準エラーとも一致しない
			name: "same_message_plain_error",
			err:  errors.New("Invalid id"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, errors.Is(tt.err, sentinel))
		})
	}
}

// TestWrap は原因を添えたエラーが sentinel・原因の両方に一致し、Error() に原因を含むことを検証する。
func TestWrap(t *testing.T) {
	cause := errors.New("unexpected EOF")

	err := Wrap(ErrInvalidRequestBody, cause)

	// 正常系: sentinel と原因の両方を errors.Is で辿れる
	assert.ErrorIs(t, err, ErrInvalidRequestBody)
	assert.ErrorIs(t, err, cause)
	// 準正常系: 別の sentinel とは一致しない
	assert.NotErrorIs(t, err, ErrAuthTokenMissing)
	// 準正常系: errors.As で取り出せるのは sentinel（レスポンスに使う情報）
	var appErr *Error
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, CodeInvalidRequestBody, appErr.Code())
	// ログ用の Error() には原因を含む
	assert.Equal(t, "Invalid request body: unexpected EOF", err.Error())
}
