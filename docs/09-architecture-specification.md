# アーキテクチャ仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo）のシステムアーキテクチャ、レイヤー構成、インフラストラクチャ、CI/CDパイプラインを定義するアーキテクチャ仕様書です。コードベースのリバースエンジニアリングに基づいて作成しています。

## 目次

- [1. システムアーキテクチャ概要](#1-システムアーキテクチャ概要)
- [2. レイヤーアーキテクチャ](#2-レイヤーアーキテクチャ)
  - [2.1 レイヤー構成図](#21-レイヤー構成図)
  - [2.2 データフロー](#22-データフロー)
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

本プロジェクトはクリーンアーキテクチャに基づく3層構成を採用している。各層はインターフェースを介して疎結合に接続されている。

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
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                 Handler Layer                        │
│  ・HTTPリクエストのパース（パスパラメータ、ボディ）      │
│  ・Cookie/JWT トークンの取得                          │
│  ・Service 呼び出し                                   │
│  ・HTTPレスポンスの生成（ステータスコード、JSON）        │
│  ・エラーハンドリング（エラーメッセージとステータス変換） │
└──────────────────────┬──────────────────────────────┘
                       │  Interface (依存性逆転)
                       ▼
┌─────────────────────────────────────────────────────┐
│                 Service Layer                        │
│  ・入力バリデーション（空文字、メール形式）              │
│  ・ビジネスロジック（タグ分割、重複チェック等）          │
│  ・Repository 呼び出し                                │
│  ・エラーの変換とログ出力                              │
└──────────────────────┬──────────────────────────────┘
                       │  Interface (依存性逆転)
                       ▼
┌─────────────────────────────────────────────────────┐
│                Repository Layer                      │
│  ・SQLクエリの構築と実行                               │
│  ・Supabase（PostgreSQL）への接続プール管理             │
│  ・結果のスキャンとモデルへのマッピング                  │
│  ・JOINによるいいね数・コメント数の集計                  │
└──────────────────────┬──────────────────────────────┘
                       │  pgx/v4 pgxpool
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
2. ミドルウェアでロギング、リカバリ、CORS処理を実行
3. ルーターが対応するHandlerにリクエストを振り分け
4. Handlerがリクエストをパースし、必要に応じてCookieからJWTを検証
5. HandlerがServiceの対応メソッドを呼び出す
6. Serviceが入力バリデーションを実行後、Repositoryを呼び出す
7. RepositoryがSQLを実行しデータを取得/変更
8. 結果がService → Handler → クライアントへ返却される

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
| DBドライバ | pgx/v4 | v4.18.3 | PostgreSQL接続（コネクションプール） |
| 認証 | golang-jwt/jwt | v3.2.2 | JWT トークン生成・検証 |
| 環境変数 | godotenv | v1.5.1 | .env ファイル読み込み |
| UUID | google/uuid | v1.6.0 | UUID生成・バリデーション |
| テスト | testify | v1.9.0 | アサーション・モック |
| 暗号化 | golang.org/x/crypto | v0.22.0（indirect） | go.mod 上は間接依存。現状コードからの直接利用はなく、将来のパスワードハッシュ化（bcrypt）導入時の利用を想定 |

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
| Logger | リクエストログの自動出力 |
| Recover | パニック時の自動リカバリ |
| CORS | `ALLOWED_ORIGINS` 環境変数で許可オリジンを制御。AllowCredentials: true |

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
| 認証 | 未認証アクセス許可 (`--allow-unauthenticated`) |
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

### 6.1 GitHub Actions ワークフロー

ファイル: `.github/workflows/deploy.yml`

#### トリガー

- `main` ブランチへの `push` 時に自動実行
- ただし `paths` フィルタにより、`**/*.go` / `go.mod` / `go.sum` / `Dockerfile` / `.github/workflows/**` のいずれかが変更された場合のみ実行
- `paths-ignore` で `**/*_test.go` を除外（テストファイルのみの変更ではデプロイされない）

#### パイプラインフロー

```
┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  Checkout    │───>│  Setup Go &  │───>│  GCP認証     │───>│  Docker認証  │───>│  Build &     │───>│  Deploy to   │
│  code        │    │  Unit Test   │    │              │    │  (GCR)       │    │  Push Image  │    │  Cloud Run   │
└──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘    └──────┬───────┘
                                                                                                          │
                                                                                                   ┌──────┴───────┐
                                                                                                   │  Cleanup old │
                                                                                                   │  images      │
                                                                                                   └──────────────┘
```

#### ステップ詳細

| ステップ | Action/コマンド | 説明 |
|---------|----------------|------|
| 1. Checkout | `actions/checkout@v3` | ソースコードの取得 |
| 2. Go セットアップ | `actions/setup-go@v5`（`go-version-file: go.mod`） | go.mod のバージョンで Go をセットアップ |
| 3. ユニットテスト | `go test ./handlers/... ./services/...` | Handler/Service 層のテストを実行（`JWT_SECRET_KEY` を env で注入） |
| 4. GCP認証 | `google-github-actions/auth@v1` | サービスアカウントキー（`credentials_json`）で認証 |
| 5. Docker認証 | `gcloud auth configure-docker` | Artifact Registryへの認証設定 |
| 6. GCloud設定 | `google-github-actions/setup-gcloud@v1` | GCloud SDKの初期化 |
| 7. ビルド・プッシュ | `docker build` + `docker push` | イメージビルドとプッシュ（タグ: コミットSHA） |
| 8. デプロイ | `gcloud run deploy` | Cloud Runへのデプロイ |
| 9. クリーンアップ | `gcloud artifacts docker images delete` | 古いイメージの削除（最新5件を保持） |

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

マルチステージビルドを採用し、最終イメージサイズを最小化している。

#### Dockerfile

```dockerfile
# ビルドステージ
FROM golang:1.19 as builder
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
| builder | `golang:1.19` | Goアプリケーションのコンパイル |
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
| `JWT_SECRET_KEY` | JWT署名キー | はい |
| `ALLOWED_ORIGINS` | CORS許可オリジン（カンマ区切り） | はい |
| `ENV` | 環境識別子（`production` で本番Cookie設定） | はい |
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

```
nextjs-echo-back-blog-app/
├── main.go                          # エントリポイント
├── go.mod                           # Go modules 設定
├── go.sum                           # 依存関係ロックファイル
├── Dockerfile                       # マルチステージDockerビルド
├── readme.md                        # プロジェクトREADME
│
├── config/
│   └── config.go                    # JWT鍵・環境判定の設定
│
├── logger/
│   └── logger.go                    # ログレベル設定（Info/Error/Warn/Debug/Test）
│
├── middlewares/
│   └── middlewares.go               # Echo ミドルウェア設定（Logger/Recover/CORS）
│
├── models/
│   ├── auth.go                      # JWT Claims 構造体
│   ├── blogs.go                     # BlogData 構造体
│   ├── blog_users.go                # BlogUsersData 構造体
│   ├── blog_likes.go                # BlogLikesData 構造体
│   └── blog_comments.go            # BlogCommentsData 構造体
│
├── routes/
│   └── routes.go                    # ルーティング定義、DI組み立て
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
│   │   ├── blog_users_GetUserByEmailAndPassword_test.go
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
│   │   ├── auth_mock.go            # モック
│   │   └── auth_test.go            # Service テスト
│   ├── blog_users/
│   │   ├── blog_users.go           # UserService インターフェース
│   │   ├── blog_users_impl.go      # FetchUser, UpdateUser 等（不正な入力でテスト）
│   │   ├── blog_users_mock.go      # モック
│   │   ├── blog_users_FetchUserById_test.go
│   │   ├── blog_users_FetchUserByEmailAndPassword_test.go
│   │   └── blog_users_UpdateUser_test.go
│   ├── blogs/
│   │   ├── blogs.go                # BlogService インターフェース
│   │   ├── blogs_impl.go           # ビジネスロジック実装
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
│   │   ├── blog_likes_mock.go      # モック
│   │   ├── blog_likes_FetchBlogLikesByVisitId_test.go
│   │   ├── blog_likes_IsBlogLiked_test.go
│   │   ├── blog_likes_CreateBlogLike_test.go
│   │   └── blog_likes_DeleteBlogLike_test.go
│   └── blog_comments/
│       ├── blog_comments.go        # CommentService インターフェース
│       ├── blog_comments_impl.go   # コメントロジック実装
│       ├── blog_comments_mock.go   # モック
│       ├── blog_comments_FetchCommentsByBlogId_test.go
│       └── blog_comments_CreateComment_test.go
│
├── repositories/                    # Repository層（データアクセス）
│   ├── auth/
│   │   ├── auth.go                 # （未使用の可能性あり）
│   │   └── auth_impl.go
│   ├── blog_users/
│   │   ├── blog_users.go           # BlogUsersRepository インターフェース
│   │   ├── blog_users_impl.go      # SQL実装
│   │   ├── blog_users_mock.go      # モック
│   │   ├── blog_users_testing_setup.go
│   │   ├── blog_users_FetchUserById_test.go
│   │   ├── blog_users_FetchUserByEmailAndPassword_test.go
│   │   └── blog_users_UpdateBlogUsers_test.go
│   ├── blogs/
│   │   ├── blogs.go                # BlogRepository インターフェース
│   │   ├── blogs_impl.go           # SQL実装（JOIN、集計クエリ含む）
│   │   ├── blogs_mock.go           # モック
│   │   ├── blogs_testing_setup.go  # テストセットアップ
│   │   └── test/                   # Repository テスト群
│   │       ├── main_test.go
│   │       ├── blogs_FetchBlogs_test.go
│   │       ├── blogs_FetchBlogById_test.go
│   │       ├── blogs_FetchBlogsByUserId_test.go
│   │       ├── blogs_FetchBlogCategories_test.go
│   │       ├── blogs_FetchBlogTags_test.go
│   │       ├── blogs_FetchBlogPopular_test.go
│   │       ├── blogs_CreateBlog_test.go
│   │       ├── blogs_UpdateBlog_test.go
│   │       ├── blogs_DeleteBlog_test.go
│   │       └── blogs_pipeline_test.go
│   ├── blog_likes/
│   │   ├── blog_likes.go           # BlogLikeRepository インターフェース
│   │   ├── blog_likes_impl.go      # SQL実装
│   │   ├── blog_likes_mock.go      # モック
│   │   ├── blog_testing_setup.go
│   │   └── blog_likes_pipeline_test.go
│   └── blog_comments/
│       ├── blog_comments.go        # CommentRepository インターフェース
│       ├── blog_comments_impl.go   # SQL実装
│       ├── blog_comments_mock.go   # モック
│       ├── blog_comments_testing_setup.go
│       └── blog_comments_FetchCommentsByBlogId_test.go
│
├── supabase/                        # Supabase接続管理
│   └── client.go                   # 接続プール初期化・テストクエリ・クローズ
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
│
├── terraform/                       # インフラ構成管理（Terraform）
│   ├── main.tf                     # Provider設定
│   ├── variables.tf                # 変数定義
│   ├── cloud-run.tf                # Cloud Runサービス定義
│   ├── cloud-run-iam.tf            # IAMロール設定
│   ├── artifact-registry.tf        # Artifact Registry定義
│   └── secrets.tf                  # Secret Manager定義
│
├── .github/
│   └── workflows/
│       └── deploy.yml              # CI/CDパイプライン
│
├── docs/                            # ドキュメント
├── manuals/                         # マニュアル
├── backup/                          # バックアップ
└── work/                            # 作業用
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
| Go | 1.20 以上（`go.mod` は `go 1.20`） | アプリケーションのビルド・実行 |
| Supabase / PostgreSQL | - | 稼働中の接続先（無料枠で可）。起動時に接続テスト `SELECT 1` を実行する |
| Docker | 任意 | コンテナで起動する場合のみ |

> 注意: `Dockerfile` のビルドステージは現状 `golang:1.19` で、`go.mod` の `go 1.20` と不一致（[`11-tasks.md`](11-tasks.md) の改善タスク参照）。ローカルでは 1.20 以上を使用すること。

### 10.2 セットアップ手順

```bash
git clone https://github.com/kojikawazu/nextjs-echo-back-blog-app.git
cd nextjs-echo-back-blog-app

# 環境変数テンプレートをコピーして値を設定
cp .env.example .env

# 依存解決して起動
go mod download
go run main.go
```

起動フローは [§7.3 サーバー起動フロー](#73-サーバー起動フロー) を参照。`http://localhost:8080/` が `Service is running` を返せば成功。

### 10.3 環境変数

アプリ起動時に `godotenv` がリポジトリ直下の `.env` を読み込む（`main.go` の `firstSetup`）。テンプレートは [`.env.example`](../.env.example)。

| 変数名 | 必須 | 説明 | 参照コード |
|--------|------|------|-----------|
| `JWT_SECRET_KEY` | ✅ | JWT署名鍵。未設定だと `config.init()` の `log.Fatal` で即終了 | `config/config.go:9,12-15` |
| `SUPABASE_URL` | ✅ | PostgreSQL接続URL。`?sslmode=require` はコードが自動付与 | `supabase/client.go:28` |
| `ALLOWED_ORIGINS` | ✅(ブラウザ利用時) | CORS許可オリジン（カンマ区切り） | `middlewares/middlewares.go:17` |
| `ENV` | - | `production` で Cookie が Secure/SameSite=None。ローカルは空 | `config/config.go:10` |
| `PORT` | - | リッスンポート（既定 8080） | `main.go:74-77` |
| `TEST_MODE` | - | `true` でログ破棄（主にテスト用） | `logger/logger.go` |

### 10.4 テスト実行

```bash
# 単体テスト（DB 不要・モック使用）。CI と同じ範囲
go test ./handlers/... ./services/...

# 全テスト（Repository 層は実 DB へ接続）
cp .env.test.example .env.test   # SUPABASE_URL / JWT_SECRET_KEY を設定
go test ./...
```

- Repository 層テストは `repositories/**/.._testing_setup.go` の `godotenv.Load("../../../.env.test")` で `.env.test` を読み、実 DB に接続する。
- Service / Handler 層テストはモックを使うため DB 不要だが、`config` パッケージの `init` が `JWT_SECRET_KEY` を要求する。
- CI（`.github/workflows/deploy.yml`）は DB 不要の `handlers` / `services` のみ実行する。詳細は [`08-test-specification.md`](08-test-specification.md)。

### 10.5 Docker での起動

```bash
docker build -t echo-blog-back .
docker run --env-file .env -p 8080:8080 echo-blog-back
```

イメージ構成（マルチステージ / Distroless）は [§7.1 Dockerイメージ構成](#71-dockerイメージ構成) を参照。

### 10.6 よくあるつまずき

| 症状 | 原因 | 対処 |
|------|------|------|
| 起動直後に `JWT_SECRET_KEY is not set` で落ちる | `.env` 未作成 or `JWT_SECRET_KEY` 空 | `.env` に値を設定する |
| `unable to connect to Supabase` で `log.Fatal` | `SUPABASE_URL` 不正 / DB 未到達 | 接続URL・ネットワーク・DB稼働を確認 |
| ブラウザから CORS エラー | `ALLOWED_ORIGINS` にフロントの URL が無い | フロントのオリジンをカンマ区切りで追加 |
| Cookie が保存されない（本番） | `ENV=production` 時は HTTPS 必須（Secure=true） | HTTPS 経由でアクセスする |
