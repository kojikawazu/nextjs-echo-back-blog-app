// Package apperrors はアプリケーション全体で使うエラーの分類（Kind）と、クライアント向けのコード・メッセージを定義する。
//
// サービス層・ハンドラ層は *Error の sentinel を返し、呼び出し側は errors.Is で判定する（エラーメッセージの文字列比較をしない）。
// HTTP ステータスとレスポンス本文への変換は middlewares.HTTPErrorHandler が 1 か所で行う。
// Kind は HTTP に依存しない分類のため、サービス層は HTTP を意識せずにエラーの種類を表現できる。
package apperrors

import "fmt"

// Kind はエラーの分類。HTTP ステータスへの対応は middlewares.HTTPErrorHandler が決める。
type Kind int

const (
	// KindInternal はサーバー内部のエラー（500）。
	KindInternal Kind = iota
	// KindInvalid は入力値の不正（400）。
	KindInvalid
	// KindUnauthorized は未認証・認証失敗（401）。
	KindUnauthorized
	// KindNotFound はリソースの不在（404）。
	KindNotFound
	// KindConflict は現在の状態と競合する操作（409）。
	KindConflict
	// KindTooManyRequests はレートリミット超過（429）。
	KindTooManyRequests
)

// クライアントが機械的に判定できるエラーコード。値は API 契約のため、一度公開したら変更しない。
const (
	CodeValidation             = "VALIDATION_ERROR"
	CodeInvalidRequestBody     = "INVALID_REQUEST_BODY"
	CodeAuthTokenMissing       = "AUTH_TOKEN_MISSING"
	CodeAuthTokenInvalid       = "AUTH_TOKEN_INVALID"
	CodeInvalidCredentials     = "INVALID_CREDENTIALS"
	CodeInvalidCurrentPassword = "INVALID_CURRENT_PASSWORD"
	CodeBlogNotFound           = "BLOG_NOT_FOUND"
	CodeUserNotFound           = "USER_NOT_FOUND"
	CodeAlreadyLiked           = "ALREADY_LIKED"
	CodeTooManyRequests        = "TOO_MANY_REQUESTS"
	CodeInternal               = "INTERNAL_ERROR"
)

// Error はクライアントに返してよいコード・メッセージを持つアプリケーションエラー。
// sentinel として定義し、ポインタの同一性で errors.Is 判定する。
type Error struct {
	kind    Kind
	code    string
	message string
}

// New はアプリケーションエラーを生成する。パッケージ変数（sentinel）の定義に用いる。
//
// 引数:
//   - kind: エラーの分類
//   - code: クライアント向けのエラーコード
//   - message: クライアント向けのメッセージ（内部情報を含めない）
//
// 戻り値:
//   - *Error: 生成したエラー
func New(kind Kind, code, message string) *Error {
	return &Error{kind: kind, code: code, message: message}
}

// Error はクライアント向けのメッセージを返す。
func (e *Error) Error() string { return e.message }

// Kind はエラーの分類を返す。
func (e *Error) Kind() Kind { return e.kind }

// Code はクライアント向けのエラーコードを返す。
func (e *Error) Code() string { return e.code }

// Message はクライアント向けのメッセージを返す。
func (e *Error) Message() string { return e.message }

// Wrap は sentinel に原因のエラーを添える。
// 戻り値は errors.Is で sentinel・原因の両方に一致し、Error() は原因を含む（ログ用。レスポンスには sentinel のメッセージのみを使う）。
//
// 引数:
//   - sentinel: 返したいアプリケーションエラー
//   - cause: 原因のエラー
//
// 戻り値:
//   - error: 原因を添えたエラー
func Wrap(sentinel *Error, cause error) error {
	return fmt.Errorf("%w: %w", sentinel, cause)
}

// ハンドラ層・ミドルウェアで共通に使う sentinel。
var (
	// ErrInvalidRequestBody はリクエストボディのバインドに失敗したことを示す。
	ErrInvalidRequestBody = New(KindInvalid, CodeInvalidRequestBody, "Invalid request body")
	// ErrAuthTokenMissing は認証トークン（Cookie）が無いことを示す。
	ErrAuthTokenMissing = New(KindUnauthorized, CodeAuthTokenMissing, "Authentication token is missing")
	// ErrAuthTokenInvalid は認証トークンが不正・期限切れであることを示す。
	ErrAuthTokenInvalid = New(KindUnauthorized, CodeAuthTokenInvalid, "Authentication token is invalid")
	// ErrTooManyRequests はレートリミットを超過したことを示す。
	ErrTooManyRequests = New(KindTooManyRequests, CodeTooManyRequests, "Too many requests")
)
