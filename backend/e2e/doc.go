// Package e2e はバックエンド API の E2E（エンドツーエンド）テストを集約する。
//
// 実サーバー（全ルート）を httptest で起動し、HTTP リクエストを通じて
// Handler → Service → Repository → DB を貫くフローを検証する。
// DB は testsupport（testcontainers の使い捨て PostgreSQL）を再利用する。
//
// テスト本体は `//go:build e2e` で分離しており、実行は次のとおり:
//
//	JWT_SECRET_KEY=... go test -tags=e2e ./e2e/...
//
// このファイル（build tag なし）は、デフォルトビルドでも e2e パッケージを
// 有効に保つためのプレースホルダであり、テストコードは含まない。
package e2e
