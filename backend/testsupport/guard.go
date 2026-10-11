//go:build integration || e2e

package testsupport

import (
	"errors"
	"net"
	"strings"

	"github.com/jackc/pgconn"
)

// errRemoteDB は SUPABASE_URL がローカル以外の DB を指している場合のエラー。
// 接続文字列・ホスト名は CI ログ（公開）に残さないため、メッセージに含めない。
var errRemoteDB = errors.New("SUPABASE_URL がローカル以外の DB を指しているため IT / E2E を中止しました。" +
	"テストはデータを変更するため、localhost / 127.0.0.1 / ::1 / Unix ソケットの DB のみ許可しています" +
	"（SUPABASE_URL を外せば testcontainers で実行されます）")

// errUnparsableDB は SUPABASE_URL を接続設定として解析できない場合のエラー。
var errUnparsableDB = errors.New("SUPABASE_URL を接続設定として解析できないため IT / E2E を中止しました")

// checkLocalDB は接続文字列の接続先がすべてローカルであることを検証する。
// IT / E2E は書き込みを伴うため、リモート（本番等）の DB への接続を拒否する。
//
// 引数:
//   - connString: 検証する接続文字列（URL 形式 / key=value 形式）
//
// 戻り値:
//   - error: 接続先にローカル以外が含まれる、または解析できない場合のエラー
func checkLocalDB(connString string) error {
	config, err := pgconn.ParseConfig(connString)
	if err != nil {
		return errUnparsableDB
	}

	hosts := []string{config.Host}
	for _, fb := range config.Fallbacks {
		hosts = append(hosts, fb.Host)
	}
	for _, h := range hosts {
		if !isLocalHost(h) {
			return errRemoteDB
		}
	}
	return nil
}

// isLocalHost はホストがローカル（ループバック / Unix ソケット）かを判定する。
//
// 引数:
//   - host: 判定するホスト（ホスト名 / IP / Unix ソケットのディレクトリ）
//
// 戻り値:
//   - bool: ローカルなら true
func isLocalHost(host string) bool {
	if strings.HasPrefix(host, "/") {
		return true
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
