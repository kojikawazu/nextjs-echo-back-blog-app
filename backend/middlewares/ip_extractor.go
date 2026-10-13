package middlewares

import (
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// clientIPHeader は BFF（フロントエンド）がクライアント IP を転送するヘッダ
	clientIPHeader = "X-Client-IP"

	// trustClientIPHeaderEnv は clientIPHeader を信頼するかを切り替える環境変数（"true" のときのみ信頼する）
	trustClientIPHeaderEnv = "TRUST_CLIENT_IP_HEADER"
)

// NewIPExtractor はクライアント IP（c.RealIP()）の抽出方法を返す。
//
// 既定では X-Forwarded-For の右端（Cloud Run のフロントエンドが付与した値）から取る。
// 既定の c.RealIP() は左端を使うため、ヘッダの偽装でレートリミットを回避できてしまう。
//
// trustClientIPHeader が true のときは、BFF が付与した X-Client-IP を優先する。
// X-Client-IP は送信者が自由に書けるため、バックエンドを BFF からのみ呼べる状態
// （Cloud Run IAM で未認証の呼び出しを拒否）でのみ true にする。
// 値が IP として解釈できない場合は既定の抽出方法に戻す（BFF の IP として数える）。
//
// 引数:
//   - trustClientIPHeader: X-Client-IP を信頼するか
//
// 戻り値:
//   - echo.IPExtractor: クライアント IP の抽出関数
func NewIPExtractor(trustClientIPHeader bool) echo.IPExtractor {
	fromXFF := echo.ExtractIPFromXFFHeader()
	if !trustClientIPHeader {
		return fromXFF
	}

	return func(req *http.Request) string {
		if ip := net.ParseIP(strings.TrimSpace(req.Header.Get(clientIPHeader))); ip != nil {
			return ip.String()
		}
		return fromXFF(req)
	}
}
