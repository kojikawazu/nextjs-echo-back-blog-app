package supabase

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"backend/logger"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	// Ctx は Supabase とのやり取りに使用するグローバルなコンテキスト。
	Ctx = context.Background()
	// Pool は Supabase との接続プール。クエリ実行時に使用する。
	Pool *pgxpool.Pool
)

// buildConnString は接続URLに sslmode クエリを付与した接続文字列を組み立てる。
// sslmode が空の場合は "require"（本番デフォルト）を用いる。
// baseURL が既にクエリ文字列（"?" を含む）を持つ場合は "&" で連結する。
//
// 引数:
//   - baseURL: SUPABASE_URL 環境変数の値
//   - sslmode: DB_SSLMODE 環境変数の値（空なら "require"）
//
// 戻り値:
//   - string: sslmode を付与した接続文字列
func buildConnString(baseURL, sslmode string) string {
	if sslmode == "" {
		sslmode = "require"
	}
	sep := "?"
	if strings.Contains(baseURL, "?") {
		sep = "&"
	}
	return baseURL + sep + "sslmode=" + sslmode
}

// InitSupabase は Supabase の接続を初期化する。
// 接続URLを環境変数から取得し、コネクションプールを設定する。
// コネクションの最大数やアイドルタイム、シンプルプロトコルの使用を設定する。
// 成功時にはnilを返し、接続に失敗した場合はエラーメッセージを返す。
//
// SSL モードは環境変数 DB_SSLMODE で制御する（未設定時は "require"）。
// URL に既にクエリ文字列が含まれる場合は "&" で連結する。
//
// 戻り値:
//   - error: 初期化に成功した場合は nil、接続に失敗した場合のエラー
func InitSupabase() error {
	logger.InfoLog.Println("Initializing Supabase client...")
	supabaseURL := buildConnString(os.Getenv("SUPABASE_URL"), os.Getenv("DB_SSLMODE"))

	config, err := pgxpool.ParseConfig(supabaseURL)
	if err != nil {
		log.Printf("Unable to parse database URL: %v", err)
		return fmt.Errorf("unable to parse database URL: %v", err)
	}

	// コネクションプールの設定
	config.MaxConns = 10
	config.MaxConnIdleTime = 30 * time.Second
	// Prepared Statementの競合を防ぐためにSimple Protocolを優先
	config.ConnConfig.PreferSimpleProtocol = true

	logger.InfoLog.Println("Connecting supabase database...")
	Pool, err = pgxpool.ConnectConfig(Ctx, config)
	if err != nil {
		logger.ErrorLog.Printf("Unable to connect to Supabase: %v", err)
		return fmt.Errorf("unable to connect to Supabase: %v", err)
	}

	// 接続の確認
	logger.InfoLog.Println("Pinging supabase database...")
	err = Pool.Ping(Ctx)
	if err != nil {
		logger.ErrorLog.Printf("Unable to ping Supabase: %v", err)
		return fmt.Errorf("unable to ping Supabase: %v", err)
	}

	log.Println("Connected to Supabase successfully")
	return nil
}

// ClosePool は Supabase のコネクションプールをクローズする。
// この関数はアプリケーションのシャットダウン時に呼び出されることを想定する。
func ClosePool() {
	if Pool != nil {
		Pool.Close()
		log.Println("Supabase connection pool closed")
	}
}

// TestQuery は Supabase に対してシンプルなクエリを実行し、接続が正しく動作しているかを確認する。
// クエリ結果として "1" を取得し、それをログに出力する。
// クエリに失敗した場合、エラーを返す。
//
// 戻り値:
//   - error: クエリ実行に成功した場合は nil、失敗した場合のエラー
func TestQuery() error {
	logger.InfoLog.Println("Testing query...")
	query := `SELECT 1`
	rows, err := Pool.Query(Ctx, query)
	if err != nil {
		logger.ErrorLog.Printf("Failed to test query: %v", err)
		return err
	}
	logger.InfoLog.Println("Test query successful")
	defer rows.Close()

	for rows.Next() {
		var num int
		err := rows.Scan(&num)
		if err != nil {
			logger.ErrorLog.Printf("Failed to scan test query result: %v", err)
			return err
		}
		logger.InfoLog.Println("Test Query Result:", num)
	}

	logger.InfoLog.Println("Test query completed")
	return rows.Err()
}
