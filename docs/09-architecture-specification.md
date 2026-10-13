# アーキテクチャ仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo）のシステムアーキテクチャ、レイヤー構成、インフラストラクチャ、CI/CDパイプラインを定義するアーキテクチャ仕様書です。コードベースのリバースエンジニアリングに基づいて作成しています。

## 目次

- [1. システムアーキテクチャ概要](#1-システムアーキテクチャ概要)
- [2. レイヤーアーキテクチャ](#2-レイヤーアーキテクチャ)
  - [2.1 レイヤー構成図](#21-レイヤー構成図)
  - [2.2 データフロー](#22-データフロー)
  - [2.3 エラーハンドリング](#23-エラーハンドリング)
- [3. DIパターン（依存性注入）](#3-diパターン依存性注入)
  - [3.1 インターフェース設計](#31-インターフェース設計)
  - [3.2 ドメイン別インターフェース](#32-ドメイン別インターフェース)
  - [3.3 依存関係の組み立て](#33-依存関係の組み立て)
- [4. 技術スタック](#4-技術スタック)
  - [4.1 アプリケーション](#41-アプリケーション)
  - [4.2 インフラストラクチャ](#42-インフラストラクチャ)
  - [4.3 ミドルウェア構成](#43-ミドルウェア構成)
- [5. インフラストラクチャ](#5-インフラストラクチャ)
  - [5.1 GCPリソース構成](#51-gcpリソース構成)
  - [5.2 Terraform構成](#52-terraform構成)
- [6. CI/CDパイプライン](#6-cicdパイプライン)
  - [6.1 GitHub Actions ワークフロー](#61-github-actions-ワークフロー)
- [7. デプロイメント](#7-デプロイメント)
  - [7.1 Dockerイメージ構成](#71-dockerイメージ構成)
  - [7.2 環境変数](#72-環境変数)
  - [7.3 サーバー起動フロー](#73-サーバー起動フロー)
- [8. ディレクトリ構成](#8-ディレクトリ構成)
- [9. セキュリティ設計](#9-セキュリティ設計)
  - [9.1 認証・認可](#91-認証認可)
  - [9.2 環境識別](#92-環境識別)
  - [9.3 シークレット管理](#93-シークレット管理)
  - [9.4 コンテナセキュリティ](#94-コンテナセキュリティ)
- [10. ローカル開発セットアップ](#10-ローカル開発セットアップ)
  - [10.1 前提ツール](#101-前提ツール)
  - [10.2 セットアップ手順](#102-セットアップ手順)
  - [10.3 環境変数](#103-環境変数)
  - [10.4 テスト実行](#104-テスト実行)
  - [10.5 Docker での起動](#105-docker-での起動)
  - [10.6 よくあるつまずき](#106-よくあるつまずき)

---

## 1. システムアーキテクチャ概要

```
┌─────────────┐       ┌──────────────────────────────────────────────┐
│             │       │            Google Cloud Platform             │
│  Frontend   │       │                                              │
│  (Next.js)  │──────>│  ┌──────────────────────────────────────┐    │
│             │ HTTPS │  │          Cloud Run                    │    │
└─────────────┘       │  │  ┌──────────────────────────────────┐│    │
                      │  │  │  Go + Echo Application           ││    │
                      │  │  │                                  ││    │
                      │  │  │  ┌──────────┐   ┌────────────┐  ││    │
                      │  │  │  │ Handler  │──>│  Service    │  ││    │
                      │  │  │  └──────────┘   └─────┬──────┘  ││    │
                      │  │  │                       │          ││    │
                      │  │  │                 ┌─────┴──────┐   ││    │
                      │  │  │                 │ Repository │   ││    │
                      │  │  │                 └─────┬──────┘   ││    │
                      │  │  │                       │          ││    │
                      │  │  └───────────────────────┼──────────┘│    │
                      │  └──────────────────────────┼───────────┘    │
                      │                             │                │
                      │  ┌──────────────────┐       │                │
                      │  │  Secret Manager  │       │                │
                      │  │  (環境変数管理)   │       │                │
                      │  └──────────────────┘       │                │
                      │                             │                │
                      │  ┌──────────────────┐       │                │
                      │  │Artifact Registry │       │                │
                      │  │ (Docker Image)   │       │                │
                      │  └──────────────────┘       │                │
                      └─────────────────────────────┼────────────────┘
                                                    │
                                              ┌─────┴──────────┐
                                              │   Supabase     │
                                              │  (PostgreSQL)  │
                                              └────────────────┘
```

---

## 2. レイヤーアーキテクチャ

本プロジェクトはクリーンアーキテクチャの考え方（インターフェースによる依存性逆転）に基づくレイヤード構成（Handler → Service → Repository の 3 層）を採用している。各層はインターフェースを介して疎結合に接続されている。

ディレクトリ・層の責務の規定は [`.claude/rules/api.md`](../.claude/rules/api.md) を正とする。Command / Query の分離や `internal/` によるパッケージの隠蔽は、現在の規模（4 テーブル・CRUD 中心）では採用せず、規模が大きくなった時点で再検討する（#125）。

### 2.1 レイヤー構成図

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Request                       │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                 Middleware Layer                      │
│  ・Logger（リクエストログ）                            │
│  ・Recover（パニックリカバリ）                          │
│  ・CORS（クロスオリジン制御）                           │
│  ・RateLimiter（レートリミット）                        │
│  ・HTTPErrorHandler（エラー → ステータス・統一形式）     │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                 Handler Layer                        │
│  ・HTTPリクエストのパース（パスパラメータ、ボディ）      │
│  ・Cookie/JWT トークンの取得                          │
│  ・Service 呼び出し                                   │
│  ・成功レスポンスの生成（ステータスコード、DTO）         │
│  ・エラーはそのまま return（変換しない）                 │
└──────────────────────┬──────────────────────────────┘
                       │  Interface (依存性逆転)
                       ▼
┌─────────────────────────────────────────────────────┐
│                 Service Layer                        │
│  ・入力バリデーション（空文字、メール形式）              │
│  ・ビジネスロジック（タグ分割、重複チェック等）          │
│  ・Repository 呼び出し                                │
│  ・sentinel error（apperrors）の返却とログ出力          │
└──────────────────────┬──────────────────────────────┘
                       │  Interface (依存性逆転)
                       ▼
┌─────────────────────────────────────────────────────┐
│                Repository Layer                      │
│  ・SQLクエリの構築と実行                               │
│  ・Supabase（PostgreSQL）への接続プール管理             │
│  ・結果のスキャンとモデルへのマッピング                  │
│  ・JOINによるいいね数・コメント数の集計                  │
│  ・pgx.ErrNoRows → ErrXxxNotFound への変換              │
└──────────────────────┬──────────────────────────────┘
                       │  pgx/v5 pgxpool
                       ▼
┌─────────────────────────────────────────────────────┐
│                  PostgreSQL (Supabase)                │
│  ・blogs テーブル                                     │
│  ・blog_users テーブル                                │
│  ・blog_likes テーブル                                │
│  ・blog_comments テーブル                             │
└─────────────────────────────────────────────────────┘
```

### 2.2 データフロー

1. クライアントからHTTPリクエストを受信
2. ミドルウェアでロギング、リカバリ、CORS処理、レートリミットを実行
3. ルーターが対応するHandlerにリクエストを振り分け
4. Handlerがリクエストをパースし、必要に応じてCookieからJWTを検証
5. HandlerがServiceの対応メソッドを呼び出す
6. Serviceが入力バリデーションを実行後、Repositoryを呼び出す
7. RepositoryがSQLを実行しデータを取得/変更
8. 結果がService → Handler → クライアントへ返却される（エラー時は 2.3 のとおり HTTPErrorHandler が変換して返す）

### 2.3 エラーハンドリング

エラーレスポンスの構造・ステータスの決定は `middlewares.HTTPErrorHandler`（`e.HTTPErrorHandler`）の 1 か所に集約する（#123）。ハンドラは `c.JSON` でエラーを書かず、エラーを `return` するだけにする。

```
Repository ── pgx.ErrNoRows を ErrBlogNotFound / ErrUserNotFound に変換（ドライバ固有のエラーを漏らさない）
    │
Service ───── 入力不正・不在・競合などを sentinel（apperrors.Error）で返す。想定外は fmt.Errorf("...: %w") で包む
    │
Handler ───── return err（Cookie 不在/不正・bind 失敗は apperrors.Wrap(ErrAuthTokenMissing 等, 原因) を返す）
    │
HTTPErrorHandler ── errors.As で apperrors.Error を取り出し、Kind → ステータス、Code・Message → 本文
```

| 項目 | 内容 |
|------|------|
| エラー定義 | `apperrors` パッケージ。`Kind`（HTTP 非依存の分類）・`Code`（クライアント向けの安定した識別子）・`Message` を持つ。共通の sentinel（`ErrInvalidRequestBody` / `ErrAuthTokenMissing` / `ErrAuthTokenInvalid` / `ErrTooManyRequests`）を定義 |
| ドメインのエラー | 各 `services/<domain>/errors.go` に sentinel を定義（例: `services_blogs.ErrBlogNotFound`）。判定は `errors.Is`（メッセージの文字列比較をしない） |
| Kind → ステータス | `KindInvalid` 400 / `KindUnauthorized` 401 / `KindNotFound` 404 / `KindConflict` 409 / `KindTooManyRequests` 429 / `KindInternal`・その他 500 |
| レスポンス本文 | `dto.ErrorResponse`（`{"code": "...", "error": "..."}`）。コード一覧は [07-api-specification.md](07-api-specification.md) |
| Echo 由来のエラー | `echo.HTTPError`（ルート不在 404・405 等）の 4xx はステータス名からコードを生成（`NOT_FOUND` 等）。メッセージはステータス名とし、独自メッセージは返さない |
| 想定外のエラー | 500 `{"code":"INTERNAL_ERROR","error":"Internal server error"}`。原因は本文に含めず、アクセスログの `error` に記録する |
| ログ | アクセスログ（`HandleError: true`）がエラー全体（原因を含む）を記録する。4xx は WARN、5xx は ERROR。HTTPErrorHandler は書き込み済み（`Committed`）なら何もしない（アクセスログ経由で処理済みのエラーが再度渡されるため） |

---

## 3. DIパターン（依存性注入）

### 3.1 インターフェース設計

各ドメインは以下の3ファイル構成で依存性注入パターンを実装している。

| ファイル | 役割 | 内容 |
|---------|------|------|
| `xxx.go` | インターフェース定義 | メソッドシグネチャとコンストラクタ |
| `xxx_impl.go` | 実装 | インターフェースの具象実装 |
| `xxx_mock.go` | モック | テスト用のモック実装（testify/mock） |

### 3.2 ドメイン別インターフェース

#### BlogRepository

```go
type BlogRepository interface {
    FetchBlogs() ([]models.BlogData, error)
    FetchBlogsByUserId(userId string) ([]models.BlogData, error)
    FetchBlogById(id string) (*models.BlogData, error)
    CreateBlog(userId, title, githubUrl, category, description, tags string) (*models.BlogData, error)
    UpdateBlog(id, title, githubUrl, category, description, tags string) (*models.BlogData, error)
    DeleteBlog(id string) error
    FetchBlogCategories() ([]string, error)
    FetchBlogTags() ([]string, error)
    FetchBlogPopular(count int) ([]models.BlogData, error)
}
```

#### BlogUsersRepository

```go
type BlogUsersRepository interface {
    FetchBlogUsersByEmailAndPassword(email, password string) (*models.BlogUsersData, error)
    FetchBlogUsersById(id string) (*models.BlogUsersData, error)
    UpdateBlogUsers(id, name, email, password string) (*models.BlogUsersData, error)
}
```

#### BlogLikeRepository

```go
type BlogLikeRepository interface {
    FetchBlogLikesByVisitId(visitId string) ([]models.BlogLikesData, error)
    IsBlogLiked(blogId, visitId string) (bool, error)
    CreateBlogLike(blogId, visitId string) (*models.BlogLikesData, error)
    DeleteBlogLike(blogId, visitId string) error
}
```

#### CommentRepository

```go
type CommentRepository interface {
    FetchCommentsByBlogId(blogId string) ([]models.BlogCommentsData, error)
    CreateComment(blogId, guestUser, comment string) (*models.BlogCommentsData, error)
}
```

#### CookieUtils

```go
type CookieUtils interface {
    GetAuthCookie(c echo.Context, tokenName string) (*http.Cookie, error)
    GetAuthCookieValue(c echo.Context, tokenName string) (string, error)
    GetAuthCookieExpirationTime() time.Time
    ExistsAuthCookie(c echo.Context, tokenName string) bool
    VerifyToken(c echo.Context, tokenString string) (*models.Claims, error)
    CreateToken(user *models.BlogUsersData) (string, error)
    AddAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
    UpdateAuthCookie(c echo.Context, tokenString string, expirationTime time.Time)
    DelAuthCookie(c echo.Context)
    GetUserIdFromToken(c echo.Context, tokenString string) (string, error)
    CreateVisitIdToken() (string, error)
    AddVisitIdCoookie(c echo.Context, tokenString string, expirationTime time.Time)
    GetVisitIdFromToken(c echo.Context, tokenString string) (string, error)
}
```

### 3.3 依存関係の組み立て

依存関係は `routes/routes.go` の `SetupRoutes` 関数で組み立てられる。

```
CookieUtils ─────────────────────────────────┐
                                              │
BlogUsersRepository ──> UserService ──────────┤──> AuthHandler
                                     │        │
AuthService ─────────────────────────┘        │
                                              │
BlogUsersRepository ──> UserService ──────────┤──> BlogUsersHandler
                                              │
CookieUtils ─────────────────────────────────┤
                                              │
BlogRepository ──> BlogService ──────────────┤──> BlogHandler
                                              │
CookieUtils ─────────────────────────────────┤
                                              │
BlogLikeRepository ──> BlogLikeService ──────┤──> BlogLikeHandler
                                              │
CookieUtils ─────────────────────────────────┤
                                              │
CommentRepository ──> CommentService ─────────┘──> CommentHandler
```

---

## 4. 技術スタック

### 4.1 アプリケーション

| 分類 | 技術 | バージョン | 用途 |
|------|------|----------|------|
| 言語 | Go | 1.20 | メイン開発言語 |
| Webフレームワーク | Echo | v4.12.0 | HTTPサーバー、ルーティング、ミドルウェア |
| DBドライバ | pgx/v5 | v5.11.0 | PostgreSQL接続（コネクションプール。Simple Protocol で実行。#153） |
| 認証 | golang-jwt/jwt/v5 | v5.3.1 | JWT トークン生成・検証（`utils/cookie/jwt.go` に集約。#152） |
| 環境変数 | godotenv | v1.5.1 | .env ファイル読み込み |
| UUID | google/uuid | v1.6.0 | UUID生成・バリデーション |
| テスト | testify | v1.9.0 | アサーション・モック |
| 暗号化 | golang.org/x/crypto | v0.54.0 | `bcrypt` でパスワードのハッシュ化・照合を行う（`services/blog_users/password.go`。#120） |

### 4.2 インフラストラクチャ

| 分類 | 技術 | 用途 |
|------|------|------|
| コンテナ実行環境 | Google Cloud Run | サーバーレスコンテナ実行 |
| コンテナレジストリ | Artifact Registry | Docker イメージ管理 |
| シークレット管理 | Secret Manager | 環境変数の安全な管理 |
| データベース | Supabase（PostgreSQL） | データ永続化 |
| IaC | Terraform | インフラ構成管理 |
| CI/CD | GitHub Actions | 自動ビルド・デプロイ |

### 4.3 ミドルウェア構成

| ミドルウェア | 設定内容 |
|-------------|---------|
| HTTPErrorHandler | `e.HTTPErrorHandler` に設定。エラーを統一エラーレスポンスに変換する（2.3 参照） |
| Logger | リクエストログの自動出力（エラー時は 4xx を WARN、5xx を ERROR で記録） |
| Recover | パニック時の自動リカバリ |
| CORS | `ALLOWED_ORIGINS` 環境変数で許可オリジンを制御。AllowCredentials: true |
| RateLimiter（全体） | クライアント IP 単位で毎秒 20 件・瞬間 60 件まで（ヘルスチェック `GET /` を除く）。CORS の後に置き、429 にも CORS ヘッダを付ける |
| RateLimiter（認証系） | ルート単位で `POST /api/users/login`・`PUT /api/users/update` に適用（5 回/分、2 つでカウントを共有）。`routes.go` で `middlewares.NewAuthRateLimiter()` を渡す |

- クライアント IP は `e.IPExtractor = middlewares.NewIPExtractor(...)` で取る。`TRUST_CLIENT_IP_HEADER=true`（本番）のときは BFF が付与した `X-Client-IP`、それ以外は `X-Forwarded-For` の右端側（`echo.ExtractIPFromXFFHeader()`）。カウントはインスタンスのメモリで持つ。詳細は [06-security-specification.md](06-security-specification.md) 8.2。

#### CORS設定詳細

```
AllowOrigins: ALLOWED_ORIGINS 環境変数（カンマ区切り）
AllowMethods: GET, POST, PUT, DELETE
AllowHeaders: Origin, Content-Type, Authorization, Access-Control-Allow-Credentials
AllowCredentials: true
```

---

## 5. インフラストラクチャ

### 5.1 GCPリソース構成

Terraformで管理される GCP リソースは以下の通り。

#### Cloud Run サービス

| 設定項目 | 値 |
|---------|-----|
| CPU | 1000m (1 vCPU) |
| メモリ | 512Mi |
| ポート | 環境変数 `api_port` で指定 |
| 認証 | 未認証アクセス拒否（`--no-allow-unauthenticated`・`--invoker-iam-check`）。Invoker は BFF のサービスアカウントのみ（`invoker_member`。#146 / #165） |
| サービスアカウント | `cloud-run-sa` |
| トラフィック | 最新リビジョンに100% |

#### Artifact Registry

| 設定項目 | 値 |
|---------|-----|
| フォーマット | DOCKER |
| リポジトリID | 変数 `repository_id` で指定 |
| 説明 | Docker repository for Echo Blog Backend App |

#### Secret Manager

以下のシークレットが管理されている。

| シークレット名 | 用途 |
|---------------|------|
| `allowed_origins` | CORS許可オリジン |
| `env_word` | 環境識別子 |
| `jwt_secret_key` | JWT署名キー |
| `supabase_url` | Supabase接続URL |
| `tmp_allowed_origins` | 一時CORS許可オリジン |
| `gcp_project_id` | GCPプロジェクトID |

#### IAMロール

- Cloud Run サービスに対する呼び出し権限（`invoker_role`）
- Cloud Run サービスアカウントに対するSecret Manager読み取り権限（`roles/secretmanager.secretAccessor`）

### 5.2 Terraform構成

```
terraform/
├── main.tf                 # Terraform設定、Providerの定義
├── variables.tf            # 変数定義（GCPプロジェクトID、リージョン等）
├── cloud-run.tf            # Cloud Runサービス、サービスアカウント
├── cloud-run-iam.tf        # IAMロール設定
├── artifact-registry.tf    # Artifact Registry リポジトリ
└── secrets.tf              # Secret Manager シークレット
```

#### Terraform要件

| 項目 | 値 |
|-----|-----|
| Terraform | >= 1.6 |
| Google Provider | ~> 5.0 |

---

## 6. CI/CDパイプライン

GitHub Actions は 3 つのワークフローに責務を分離している。

| ファイル | トリガー | 役割 |
|---------|---------|------|
| `.github/workflows/ci.yml` | PR・`main` への push | 品質ゲート（lint / 静的解析 + UT + IT + E2E）。マージ前に PR で実行される |
| `.github/workflows/deploy.yml` | `main` への push | Docker ビルド → Artifact Registry へ push → Cloud Run へデプロイ |
| `.github/workflows/secret-scan.yml` | PR・`main` への push（パスフィルタなし） | 秘匿ファイル（鍵・`.env`・Terraform の state / 変数）が Git の追跡対象に入っていないかを検査する（6.0.2 参照） |

> テスト実行は `ci.yml` に一本化している。`deploy.yml` はデプロイ専任で、テスト・Go セットアップは持たない（PR 段階で `ci.yml` が品質を担保済みのため）。

### 6.0 CI ワークフロー（品質ゲート）

ファイル: `.github/workflows/ci.yml`

PR・`main` への push で実行し、独立した 4 ジョブを並列に走らせる。

| ジョブ | ステップ | 説明 |
|-------|---------|------|
| `lint` | gofmt（`gofmt -l .`）/ go vet / golangci-lint（v2.12.2） | フォーマット・静的解析・統合リンタ |
| `test` | `go test ./...`（`working-directory: backend`、`JWT_SECRET_KEY` にテスト用ダミー値を env で注入） | 全パッケージの UT を実行（全 mock。DB を必要としない）。対象を列挙しないため、追加したパッケージも自動で対象になる（#158）。IT / E2E は build tag で分離しているため含まれない |
| `integration-test` | `go test -tags=integration ./repositories/... ./testsupport/...`（`working-directory: backend`） | Repository 層の IT を testcontainers（使い捨て PostgreSQL）で実行。実 Supabase には接続しない。`testsupport` の接続先検査とスキーマの検査（監査列・`updated_at` のトリガ。#117）のテストも実行 |
| `e2e-test` | `go test -tags=e2e ./e2e/...`（`working-directory: backend`、`JWT_SECRET_KEY` にテスト用ダミー値を env で注入） | httptest で起動した実サーバーへ HTTP を投げ、Handler→Service→Repository→DB を貫く API-E2E を testcontainers 上で実行 |

> IT・E2E とも `testsupport` パッケージが testcontainers で PostgreSQL コンテナを起動し、`testsupport/testdata/schema.sql` / `seed.sql` を適用して実行する。GitHub Actions の ubuntu ランナーは Docker 同梱のため追加設定は不要。`SUPABASE_URL` を環境変数で指定した場合のみ、コンテナを起動せずその DB に接続する（接続先はローカルに限り、本番等のリモートは `testsupport` が拒否する。#140）。UT・E2E は `config` パッケージが起動時に `JWT_SECRET_KEY` を要求するため env で注入する。値は本番の secret ではなく CI 専用のダミー値（`ci-test-dummy-jwt-secret`）とする（Dependabot が起動したワークフローは Actions secrets を参照できず空になるため、また本番鍵をテストに露出させないため）。

### 6.0.1 依存関係の自動更新（Dependabot）

ファイル: `.github/dependabot.yml`

| エコシステム | 対象 | 更新方針 |
|-------------|------|---------|
| `github-actions` | `.github/workflows/*.yml` のアクション | 週次でバージョン更新。全アクションを 1 つの PR（グループ `github-actions`）にまとめる |
| `gomod` | `backend/go.mod` | **セキュリティ更新のみ**（`open-pull-requests-limit: 0` でバージョン更新 PR は無効化）。セキュリティ更新は 1 つの PR（グループ `go-security`）にまとめ、`go.sum` の連鎖コンフリクトを防ぐ |

> セキュリティ更新はリポジトリ設定（Dependabot security updates）で有効化されている。Dependabot が起動した CI は Actions secrets を参照できないため、CI はテスト用ダミー値で完結させている（6.0 参照）。Dependabot の PR は手で push せず、コンフリクト時は `@dependabot rebase` で解消する（手で push すると以後 Dependabot が当該 PR を更新しなくなる）。

### 6.0.2 Secret scan（秘匿ファイルの混入検出）

ファイル: `.github/workflows/secret-scan.yml`（判定ロジックは `scripts/check-secret-files.sh`、そのテストは `scripts/check-secret-files_test.sh`）

`git ls-files` で引いた「git が見ているファイル」（追跡済み + 未追跡。`.gitignore` 対象は除く）を検査し、秘匿ファイルが含まれていたら失敗させる（#126）。`.gitignore` は未追跡ファイルにしか効かず（`git add -f` や書き漏れを止められない）、一度 push したファイルは追跡除外しても履歴に残るため、「追跡された時点で落とす」検出側として置く。

| ステップ | 説明 |
|---------|------|
| `scripts/check-secret-files_test.sh` | 検出ルールの退行を、実スキャンの前に検知する（ケースごとに使い捨ての git リポジトリを作って実行する） |
| `scripts/check-secret-files.sh` | 追跡済みのファイル（混入）と、`.gitignore` されていない未追跡のファイル（`git add` すると混入する）を検査する。`git ls-files` はインデックスを読むだけで履歴を走査しないため数秒で終わる |

| 検査 | 対象 |
|------|------|
| ファイル名 | `.env` / `.env.*`、`*.key` / `*.pem` / `*.p12` / `*.pfx` / `*.jks` / `*.keystore`、`id_rsa` / `id_ed25519` / `id_dsa`、`credentials.json` / `serviceAccountKey.json`、Terraform の `*.tfstate` / `*.tfstate.backup` / `*.tfvars`（state には Secret Manager に入れる値が平文で残る） |
| ファイル名（除外） | `*.example` / `*.sample` / `*.template` / `*.dist`、`*.env.d.ts`（テンプレートと型定義は秘匿情報を含まない） |
| 内容 | PEM 秘密鍵のヘッダ（`-----BEGIN … PRIVATE KEY-----`）、GCP サービスアカウントキー JSON の `private_key` フィールド。サービスアカウントキーは任意のファイル名で保存されるため、名前では捕まらない |

> `ci.yml` に同居させず独立したワークフローにしているのは、`ci.yml` がワークフロー単位の `paths`（`backend/**`）で絞っており、ルートや `terraform/` への混入では起動しないため。秘匿ファイルはどの変更種別でも混入しうるため、`secret-scan.yml` はパスフィルタを付けずに常に実行する。常に起動するので、ブランチ保護の必須チェックにしても `pending` のまま詰まらない。
>
> 検出時の対処は追跡除外ではなく、鍵・トークンのローテーションになる（push 済みの内容は履歴から取得できるため）。ローカルでは `scripts/check-secret-files.sh` で同じ検査を実行できる（`/pr-create`・`/self-review` もこのスクリプトがあれば実行する）。履歴上のパスは `scripts/check-secret-files.sh --history`（ファイル名のみ。導入時に一度だけ実行し、0 件を確認済み）で照合する。

### 6.1 デプロイワークフロー

ファイル: `.github/workflows/deploy.yml`

#### トリガー

- `main` ブランチへの `push` 時に自動実行
- ただし `paths` フィルタにより、`backend/**/*.go` / `backend/go.mod` / `backend/go.sum` / `backend/Dockerfile` / `.github/workflows/**` のいずれかが変更された場合のみ実行（Go アプリは `backend/` 配下に集約されているため `backend/` プレフィックスで限定）
- `_test.go` は `paths` 内の否定パターン `'!backend/**/*_test.go'` で除外（テストファイルのみの変更ではデプロイされない）

> 注意: GitHub Actions は同一イベントで `paths` と `paths-ignore` を併用できない（併用するとワークフローが startup failure になる）。テスト除外は `paths-ignore` ではなく `paths` 内の否定パターンで表現している。

#### パイプラインフロー

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Checkout    │───>│  GCP認証     │───>│  Docker認証  │───>│  Build &     │───>│  Deploy to   │
│  code        │    │              │    │  (GCR)       │    │  Push Image  │    │  Cloud Run   │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘    └──────┬───────┘
                                                                                       │
                                                                                ┌──────┴───────┐
                                                                                │  Cleanup old │
                                                                                │  images      │
                                                                                └──────────────┘
```

> テストは `ci.yml`（PR・push 時）で実行済みのため、`deploy.yml` にテストステップは持たない。

#### ステップ詳細

| ステップ | Action/コマンド | 説明 |
|---------|----------------|------|
| 1. Checkout | `actions/checkout@v3` | ソースコードの取得 |
| 2. GCP認証 | `google-github-actions/auth@v1` | サービスアカウントキー（`credentials_json`）で認証 |
| 3. Docker認証 | `gcloud auth configure-docker` | Artifact Registryへの認証設定 |
| 4. GCloud設定 | `google-github-actions/setup-gcloud@v1` | GCloud SDKの初期化 |
| 5. ビルド・プッシュ | `docker build ./backend` + `docker push` | イメージビルドとプッシュ（ビルドコンテキストは `backend/`。タグ: コミットSHA） |
| 6. デプロイ | `gcloud run deploy` | Cloud Runへのデプロイ（`--max-instances 3`。レートリミットの実効上限とコストを固定する。`--no-allow-unauthenticated`・`--invoker-iam-check` と `TRUST_CLIENT_IP_HEADER=true` を同じコマンドで設定する。#146 / #165） |
| 7. クリーンアップ | `gcloud artifacts docker images delete` | 古いイメージの削除（最新5件を保持） |

#### 使用するGitHub Secrets

| シークレット名 | 用途 |
|---------------|------|
| `GCP_SERVICE_ACCOUNT_KEY` | GCPサービスアカウントの認証キー（JSON） |
| `GCP_PROJECT_ID` | GCPプロジェクトID |
| `GCP_REGION` | GCPリージョン |
| `REPO_NAME` | Artifact Registryリポジトリ名 |
| `APP_NAME` | アプリケーション名 |
| `GCP_CLOUD_RUN_SERVICE_NAME` | Cloud Runサービス名 |

---

## 7. デプロイメント

### 7.1 Dockerイメージ構成

マルチステージビルドを採用し、最終イメージサイズを最小化している。`Dockerfile` は `backend/Dockerfile` に配置し、ビルドコンテキストは `backend/` を指定する（CI の `docker build ./backend`）。`COPY . .` はこのコンテキスト（＝`backend/` の中身）をコピーする。

#### Dockerfile（backend/Dockerfile）

```dockerfile
# ビルドステージ
FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

# 実行ステージ
FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
```

| ステージ | ベースイメージ | 用途 |
|---------|--------------|------|
| builder | `golang:1.25` | Goアプリケーションのコンパイル |
| runtime | `gcr.io/distroless/base` | 最小限のランタイム環境 |

#### Distrolessイメージの特徴

- シェルやパッケージマネージャを含まない最小構成
- セキュリティリスクの低減
- イメージサイズの削減

### 7.2 環境変数

アプリケーション起動時に必要な環境変数。Cloud Run上ではSecret Managerから注入される。

| 環境変数 | 説明 | 必須 |
|---------|------|------|
| `SUPABASE_URL` | Supabase（PostgreSQL）接続URL | はい |
| `DB_SSLMODE` | 接続時の SSL モード（未設定時 `require`） | いいえ |
| `JWT_SECRET_KEY` | JWT署名キー | はい |
| `ALLOWED_ORIGINS` | CORS許可オリジン（カンマ区切り） | はい |
| `ENV` | 環境識別子（`production` で本番Cookie設定） | はい |
| `TRUST_CLIENT_IP_HEADER` | `true` で BFF が付与する `X-Client-IP` をクライアント IP として信頼する（非公開化が前提。`deploy.yml` と Terraform で設定） | 本番のみ |
| `PORT` | サーバーリスンポート（デフォルト: 8080） | いいえ |

### 7.3 サーバー起動フロー

```
main()
  ├── firstSetup()
  │     ├── godotenv.Load()          # .envファイル読み込み（存在する場合）
  │     ├── logger.InitLogger()      # ログ設定初期化
  │     ├── supabase.InitSupabase()  # DB接続プール初期化
  │     └── supabase.TestQuery()     # DB接続テスト
  │
  ├── echo.New()                     # Echoインスタンス生成
  ├── middlewares.SetupMiddlewares() # ミドルウェア設定
  ├── routes.SetupRoutes()           # ルーティング・DI設定
  │
  ├── signal.Notify()                # シグナルハンドラ設定（SIGINT, SIGTERM）
  │     ├── e.Close()               # Echoサーバーシャットダウン
  │     └── supabase.ClosePool()    # DB接続プールクローズ
  │
  └── e.Start(":PORT")              # サーバー起動
```

---

## 8. ディレクトリ構成

本プロジェクトはモノレポ構成を採用しており、Go アプリケーション一式は `backend/` 配下に集約している。CI（`.github/`）・インフラ（`terraform/`）・ドキュメント（`docs/`）はリポジトリルートで全体共通として管理する。

**リポジトリ全体**

```
nextjs-echo-back-blog-app/           # リポジトリルート（モノレポ）
├── backend/                         # Go + Echo アプリ一式（module backend）。詳細は下記
├── terraform/                       # インフラ構成管理（Terraform）
├── docs/                            # 仕様・設計ドキュメント
├── .github/                         # CI/CD（GitHub Actions）
├── scripts/                         # CI 用スクリプト（Secret scan とそのテスト）
├── manuals/                         # マニュアル
├── backup/                          # 旧 AWS App Runner 構成のバックアップ
├── work/                            # 作業用（gitignore）
├── readme.md                        # プロジェクトREADME
├── CLAUDE.md                        # 開発ルール索引
└── LICENSE
```

**`backend/` 配下の詳細**

```
backend/
├── main.go                          # エントリポイント
├── go.mod                           # Go modules 設定（module backend）
├── go.sum                           # 依存関係ロックファイル
├── Dockerfile                       # マルチステージDockerビルド
│
├── apperrors/
│   ├── apperrors.go                 # エラーの分類（Kind）・コード・共通 sentinel（2.3 参照）
│   └── apperrors_test.go            # errors.Is / Wrap の UT
│
├── config/
│   └── config.go                    # JWT鍵・環境判定の設定
│
├── dto/                             # レスポンス DTO。モデルを直接返さない（JSON タグは DTO 側で管理）
│   ├── error.go                     # ErrorResponse（統一エラーレスポンス）
│   ├── blog_users.go                # UserResponse（パスワードを含めない）
│   ├── blogs.go                     # BlogResponse
│   ├── blog_likes.go                # BlogLikeResponse
│   ├── blog_comments.go             # CommentResponse
│   └── *_test.go                    # DTO 化以前の JSON との一致（ゴールデン）を検証
│
├── logger/
│   └── logger.go                    # ログレベル設定（Info/Error/Warn/Debug/Test）
│
├── middlewares/
│   ├── middlewares.go               # Echo ミドルウェア設定（HTTPErrorHandler/IPExtractor/RequestLogger/Recover/CORS/RateLimiter）
│   ├── ip_extractor.go              # クライアント IP の抽出（X-Client-IP / X-Forwarded-For）
│   ├── ip_extractor_test.go         # IP 抽出・TRUST_CLIENT_IP_HEADER の UT
│   ├── middlewares_test.go          # アクセスログミドルウェアの UT
│   ├── error_handler.go             # HTTPErrorHandler（エラー → ステータス・統一エラーレスポンス）
│   ├── error_handler_test.go        # エラー変換の UT
│   ├── rate_limiter.go              # レートリミット（全体・認証系）
│   └── rate_limiter_test.go         # レートリミット（BFF 経由の数え方を含む）の UT
│
├── models/
│   ├── auth.go                      # JWT Claims 構造体
│   ├── blogs.go                     # BlogData 構造体
│   ├── blog_users.go                # BlogUsersData 構造体
│   ├── blog_likes.go                # BlogLikesData 構造体
│   └── blog_comments.go            # BlogCommentsData 構造体
│
├── routes/
│   ├── routes.go                    # ルーティング定義、DI組み立て
│   └── routes_test.go               # 認証系レートリミットの適用先の UT
│
├── handlers/                        # Handler層（HTTPリクエスト処理）
│   ├── auth/
│   │   ├── auth.go                  # AuthHandler 構造体・コンストラクタ
│   │   ├── auth_impl.go            # Login, CheckAuth, Logout 実装
│   │   ├── auth_testing_setup.go   # テストセットアップ
│   │   ├── auth_test.go            # Handler テスト
│   │   ├── auth_Login_test.go
│   │   ├── auth_CheckAuth_test.go
│   │   ├── auth_Logout_test.go
│   │   └── auth_pipeline_test.go
│   ├── blog_users/
│   │   ├── blog_users.go           # BlogUsersHandler 構造体
│   │   ├── blog_users_impl.go      # FetchBlogUsers, UpdateBlogUsers 等
│   │   ├── blog_users_cookie_mock.go  # Cookie モックヘルパー
│   │   ├── blog_users_FetchUser_test.go
│   │   ├── blog_users_response_test.go  # レスポンス DTO の検証
│   │   └── blog_users_UpdateUser_test.go
│   ├── blogs/
│   │   ├── blogs.go                # BlogHandler 構造体
│   │   ├── blogs_impl.go           # FetchBlogs, CreateBlog 等
│   │   ├── blogs_cookie_mock.go    # Cookie モックヘルパー
│   │   ├── blogs_testing_setup.go  # テストセットアップ
│   │   └── test/                   # Handler テスト群
│   │       ├── main_test.go
│   │       ├── blogs_FetchBlogs_test.go
│   │       ├── blogs_FetchBlogById_test.go
│   │       ├── blogs_FetchBlogsByUserId_test.go
│   │       ├── blogs_FetchBlogCategories_test.go
│   │       ├── blogs_FetchBlogTags_test.go
│   │       ├── blogs_FetchBlogPopular_test.go
│   │       ├── blogs_CreateBlog_test.go
│   │       ├── blogs_UpdateBlog_test.go
│   │       └── blogs_DeleteBlog_test.go
│   ├── blog_likes/
│   │   ├── blog_likes.go           # BlogLikeHandler 構造体
│   │   ├── blog_likes_impl.go      # FetchBlogLikesByVisitId 等
│   │   ├── blog_likes_cookie_mock.go
│   │   ├── blog_likes_FetchBlogLikesByVisitId_test.go
│   │   ├── blog_likes_IsBlogLiked_test.go
│   │   ├── blog_likes_CreateBlogLike_test.go
│   │   ├── blog_likes_DeleteBlogLike_test.go
│   │   └── blog_likes_GenerateVisitorId_test.go
│   └── blog_comments/
│       ├── blog_comments.go        # CommentHandler 構造体
│       ├── blog_comments_impl.go   # FetchCommentsByBlogId, CreateComment
│       ├── blog_comments_FetchCommentsByBlogId_test.go
│       └── blog_comments_CreateComment_test.go
│
├── services/                        # Service層（ビジネスロジック）
│   ├── auth/
│   │   ├── auth.go                 # AuthService インターフェース
│   │   ├── auth_impl.go            # Login バリデーション実装
│   │   ├── errors.go               # サービスが返す sentinel error
│   │   ├── auth_mock.go            # モック
│   │   └── auth_test.go            # Service テスト
│   ├── blog_users/
│   │   ├── blog_users.go           # UserService インターフェース
│   │   ├── blog_users_impl.go      # FetchUser, UpdateUser 等（不正な入力でテスト）
│   │   ├── errors.go               # サービスが返す sentinel error
│   │   ├── blog_users_mock.go      # モック
│   │   ├── blog_users_FetchUserById_test.go
│   │   ├── blog_users_FetchUserByEmailAndPassword_test.go
│   │   └── blog_users_UpdateUser_test.go
│   ├── blogs/
│   │   ├── blogs.go                # BlogService インターフェース
│   │   ├── blogs_impl.go           # ビジネスロジック実装
│   │   ├── errors.go               # サービスが返す sentinel error
│   │   ├── blogs_mock.go           # モック
│   │   ├── blogs_testing_setup.go  # テストセットアップ
│   │   └── test/                   # Service テスト群
│   │       ├── main_test.go
│   │       ├── blogs_FetchBlogs_test.go
│   │       ├── blogs_FetchBlogById_test.go
│   │       ├── blogs_FetchBlogsByUserId_test.go
│   │       ├── blogs_FetchBlogCategories_test.go
│   │       ├── blogs_FetchBlogTags_test.go
│   │       ├── blogs_FetchBlogPopular_test.go
│   │       ├── blogs_CreateBlog_test.go
│   │       ├── blogs_UpdateBlog_test.go
│   │       └── blogs_DeleteBlog_test.go
│   ├── blog_likes/
│   │   ├── blog_likes.go           # BlogLikeService インターフェース
│   │   ├── blog_likes_impl.go      # いいねロジック実装
│   │   ├── errors.go               # サービスが返す sentinel error
│   │   ├── blog_likes_mock.go      # モック
│   │   ├── blog_likes_FetchBlogLikesByVisitId_test.go
│   │   ├── blog_likes_IsBlogLiked_test.go
│   │   ├── blog_likes_CreateBlogLike_test.go
│   │   └── blog_likes_DeleteBlogLike_test.go
│   └── blog_comments/
│       ├── blog_comments.go        # CommentService インターフェース
│       ├── blog_comments_impl.go   # コメントロジック実装
│       ├── errors.go               # サービスが返す sentinel error
│       ├── blog_comments_mock.go   # モック
│       ├── blog_comments_FetchCommentsByBlogId_test.go
│       └── blog_comments_CreateComment_test.go
│
├── repositories/                    # Repository層（データアクセス）
│   ├── blog_users/
│   │   ├── blog_users.go           # BlogUsersRepository インターフェース
│   │   ├── blog_users_impl.go      # SQL実装
│   │   ├── blog_users_mock.go      # モック
│   │   ├── main_test.go            # TestMain（testsupport 経由でIT用DB準備）
│   │   ├── blog_users_FetchUserById_test.go
│   │   ├── blog_users_FetchUserByEmailAndPassword_test.go
│   │   ├── blog_users_UpdateBlogUsers_test.go
│   │   └── zz_connection_closed_test.go  # 異常系（接続断で安全に失敗）
│   ├── blogs/
│   │   ├── blogs.go                # BlogRepository インターフェース
│   │   ├── blogs_impl.go           # SQL実装（JOIN、集計クエリ含む）
│   │   ├── blogs_mock.go           # モック
│   │   └── test/                   # Repository テスト群（IT）
│   │       ├── main_test.go        # TestMain（testsupport 経由でIT用DB準備）
│   │       ├── blogs_FetchBlogs_test.go
│   │       ├── blogs_FetchBlogById_test.go
│   │       ├── blogs_FetchBlogsByUserId_test.go
│   │       ├── blogs_FetchBlogCategories_test.go
│   │       ├── blogs_FetchBlogTags_test.go
│   │       ├── blogs_FetchBlogPopular_test.go
│   │       ├── blogs_CreateBlog_test.go
│   │       ├── blogs_UpdateBlog_test.go
│   │       ├── blogs_DeleteBlog_test.go
│   │       ├── blogs_pipeline_test.go
│   │       └── zz_connection_closed_test.go  # 異常系（接続断で安全に失敗）
│   ├── blog_likes/
│   │   ├── blog_likes.go           # BlogLikeRepository インターフェース
│   │   ├── blog_likes_impl.go      # SQL実装
│   │   ├── blog_likes_mock.go      # モック
│   │   ├── main_test.go            # TestMain（testsupport 経由でIT用DB準備）
│   │   ├── blog_likes_pipeline_test.go
│   │   └── zz_connection_closed_test.go  # 異常系（接続断で安全に失敗）
│   └── blog_comments/
│       ├── blog_comments.go        # CommentRepository インターフェース
│       ├── blog_comments_impl.go   # SQL実装
│       ├── blog_comments_mock.go   # モック
│       ├── main_test.go            # TestMain（testsupport 経由でIT用DB準備）
│       ├── blog_comments_CreateComment_test.go
│       ├── blog_comments_FetchCommentsByBlogId_test.go
│       └── zz_connection_closed_test.go  # 異常系（接続断で安全に失敗）
│
├── supabase/                        # Supabase接続管理
│   └── client.go                   # 接続プール初期化・テストクエリ・クローズ
│
├── testsupport/                     # IT/E2E 共有ヘルパー（build tag: integration || e2e）
│   ├── testdb.go                   # testcontainers 起動 or 実DB接続の切替（Setup / Start）
│   └── testdata/
│       ├── schema.sql              # スキーマ（コードのテーブル名に一致）
│       └── seed.sql                # 決定的シード（固定UUID）
│
├── e2e/                             # API-E2E テスト（build tag: e2e）
│   ├── doc.go                      # パッケージ宣言（デフォルトビルドで有効化するプレースホルダ）
│   ├── main_test.go                # TestMain（DB準備 + httptest サーバ起動）
│   ├── helpers_test.go             # HTTP クライアント/リクエストヘルパ
│   ├── health_test.go              # ヘルスチェック
│   ├── auth_flow_test.go           # ログイン→認証確認→ログアウト
│   ├── blogs_read_test.go          # ブログ参照系
│   ├── blogs_crud_test.go          # ブログ CRUD（認証あり）
│   ├── likes_comments_test.go      # いいね/コメントフロー
│   ├── users_test.go               # ユーザー取得（パスワード非露出）
│   ├── rate_limit_test.go          # レートリミット
│   └── errors_test.go              # 統一エラーレスポンス・ステータス
│
├── utils/                           # ユーティリティ
│   ├── cookie/
│   │   ├── cookie_impl.go          # CookieUtils インターフェース定義
│   │   ├── cookie_common_di.go     # 共通Cookie操作
│   │   ├── cookie_token_di.go      # 認証トークンCookie操作
│   │   ├── cookie_visitId_di.go    # 訪問者IDトークンCookie操作
│   │   ├── cookie_mock.go          # CookieUtils モック
│   │   └── utils_cookie.go         # Cookie追加・削除のスタティック関数
│   └── log/
│       ├── utils_log.go            # ログユーティリティ
│       └── utils_log_test.go       # ログユーティリティテスト
```

> `terraform/` ・ `.github/` ・ `docs/` ・ `manuals/` ・ `backup/` ・ `work/` はリポジトリルート直下（`backend/` と同階層）。構成は上記「リポジトリ全体」ツリーを参照。

#### インフラ（terraform/）の構成

```
terraform/
├── main.tf                     # Provider設定
├── variables.tf                # 変数定義
├── cloud-run.tf                # Cloud Runサービス定義
├── cloud-run-iam.tf            # IAMロール設定
├── artifact-registry.tf        # Artifact Registry定義
└── secrets.tf                  # Secret Manager定義
```

---

## 9. セキュリティ設計

### 9.1 認証・認可

| 項目 | 実装 |
|------|------|
| 認証方式 | JWT（HS256） |
| トークン保存 | HTTP-Only Cookie |
| トークン有効期限 | 1時間 |
| Cookie設定（本番） | Secure: true, SameSite: None, HttpOnly: true |
| Cookie設定（開発） | Secure: false, SameSite: Lax, HttpOnly: true |

### 9.2 環境識別

`ENV` 環境変数が `production` の場合に本番用Cookie設定が適用される。

```go
var IsProduction = os.Getenv("ENV") == "production"
```

### 9.3 シークレット管理

- 環境変数はGCP Secret Managerで管理
- Cloud Runサービスアカウントにのみ `secretmanager.secretAccessor` ロールを付与
- `JWT_SECRET_KEY` が未設定の場合、アプリケーション起動時に `log.Fatal` で即座に停止

### 9.4 コンテナセキュリティ

- Distrolessベースイメージにより攻撃対象を最小化
- シェルアクセス不可
- 不要なOSパッケージを含まない

---

## 10. ローカル開発セットアップ

ゼロから動かすための手順。コマンドの正準はリポジトリ直下の [`README.md`](../README.md) にも記載しており、本節はその背景・補足を含む詳細版である。

### 10.1 前提ツール

| ツール | バージョン | 用途 |
|--------|-----------|------|
| Go | 1.25 以上（`go.mod` は `go 1.25.0`） | アプリケーションのビルド・実行 |
| Supabase / PostgreSQL | - | 稼働中の接続先（無料枠で可）。起動時に接続テスト `SELECT 1` を実行する |
| Docker | アプリのコンテナ起動時、および IT（testcontainers）実行時に必須 | Repository 層 IT はコンテナで PostgreSQL を起動する |

> `go.mod`（`go 1.25.0`）と `Dockerfile` のビルドステージ（`golang:1.25`）はバージョンを一致させている。依存更新で `go.mod` の `go` 指令が引き上げられた場合は、`Dockerfile` も同一 PR で揃える（公式 golang イメージは `GOTOOLCHAIN=local` のため新しいツールチェーンを自動取得しない）。

### 10.2 セットアップ手順

Go アプリは `backend/` 配下に集約されている（モノレポ構成）。以降の `go` コマンド・`.env` は `backend/` ディレクトリ内で扱う。

```bash
git clone https://github.com/kojikawazu/nextjs-echo-back-blog-app.git
cd nextjs-echo-back-blog-app/backend

# 環境変数テンプレートをコピーして値を設定
cp .env.example .env

# 依存解決して起動
go mod download
go run main.go
```

起動フローは [§7.3 サーバー起動フロー](#73-サーバー起動フロー) を参照。`http://localhost:8080/` が `Service is running` を返せば成功。

### 10.3 環境変数

アプリ起動時に `godotenv` が `backend/.env` を読み込む（`backend/main.go` の `firstSetup`）。テンプレートは [`../backend/.env.example`](../backend/.env.example)。参照コードのパスはいずれも `backend/` 配下。

| 変数名 | 必須 | 説明 | 参照コード |
|--------|------|------|-----------|
| `JWT_SECRET_KEY` | ✅ | JWT署名鍵。未設定だと `config.init()` の `log.Fatal` で即終了 | `backend/config/config.go:9,12-15` |
| `SUPABASE_URL` | ✅ | PostgreSQL接続URL。`?sslmode=<DB_SSLMODE>` はコードが自動付与 | `backend/supabase/client.go` |
| `DB_SSLMODE` | - | 接続時の SSL モード。未設定時 `require`（本番）。非SSLのDBは `disable` | `backend/supabase/client.go` |
| `ALLOWED_ORIGINS` | ✅(ブラウザ利用時) | CORS許可オリジン（カンマ区切り） | `backend/middlewares/middlewares.go:29` |
| `ENV` | - | `production` で Cookie が Secure/SameSite=None。ローカルは空 | `backend/config/config.go:10` |
| `TRUST_CLIENT_IP_HEADER` | - | `true` で `X-Client-IP` を信頼する。ローカルでは設定しない（偽装でレートリミットを回避できるため） | `backend/middlewares/middlewares.go:23` |
| `PORT` | - | リッスンポート（既定 8080） | `backend/main.go:74-77` |
| `TEST_MODE` | - | `true` でログ破棄（主にテスト用） | `backend/logger/logger.go` |

### 10.4 テスト実行

`go` コマンドは `backend/` 内で実行する。

```bash
cd backend

# 単体テスト（UT / DB 不要・全 mock。全パッケージ）
go test ./...

# インテグレーションテスト（IT / Repository 層）。既定で testcontainers が
# 使い捨ての PostgreSQL を起動する（Docker 稼働が前提）。
go test -tags=integration ./repositories/...

# E2E（API-E2E）。実サーバーを起動し HTTP フローを検証（Docker 稼働が前提）。
# config が起動時に JWT_SECRET_KEY を要求するため env で渡す。
JWT_SECRET_KEY=... go test -tags=e2e ./e2e/...
```

- IT は `//go:build integration`、E2E は `//go:build e2e` タグで分離されており、通常の `go test ./...` には含まれない。IT は各パッケージの `TestMain` が `backend/testsupport` の `Start` を、E2E は `main_test.go` が `Setup` を呼び、コンテナ起動・スキーマ/シード適用・接続を一括で行う。
- IT は `SUPABASE_URL` を環境変数として指定した場合のみ、コンテナを起動せずその DB に接続する（`TEST_*` も併せて指定すること）。`testsupport` は `.env` ファイルを自動読み込みしない（古い `.env.test` によるコンテナ経路の意図しない乗っ取りを防ぐため）。
- Service / Handler 層テストはモックを使うため DB 不要だが、`config` パッケージの `init` が `JWT_SECRET_KEY` を要求する。
- CI（`.github/workflows/ci.yml`）は PR・`main` への push で `test` ジョブ（UT）と `integration-test` ジョブ（IT / testcontainers）を `working-directory: backend` で実行する。詳細は [`08-test-specification.md`](08-test-specification.md)。

### 10.5 Docker での起動

ビルドコンテキストは `backend/` を指定する。

```bash
docker build -t echo-blog-back ./backend
docker run --env-file backend/.env -p 8080:8080 echo-blog-back
```

イメージ構成（マルチステージ / Distroless）は [§7.1 Dockerイメージ構成](#71-dockerイメージ構成) を参照。

### 10.6 よくあるつまずき

| 症状 | 原因 | 対処 |
|------|------|------|
| 起動直後に `JWT_SECRET_KEY is not set` で落ちる | `.env` 未作成 or `JWT_SECRET_KEY` 空 | `.env` に値を設定する |
| `unable to connect to Supabase` で `log.Fatal` | `SUPABASE_URL` 不正 / DB 未到達 | 接続URL・ネットワーク・DB稼働を確認 |
| ブラウザから CORS エラー | `ALLOWED_ORIGINS` にフロントの URL が無い | フロントのオリジンをカンマ区切りで追加 |
| Cookie が保存されない（本番） | `ENV=production` 時は HTTPS 必須（Secure=true） | HTTPS 経由でアクセスする |
