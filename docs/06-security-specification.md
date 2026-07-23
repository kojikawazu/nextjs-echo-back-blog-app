# セキュリティ仕様書

本ドキュメントは、ブログWebアプリケーションバックエンドにおけるセキュリティ関連の実装仕様を、ソースコードから逆引きして整理したものである。

## 目次

- [1. 認証（Authentication）](#1-認証authentication)
  - [1.1 JWT認証方式](#11-jwt認証方式)
  - [1.2 認証トークンのペイロード](#12-認証トークンのペイロード)
  - [1.3 認証フロー](#13-認証フロー)
  - [1.4 JWT鍵の管理](#14-jwt鍵の管理)
- [2. 認可（Authorization）](#2-認可authorization)
  - [2.1 Cookie認証チェックによる書き込み操作の保護](#21-cookie認証チェックによる書き込み操作の保護)
  - [2.2 認証不要のエンドポイント](#22-認証不要のエンドポイント)
  - [2.3 認可の制限事項](#23-認可の制限事項)
- [3. Cookie セキュリティ](#3-cookie-セキュリティ)
  - [3.1 認証Cookie（`token`）](#31-認証cookietoken)
  - [3.2 訪問者ID Cookie（`visit-id-token`）](#32-訪問者id-cookievisit-id-token)
  - [3.3 環境判定](#33-環境判定)
  - [3.4 Cookie削除方式](#34-cookie削除方式)
- [4. CORS（Cross-Origin Resource Sharing）](#4-corscross-origin-resource-sharing)
  - [4.1 CORS設定](#41-cors設定)
  - [4.2 CORS設定の留意事項](#42-cors設定の留意事項)
- [5. 訪問者追跡（Visitor Tracking）](#5-訪問者追跡visitor-tracking)
  - [5.1 概要](#51-概要)
  - [5.2 訪問者ID生成フロー](#52-訪問者id生成フロー)
  - [5.3 訪問者IDの使用](#53-訪問者idの使用)
- [6. インフラストラクチャセキュリティ](#6-インフラストラクチャセキュリティ)
  - [6.1 Secret Manager](#61-secret-manager)
  - [6.2 Distrolessベースイメージ](#62-distrolessベースイメージ)
  - [6.3 IAMロール](#63-iamロール)
  - [6.4 Artifact Registry](#64-artifact-registry)
  - [6.5 SSL/TLS](#65-ssltls)
- [7. SQLインジェクション対策](#7-sqlインジェクション対策)
- [8. 既知のセキュリティ課題](#8-既知のセキュリティ課題)
  - [8.1 パスワードの平文保存（重大）](#81-パスワードの平文保存重大)
  - [8.2 レート制限の欠如](#82-レート制限の欠如)
  - [8.3 CSRF保護の欠如](#83-csrf保護の欠如)
  - [8.4 認可の不足（所有権チェック）](#84-認可の不足所有権チェック)
  - [8.5 入力長の制限なし](#85-入力長の制限なし)
  - [8.6 エラーメッセージの文字列ベース判定](#86-エラーメッセージの文字列ベース判定)

---

## 1. 認証（Authentication）

### 1.1 JWT認証方式

本アプリケーションはJWT（JSON Web Token）を使用してユーザー認証を行う。

| 項目 | 値 |
|---|---|
| 署名アルゴリズム | HS256（HMAC-SHA256） |
| 署名鍵 | 環境変数 `JWT_SECRET_KEY` から取得 |
| トークン有効期限 | 1時間（`time.Now().Add(1 * time.Hour)`） |
| トークン格納先 | HTTP-only Cookie（Cookie名: `token`） |
| JWTライブラリ | `github.com/golang-jwt/jwt v3.2.2+incompatible` |

### 1.2 認証トークンのペイロード

```go
type Claims struct {
    UserID   string `json:"user_id"`   // ユーザーUUID
    Email    string `json:"email"`     // メールアドレス
    Username string `json:"username"`  // ユーザー名
    jwt.StandardClaims                 // ExpiresAt（有効期限）
}
```

### 1.3 認証フロー

#### ログイン（`POST /api/users/login`）

1. リクエストボディから `email` と `password` を取得する。
2. Service層でバリデーションを実施する（空チェック、メールフォーマット検証）。
3. データベースに対して `email` と `password` の平文一致でユーザーを検索する。
4. JWTトークンを生成する（Claims にユーザーID、メールアドレス、ユーザー名を含む）。
5. HTTP-only Cookieにトークンをセットする。
6. `{"message": "Login successful"}` を返却する。

#### 認証確認（`GET /api/users/auth-check`）

1. `token` Cookieからトークン文字列を取得する。
2. `jwt.ParseWithClaims` でトークンを検証する（署名・有効期限）。
3. 検証成功時に `user_id`, `username`, `email` を返却する。

#### ログアウト（`POST /api/users/logout`）

1. `token` Cookieの値を空文字列に設定し、有効期限を過去（`time.Unix(0, 0)`）に設定する。
2. `{"message": "Logout successful"}` を返却する。

### 1.4 JWT鍵の管理

- `JWT_SECRET_KEY` 環境変数が未設定の場合、アプリケーション起動時に `log.Fatal` で即座に終了する。
- 本番環境ではGCP Secret Managerに格納され、Cloud Runのコンテナ環境変数として注入される。

## 2. 認可（Authorization）

### 2.1 Cookie認証チェックによる書き込み操作の保護

以下のエンドポイントでは、リクエスト処理前にCookieからJWTトークンを取得・検証し、認証済みユーザーであることを確認する。

| エンドポイント | HTTPメソッド | 認証要否 |
|---|---|---|
| `/api/blogs/create` | POST | 必要 |
| `/api/blogs/update/:id` | PUT | 必要 |
| `/api/blogs/delete/:id` | DELETE | 必要 |
| `/api/users/detail` | GET | 必要 |
| `/api/users/update` | PUT | 必要 |

認証チェックの処理フロー:

1. `CookieUtils.GetAuthCookieValue(c, "token")` でCookieの値を取得する。
2. `CookieUtils.GetUserIdFromToken(c, cookieValue)` でトークンを検証し、ユーザーIDを取得する。
3. 検証失敗時は `401 Unauthorized` を返却する。

### 2.2 認証不要のエンドポイント

以下のエンドポイントは認証なしでアクセス可能。

| エンドポイント | HTTPメソッド | 説明 |
|---|---|---|
| `/` | GET | ヘルスチェック |
| `/api/blogs` | GET | ブログ一覧取得 |
| `/api/blogs/user/:userId` | GET | ユーザー別ブログ一覧 |
| `/api/blogs/detail/:id` | GET | ブログ詳細取得 |
| `/api/blogs/categories` | GET | カテゴリ一覧取得 |
| `/api/blogs/tags` | GET | タグ一覧取得 |
| `/api/blogs/popular/:count` | GET | 人気ブログ取得 |
| `/api/users/login` | POST | ログイン |
| `/api/users/auth-check` | GET | 認証確認 |
| `/api/users/logout` | POST | ログアウト |
| `/api/blog-likes/*` | GET/POST/DELETE | いいね関連（visit-id認証） |
| `/api/comments/*` | GET/POST | コメント関連 |

### 2.3 認可の制限事項

- ブログの更新・削除時に、操作ユーザーがブログの作成者であるかの確認（所有権チェック）は行われていない。認証済みであれば他ユーザーのブログも操作可能。
- ミドルウェアレベルの認証ガードは存在せず、各ハンドラが個別にCookie認証を実装している。

## 3. Cookie セキュリティ

### 3.1 認証Cookie（`token`）

| 属性 | 本番環境（production） | 開発環境 |
|---|---|---|
| Name | `token` | `token` |
| HttpOnly | `true` | `true` |
| Secure | `true` | `false` |
| SameSite | `None` | `Lax` |
| Path | `/` | `/` |
| Expires | 現在時刻 + 1時間 | 現在時刻 + 1時間 |

### 3.2 訪問者ID Cookie（`visit-id-token`）

| 属性 | 本番環境（production） | 開発環境 |
|---|---|---|
| Name | `visit-id-token` | `visit-id-token` |
| HttpOnly | `true` | `true` |
| Secure | `true` | `false` |
| SameSite | `None` | `Lax` |
| Path | `/` | `/` |
| Expires | 現在時刻 + 1時間 | 現在時刻 + 1時間 |

### 3.3 環境判定

環境変数 `ENV` の値が `"production"` の場合に本番環境と判定される。

```go
var IsProduction = os.Getenv("ENV") == "production"
```

### 3.4 Cookie削除方式

ログアウト時のCookie削除は、有効期限を `time.Unix(0, 0)`（1970年1月1日）に設定することで実現する。

## 4. CORS（Cross-Origin Resource Sharing）

### 4.1 CORS設定

| 設定項目 | 値 |
|---|---|
| AllowOrigins | 環境変数 `ALLOWED_ORIGINS` をカンマ区切りで分割 |
| AllowMethods | GET, POST, PUT, DELETE |
| AllowHeaders | Origin, Content-Type, Authorization, Access-Control-Allow-Credentials |
| AllowCredentials | true |

### 4.2 CORS設定の留意事項

- `AllowCredentials: true` により、クライアント側でも `withCredentials: true` の設定が必要。
- `AllowOrigins` にワイルドカード（`*`）を指定すると `AllowCredentials: true` と競合するため、具体的なオリジンを指定する必要がある。
- `ALLOWED_ORIGINS` はカンマ区切りで複数オリジンを指定可能（例: `https://example.com,https://staging.example.com`）。

## 5. 訪問者追跡（Visitor Tracking）

### 5.1 概要

未認証の訪問者がブログにいいねを付けるための仕組みとして、訪問者ID（visit-id）を使用する。

### 5.2 訪問者ID生成フロー

1. `GET /api/blog-likes/generate-visit-id` にアクセスする。
2. 既存の `visit-id-token` Cookieが有効か確認する。
3. 有効なCookieが存在しない場合、以下を実行する:
   - `uuid.New().String()` でUUIDを生成する。
   - UUIDをClaimsに含めたJWTトークンを作成する（HS256、1時間有効）。
   - `visit-id-token` CookieにJWTトークンを格納する。
4. 有効なCookieが存在する場合、生成をスキップする。

### 5.3 訪問者IDの使用

- いいね操作（作成・削除・確認・一覧取得）時に `visit-id-token` Cookieからvisit-idを抽出する。
- データベースの `blog_likes.visit_id` カラムにUUID文字列として保存される。
- `blog_id` と `visit_id` の組み合わせで重複いいねを防止する。

## 6. インフラストラクチャセキュリティ

### 6.1 Secret Manager

機密情報はGCP Secret Managerで管理される。以下のシークレットが登録されている。

| シークレットキー | 説明 | Terraform変数 |
|---|---|---|
| allowed_origins | CORSの許可オリジン | `var.allowed_origins` |
| env_word | 環境識別子 | `var.env_word` |
| jwt_secret_key | JWT署名鍵 | `var.jwt_secret_key` |
| supabase_url | Supabase接続URL | `var.supabase_url` |
| tmp_allowed_origins | 一時的な許可オリジン | `var.tmp_allowed_origins` |
| gcp_project_id | GCPプロジェクトID | `var.gcp_project_id` |

- Cloud Runのサービスアカウントに `roles/secretmanager.secretAccessor` ロールが付与されている。
- Terraform変数の `sensitive = true` により、計画・適用時のログに値が出力されない。

### 6.2 Distrolessベースイメージ

- 実行ステージのDockerイメージとして `gcr.io/distroless/base` を使用。
- シェル、パッケージマネージャー、その他の不要なプログラムが含まれないため、攻撃対象面が最小化される。
- マルチステージビルド（ビルド: `golang:1.19`、実行: `distroless`）により、ビルドツールが実行環境に残らない。

### 6.3 IAMロール

- Cloud Run用のサービスアカウント（`cloud-run-sa`）が作成される。
- Terraformでは `google_cloud_run_service_iam_member` でInvokerロールが定義されているが、GitHub Actionsのデプロイでは `--allow-unauthenticated` フラグを使用しており、Cloud Runサービスへのパブリックアクセスが許可されている。
- Secret ManagerへのアクセスはサービスアカウントベースのIAMで制限される。

### 6.4 Artifact Registry

- DockerイメージはGCP Artifact Registryに格納される。
- Cloud RunはArtifact Registry上のコミットSHAタグ付きイメージ（`:${{ github.sha }}`）からデプロイされる。`:latest` タグは使用していない。

### 6.5 SSL/TLS

- Supabase接続の SSL モードは環境変数 `DB_SSLMODE` で制御し、**未設定時は `require`**（本番デフォルト）。本番では `require` を維持する。
- `disable` はローカルの非SSL DB（IT の testcontainers など）に接続する場合のみ使用する。接続文字列は `backend/supabase/client.go` の `buildConnString` が `SUPABASE_URL` に `sslmode` を付与して組み立てる。
- Cloud RunはHTTPSを標準で提供するため、クライアント-サーバー間の通信は暗号化される。

## 7. SQLインジェクション対策

- すべてのSQLクエリでパラメータ化クエリ（`$1`, `$2`, ...）を使用しており、SQLインジェクションを防止している。
- ユーザー入力値が直接SQLに埋め込まれる箇所は存在しない。

## 8. 既知のセキュリティ課題

以下は現在の実装における既知のセキュリティ課題である。

### 8.1 パスワードの平文保存（重大）

- `blog_users` テーブルの `password` カラムにパスワードが平文で格納されている。
- ログイン認証時のSQLクエリ（`WHERE email = $1 AND password = $2`）で平文比較を行っている。
- **推奨対策**: bcryptなどのパスワードハッシュアルゴリズムを使用してハッシュ化すべきである。

### 8.2 レート制限の欠如

- ログインエンドポイントを含む全APIエンドポイントにレート制限が設定されていない。
- ブルートフォース攻撃やDDoS攻撃に対して脆弱である。
- **推奨対策**: Echo ミドルウェアまたはCloud Runレベルでレート制限を実装すべきである。

### 8.3 CSRF保護の欠如

- クロスサイトリクエストフォージェリ（CSRF）トークンの検証が実装されていない。
- SameSite=None（本番環境）の設定により、クロスオリジンからのCookieが自動送信される。
- `AllowCredentials: true` のCORS設定と組み合わせることで、CSRF攻撃のリスクが存在する。
- **推奨対策**: CSRFトークンの実装、またはSameSite=Strictの検討を行うべきである。

### 8.4 認可の不足（所有権チェック）

- ブログの更新・削除時に、リクエストユーザーがブログの作成者であるかの確認が行われていない。
- 認証済みユーザーであれば他のユーザーが作成したブログを変更・削除可能。
- **推奨対策**: `blogs.blog_user_id` とトークンのUserIDを比較する所有権チェックを追加すべきである。

### 8.5 入力長の制限なし

- リクエストボディの各フィールドに対する最大長の制限が実装されていない。
- 極端に大きなデータの送信によるリソース消費が可能。
- **推奨対策**: バリデーション層で各フィールドの最大長を定義すべきである。

### 8.6 エラーメッセージの文字列ベース判定

- Handler層のエラーハンドリングがエラーメッセージの文字列マッチング（`err.Error()` のswitch-case）に依存しており、エラーメッセージの変更がAPIレスポンスに意図しない影響を与える可能性がある。
- **推奨対策**: カスタムエラー型またはエラーコードの導入を検討すべきである。
