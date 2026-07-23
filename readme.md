# nextjs-echo-back-blog-app

[![Build and Deploy](https://github.com/kojikawazu/nextjs-echo-back-blog-app/actions/workflows/deploy.yml/badge.svg)](https://github.com/kojikawazu/nextjs-echo-back-blog-app/actions/workflows/deploy.yml)

Zenn・Qiita・GitHub に散在する技術記事を集約するブログプラットフォームの **バックエンド API**。
Go + Echo 製の REST API で、記事本文は GitHub 上の Markdown を参照し、メタデータ・いいね・コメントを Supabase(PostgreSQL) で管理する。フロントエンドは [別リポジトリ](https://github.com/kojikawazu/nextjs-echo-blog-front-app-renew)（Next.js）。

- 認証は JWT を HTTP-Only Cookie で伝達（有効期限1時間）
- クリーンアーキテクチャ（Handler → Service → Repository）＋インターフェース DI
- Google Cloud Run へ GitHub Actions で自動デプロイ、インフラは Terraform 管理

## 主な機能

| 機能 | 代表的なエンドポイント | 認証 |
|---|---|---|
| 認証 | `POST /api/users/login` ・ `GET /api/users/auth-check` ・ `POST /api/users/logout` | JWT Cookie |
| ブログ CRUD | `GET /api/blogs` ・ `GET /api/blogs/detail/:id` ・ `POST /api/blogs/create` ・ `PUT /api/blogs/update/:id` ・ `DELETE /api/blogs/delete/:id` | 参照は不要 / 更新系は要 |
| カテゴリ・タグ・人気 | `GET /api/blogs/categories` ・ `GET /api/blogs/tags` ・ `GET /api/blogs/popular/:count` | 不要 |
| いいね（匿名） | `GET /api/blog-likes/generate-visit-id` ・ `GET /api/blog-likes/is-liked/:blogId` ・ `POST /api/blog-likes/create/:blogId` ・ `DELETE /api/blog-likes/delete/:blogId` | visit-id Cookie |
| コメント（ゲスト） | `GET /api/comments/blog/:blogId` ・ `POST /api/comments/create` | 不要 |
| ヘルスチェック | `GET /` | 不要 |

> 全エンドポイントの仕様（リクエスト/レスポンス/ステータスコード）は [`docs/07-api-specification.md`](docs/07-api-specification.md) を参照。

## 技術スタック

| 分類 | 技術 | バージョン |
|---|---|---|
| 言語 | Go | 1.22 |
| Web フレームワーク | Echo | v4.12.0 |
| DB ドライバ | pgx/v4（pgxpool） | v4.18.3 |
| 認証 | golang-jwt/jwt（HS256） | v3.2.2 |
| データベース | Supabase（PostgreSQL） | - |
| 実行環境 | Google Cloud Run | - |
| IaC | Terraform | >= 1.6 |
| CI/CD | GitHub Actions | - |

## リポジトリ構成

モノレポ構成を採用しており、Go アプリケーション一式は [`backend/`](backend/) 配下に集約している。CI（`.github/`）・インフラ（`terraform/`）・ドキュメント（`docs/`）はリポジトリルートで全体共通として管理する。

```
.
├── backend/     # Go + Echo アプリ（main.go / go.mod / Dockerfile 含む）
├── terraform/   # インフラ構成（Terraform）
├── docs/        # 仕様・設計ドキュメント
├── .github/     # CI/CD（GitHub Actions）
└── README.md
```

> Go のモジュールパスは `module backend`。以降の `go` コマンドは `backend/` ディレクトリ内で実行する。

## セットアップ（ローカル起動）

前提: Go 1.20 以上、稼働中の Supabase(PostgreSQL) への接続情報。

```bash
git clone https://github.com/kojikawazu/nextjs-echo-back-blog-app.git
cd nextjs-echo-back-blog-app/backend

# 環境変数を用意（値は各自で設定）
cp .env.example .env
#   - JWT_SECRET_KEY: 必須。未設定だと起動時に log.Fatal で即終了する
#   - SUPABASE_URL  : 必須。Supabase の接続URL（?sslmode=require はコードが自動付与）
#   - ALLOWED_ORIGINS: フロントの URL（例: http://localhost:3000）

# 依存解決して起動
go mod download
go run main.go
```

起動後、`http://localhost:8080/` で `Service is running` が返れば成功（ポートは `PORT` で変更可、既定 8080）。

### Docker で起動する場合

ビルドコンテキストは `backend/` を指定する（リポジトリルートから実行する場合）。

```bash
docker build -t echo-blog-back ./backend
docker run --env-file backend/.env -p 8080:8080 echo-blog-back
```

## テスト

`go` コマンドは `backend/` 内で実行する。テストは 2 種類に分かれる。

```bash
cd backend

# 単体テスト（UT）: DB 不要・全 mock。handlers / services。
go test ./handlers/... ./services/...

# インテグレーションテスト（IT）: repositories 層。
# 既定では testcontainers が使い捨ての PostgreSQL を起動する（Docker 稼働が前提）。
go test -tags=integration ./repositories/...
```

- IT は `//go:build integration` タグで分離しているため、通常の `go test ./...` には含まれない。
- IT は既定でコンテナを起動するため `.env.test` は不要。実 DB に対して実行したい場合のみ `SUPABASE_URL`（と `TEST_*`）を環境変数で渡す（[`backend/.env.test.example`](backend/.env.test.example) 参照）。

> CI（`ci.yml`）は PR・`main` への push で、`test` ジョブ（UT）と `integration-test` ジョブ（IT / testcontainers）の両方を実行する。テスト方針の詳細は [`docs/08-test-specification.md`](docs/08-test-specification.md)。

## Lint / 静的解析

`go` コマンド・`golangci-lint` は `backend/` 内で実行する。CI（`ci.yml`）は PR・`main` への push で同じチェックを実行する。

```bash
cd backend

# フォーマット確認（差分があるファイルを列挙。整形は gofmt -w .）
gofmt -l .

# 標準の静的解析
go vet ./...

# 統合リンタ（設定は backend/.golangci.yml。v2 系を使用）
#   未導入の場合: https://golangci-lint.run/welcome/install/
golangci-lint run
```

## ドキュメント

仕様書・設計書は [`docs/`](docs/) に集約。索引は [`docs/README.md`](docs/README.md)。

| ドキュメント | 内容 |
|---|---|
| [01 要求](docs/01-business-requirements.md) / [02 要件](docs/02-requirements-specification.md) | 背景・目的・スコープ・機能要件 |
| [03 機能](docs/03-functional-specification.md) | 処理フロー・バリデーション |
| [05 データ](docs/05-data-specification.md) | テーブル・モデル定義 |
| [06 セキュリティ](docs/06-security-specification.md) | 認証・Cookie・CORS |
| [07 API](docs/07-api-specification.md) | 全エンドポイント仕様 |
| [09 アーキテクチャ](docs/09-architecture-specification.md) | 構成・DI・インフラ・**ローカル開発セットアップ** |
| [11 タスク](docs/11-tasks.md) | 実績・改善タスク（技術的負債） |

開発ルールは [`CLAUDE.md`](CLAUDE.md) と [`.claude/rules/`](.claude/rules/) を参照。

## 関連リポジトリ

- フロントエンド: [nextjs-echo-blog-front-app-renew](https://github.com/kojikawazu/nextjs-echo-blog-front-app-renew)（Next.js + Hono）

## ライセンス

[LICENSE](LICENSE) を参照。
