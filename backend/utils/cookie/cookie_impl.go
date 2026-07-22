package utils_cookie

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// CookieUtils インターフェース
type CookieUtils interface {
	// Token用

	// GetAuthCookie は認証用のCookieを取得する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenName: 取得対象のCookie名
	//
	// 戻り値:
	//   - *http.Cookie: 取得したCookie
	//   - error: Cookieが存在しない等、取得に失敗した場合のエラー
	GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error)

	// GetAuthCookieValue は認証用のCookieの値を取得する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenName: 取得対象のCookie名
	//
	// 戻り値:
	//   - string: 取得したCookieの値
	//   - error: Cookieが存在しない等、取得に失敗した場合のエラー
	GetAuthCookieValue(c echo.Context, tokenName string) (string, error)

	// GetAuthCookieExpirationTime は認証用のCookieの有効期限を取得する。
	//
	// 戻り値:
	//   - time.Time: 現在時刻から1時間後の有効期限
	GetAuthCookieExpirationTime() time.Time

	// ExistsAuthCookie は認証用のCookieが存在するか確認する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenName: 確認対象のCookie名
	//
	// 戻り値:
	//   - bool: Cookieが存在する場合はtrue、存在しない場合はfalse
	ExistsAuthCookie(c echo.Context, tokenName string) bool

	// VerifyToken はJWTトークンを検証する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: 検証対象のJWTトークン文字列
	//
	// 戻り値:
	//   - *models.Claims: 検証に成功したトークンのクレーム情報
	//   - error: トークンが無効または有効期限切れの場合のエラー
	VerifyToken(c echo.Context, tokenString string) (*models.Claims, error)

	// 認証Token用

	// CreateToken はJWTトークンを作成する。
	//
	// 引数:
	//   - user: トークンに埋め込むユーザー情報
	//
	// 戻り値:
	//   - string: 生成したJWTトークン文字列
	//   - error: トークンの署名に失敗した場合のエラー
	CreateToken(user *models.BlogUsersData) (string, error)

	// AddAuthCookie は認証用のCookieを追加する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: Cookieに保存するJWTトークン文字列
	//   - expirationTime: Cookieの有効期限
	AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)

	// UpdateAuthCookie は認証用のCookieを更新する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: Cookieに保存する更新後のJWTトークン文字列
	//   - expirationTime: Cookieの有効期限
	UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)

	// DelAuthCookie は認証用のCookieを削除する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	DelAuthCookie(c echo.Context)

	// GetUserIdFromToken はJWTトークンを解析してユーザーIDを取得する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: 解析対象のJWTトークン文字列
	//
	// 戻り値:
	//   - string: トークンから取得したユーザーID
	//   - error: トークンが無効または有効期限切れの場合のエラー
	GetUserIdFromToken(c echo.Context, tokenString string) (string, error)

	// VisitId用

	// CreateVisitIdToken はvisitId用JWTトークンを作成する。
	//
	// 戻り値:
	//   - string: 生成したvisitId用JWTトークン文字列
	//   - error: トークンの署名に失敗した場合のエラー
	CreateVisitIdToken() (string, error)

	// AddVisitIdCoookie は訪問者ID用のCookieを追加する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: Cookieに保存するvisitId用JWTトークン文字列
	//   - expirationTime: Cookieの有効期限
	AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time)

	// GetVisitIdFromToken は訪問者IDを取得する。
	//
	// 引数:
	//   - c: Echoのリクエストコンテキスト
	//   - tokenString: 解析対象のvisitId用JWTトークン文字列
	//
	// 戻り値:
	//   - string: トークンから取得した訪問者ID
	//   - error: トークンが無効または有効期限切れの場合のエラー
	GetVisitIdFromToken(c echo.Context, tokenString string) (string, error)
}

// CookieUtilsImpl は CookieUtils の実装。
type CookieUtilsImpl struct{}

// NewCookieUtils は CookieUtils を生成する。
//
// 戻り値:
//   - CookieUtils: 生成したCookieUtilsの実装
func NewCookieUtils() CookieUtils {
	return &CookieUtilsImpl{}
}
