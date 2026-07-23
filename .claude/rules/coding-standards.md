---
description: コーディング規約
globs: 
---

# コーディング規約

- **言語**: Go 1.22（`go.mod` の要求バージョンに準拠。Dockerfile のビルドステージも同一メジャーに揃える）
- **パッケージマネージャ**: go mod を使用（`go mod download` / `go mod tidy` で依存管理。`go` コマンドは `backend/` 内で実行）
- **Linter / Formatter**: golangci-lint + gofmt / goimports でコード品質を担保
- **環境変数**: 設定値は環境変数で管理（.env）
- **シークレット禁止**: シークレット・認証情報をハードコードしない
