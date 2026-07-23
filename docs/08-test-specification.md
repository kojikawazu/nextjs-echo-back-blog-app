# テスト仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo）のテスト戦略、テスト構成、およびテスト実行方法を定義するテスト仕様書です。コードベースのリバースエンジニアリングに基づいて作成しています。

## 目次

- [1. テスト戦略概要](#1-テスト戦略概要)
  - [テストピラミッド](#テストピラミッド)
- [2. テストツール・ライブラリ](#2-テストツールライブラリ)
- [3. テストファイル構成](#3-テストファイル構成)
  - [3.1 テストファイル一覧（全63ファイル）](#31-テストファイル一覧全63ファイル)
- [4. テストインフラストラクチャ](#4-テストインフラストラクチャ)
  - [4.1 TestMain パターン](#41-testmain-パターン)
  - [4.2 セットアップ関数](#42-セットアップ関数)
  - [4.3 ログ制御](#43-ログ制御)
- [5. モックパターン](#5-モックパターン)
  - [5.1 Repository モック](#51-repository-モック)
  - [5.2 Service モック](#52-service-モック)
  - [5.3 Cookie モック](#53-cookie-モック)
- [6. テストカテゴリ別詳細](#6-テストカテゴリ別詳細)
  - [6.1 Repository テスト（結合テスト）](#61-repository-テスト結合テスト)
  - [6.2 Service テスト（単体テスト）](#62-service-テスト単体テスト)
  - [6.3 Handler テスト（単体テスト）](#63-handler-テスト単体テスト)
- [7. テスト実行方法](#7-テスト実行方法)
  - [全テスト実行](#全テスト実行)
  - [特定パッケージのテスト実行](#特定パッケージのテスト実行)
  - [特定テスト関数の実行](#特定テスト関数の実行)
  - [テストカバレッジの計測](#テストカバレッジの計測)
- [8. カバレッジ目標](#8-カバレッジ目標)
  - [ドメイン別テストカバレッジ](#ドメイン別テストカバレッジ)
- [9. テスト環境構成](#9-テスト環境構成)
  - [環境変数ファイル](#環境変数ファイル)
  - [注意事項](#注意事項)

---

## 1. テスト戦略概要

本プロジェクトでは、クリーンアーキテクチャの各層に対応した3レベルのテスト戦略を採用している。

| テストレベル | 対象層 | テスト種別 | DB接続 | モック使用 |
|-------------|--------|-----------|--------|-----------|
| Repository テスト | Repository層 | インテグレーション（IT） | 必要（既定は testcontainers の使い捨て PostgreSQL。`SUPABASE_URL` 指定時は実DB） | なし |
| Service テスト | Service層 | 単体テスト（UT） | 不要 | Repository モック |
| Handler テスト | Handler層 | 単体テスト（UT） | 不要 | Service モック + Cookie モック |

> IT は `//go:build integration` タグで分離しており、実行は `go test -tags=integration ./repositories/...`。通常の `go test ./...`（UT）には含まれない。

### テストピラミッド

```
        /  Handler テスト  \        ← HTTPリクエスト/レスポンスの検証
       / Service テスト      \      ← ビジネスロジックの検証
      / Repository テスト      \    ← データアクセスの検証（testcontainers / 実DB）
     /_________________________\
```

---

## 2. テストツール・ライブラリ

| ツール | バージョン | 用途 |
|--------|----------|------|
| testing（標準パッケージ） | Go 1.22 | テストフレームワーク |
| github.com/stretchr/testify | v1.9.0 | アサーション（assert）・モック（mock） |
| net/http/httptest（標準パッケージ） | Go 1.22 | HTTPリクエスト/レスポンスのテスト |
| github.com/testcontainers/testcontainers-go | v0.35.0 | IT用の使い捨て PostgreSQL コンテナ（postgres モジュール） |
| github.com/jackc/pgconn | v1.14.3 | IT のエラー検証（`PgError.Code` で SQLSTATE を判定） |
| github.com/labstack/echo/v4 | v4.12.0 | Echoコンテキストの生成 |

---

## 3. テストファイル構成

### 3.1 テストファイル一覧（主要ファイル）

#### Repository層テスト（結合テスト）

| ファイルパス | テスト対象 |
|-------------|-----------|
| `repositories/blogs/test/main_test.go` | TestMain（`testsupport.Start` 委譲） |
| `repositories/blogs/test/blogs_FetchBlogs_test.go` | 全ブログ取得 |
| `repositories/blogs/test/blogs_FetchBlogById_test.go` | ブログID指定取得 + 準正常系（不正UUID） |
| `repositories/blogs/test/blogs_FetchBlogsByUserId_test.go` | ユーザーID指定取得 + 準正常系（不正UUID） |
| `repositories/blogs/test/blogs_FetchBlogCategories_test.go` | カテゴリ取得 |
| `repositories/blogs/test/blogs_FetchBlogTags_test.go` | タグ取得 |
| `repositories/blogs/test/blogs_FetchBlogPopular_test.go` | 人気ブログ取得 + 境界（0件） |
| `repositories/blogs/test/blogs_CreateBlog_test.go` | ブログ作成（異常系: 空userId） |
| `repositories/blogs/test/blogs_UpdateBlog_test.go` | ブログ更新（異常系: 空id） |
| `repositories/blogs/test/blogs_DeleteBlog_test.go` | ブログ削除（異常系: 空id） |
| `repositories/blogs/test/blogs_pipeline_test.go` | ブログCRUDパイプライン |
| `repositories/blogs/test/zz_connection_closed_test.go` | 異常系（接続断で安全に失敗） |
| `repositories/blog_users/main_test.go` | TestMain（`testsupport.Start` 委譲） |
| `repositories/blog_users/blog_users_FetchUserById_test.go` | ユーザーID指定取得 + 準正常系（空id） |
| `repositories/blog_users/blog_users_FetchUserByEmailAndPassword_test.go` | メール・パスワード指定取得 + 準正常系 |
| `repositories/blog_users/blog_users_UpdateBlogUsers_test.go` | ユーザー更新 + 準正常系（存在しないID） |
| `repositories/blog_users/zz_connection_closed_test.go` | 異常系（接続断で安全に失敗） |
| `repositories/blog_comments/main_test.go` | TestMain（`testsupport.Start` 委譲） |
| `repositories/blog_comments/blog_comments_FetchCommentsByBlogId_test.go` | ブログID指定コメント取得 + 準正常系（不正UUID） |
| `repositories/blog_comments/blog_comments_CreateComment_test.go` | コメント作成 + 準正常系（不正UUID） |
| `repositories/blog_comments/zz_connection_closed_test.go` | 異常系（接続断で安全に失敗） |
| `repositories/blog_likes/main_test.go` | TestMain（`testsupport.Start` 委譲） |
| `repositories/blog_likes/blog_likes_pipeline_test.go` | いいねCRUDパイプライン |
| `repositories/blog_likes/zz_connection_closed_test.go` | 異常系（接続断で安全に失敗） |
| `testsupport/testdb.go` | IT共有ヘルパー（testcontainers 起動 / 実DB切替、TestMain 委譲先） |

#### Service層テスト（単体テスト）

| ファイルパス | テスト対象 |
|-------------|-----------|
| `services/blogs/test/main_test.go` | TestMain（セットアップ） |
| `services/blogs/test/blogs_FetchBlogs_test.go` | 全ブログ取得 |
| `services/blogs/test/blogs_FetchBlogById_test.go` | ブログID指定取得 |
| `services/blogs/test/blogs_FetchBlogsByUserId_test.go` | ユーザーID指定取得 |
| `services/blogs/test/blogs_FetchBlogCategories_test.go` | カテゴリ取得 |
| `services/blogs/test/blogs_FetchBlogTags_test.go` | タグ取得 |
| `services/blogs/test/blogs_FetchBlogPopular_test.go` | 人気ブログ取得 |
| `services/blogs/test/blogs_CreateBlog_test.go` | ブログ作成 |
| `services/blogs/test/blogs_UpdateBlog_test.go` | ブログ更新 |
| `services/blogs/test/blogs_DeleteBlog_test.go` | ブログ削除 |
| `services/blog_users/blog_users_FetchUserById_test.go` | ユーザーID指定取得 |
| `services/blog_users/blog_users_FetchUserByEmailAndPassword_test.go` | メール・パスワード指定取得 |
| `services/blog_users/blog_users_UpdateUser_test.go` | ユーザー更新 |
| `services/blog_comments/blog_comments_FetchCommentsByBlogId_test.go` | コメント取得 |
| `services/blog_comments/blog_comments_CreateComment_test.go` | コメント作成 |
| `services/blog_likes/blog_likes_FetchBlogLikesByVisitId_test.go` | 訪問者IDいいね取得 |
| `services/blog_likes/blog_likes_IsBlogLiked_test.go` | いいね存在確認 |
| `services/blog_likes/blog_likes_CreateBlogLike_test.go` | いいね作成 |
| `services/blog_likes/blog_likes_DeleteBlogLike_test.go` | いいね削除 |
| `services/auth/auth_test.go` | ログインバリデーション |

#### Handler層テスト（単体テスト）

| ファイルパス | テスト対象 |
|-------------|-----------|
| `handlers/blogs/test/main_test.go` | TestMain（セットアップ） |
| `handlers/blogs/test/blogs_FetchBlogs_test.go` | 全ブログ取得ハンドラ |
| `handlers/blogs/test/blogs_FetchBlogById_test.go` | ブログID指定取得ハンドラ |
| `handlers/blogs/test/blogs_FetchBlogsByUserId_test.go` | ユーザーID指定取得ハンドラ |
| `handlers/blogs/test/blogs_FetchBlogCategories_test.go` | カテゴリ取得ハンドラ |
| `handlers/blogs/test/blogs_FetchBlogTags_test.go` | タグ取得ハンドラ |
| `handlers/blogs/test/blogs_FetchBlogPopular_test.go` | 人気ブログ取得ハンドラ |
| `handlers/blogs/test/blogs_CreateBlog_test.go` | ブログ作成ハンドラ |
| `handlers/blogs/test/blogs_UpdateBlog_test.go` | ブログ更新ハンドラ |
| `handlers/blogs/test/blogs_DeleteBlog_test.go` | ブログ削除ハンドラ |
| `handlers/blog_users/blog_users_FetchUser_test.go` | ユーザー取得ハンドラ |
| `handlers/blog_users/blog_users_GetUserByEmailAndPassword_test.go` | メール・パスワード取得ハンドラ |
| `handlers/blog_users/blog_users_UpdateUser_test.go` | ユーザー更新ハンドラ |
| `handlers/blog_comments/blog_comments_FetchCommentsByBlogId_test.go` | コメント取得ハンドラ |
| `handlers/blog_comments/blog_comments_CreateComment_test.go` | コメント作成ハンドラ |
| `handlers/blog_likes/blog_likes_FetchBlogLikesByVisitId_test.go` | いいね取得ハンドラ |
| `handlers/blog_likes/blog_likes_IsBlogLiked_test.go` | いいね存在確認ハンドラ |
| `handlers/blog_likes/blog_likes_CreateBlogLike_test.go` | いいね作成ハンドラ |
| `handlers/blog_likes/blog_likes_DeleteBlogLike_test.go` | いいね削除ハンドラ |
| `handlers/blog_likes/blog_likes_GenerateVisitorId_test.go` | 訪問者ID生成ハンドラ |
| `handlers/auth/auth_test.go` | 認証ハンドラ |
| `handlers/auth/auth_Login_test.go` | ログインハンドラ |
| `handlers/auth/auth_CheckAuth_test.go` | 認証確認ハンドラ |
| `handlers/auth/auth_Logout_test.go` | ログアウトハンドラ |
| `handlers/auth/auth_pipeline_test.go` | 認証フローパイプライン |

#### その他

| ファイルパス | テスト対象 |
|-------------|-----------|
| `utils/log/utils_log_test.go` | ログユーティリティ |

---

## 4. テストインフラストラクチャ

### 4.1 TestMain パターン

各テストパッケージでは `TestMain` 関数を使用して、テスト実行前のセットアップを行う。

```go
func TestMain(m *testing.M) {
    // テスト前のセットアップ
    パッケージ名.SetupTest(nil)

    // テスト実行
    code := m.Run()

    // 終了コードを返す
    os.Exit(code)
}
```

### 4.2 セットアップ関数

テストセットアップは、テストレベルに応じて2種類存在する。

#### Repository層セットアップ（IT / testsupport 経由）

各 Repository テストパッケージの `TestMain` は、共有ヘルパー `backend/testsupport` の `Start` に委譲する（`//go:build integration`）。

```go
//go:build integration

func TestMain(m *testing.M) {
    os.Exit(testsupport.Start(m))
}
```

`testsupport.Start` の動作:

- `SUPABASE_URL` が**実環境変数として**設定済み → その DB を使う（実 Supabase 等。`TEST_*` も呼び出し側が指定）。
- 未設定（既定） → `testcontainers` で使い捨ての PostgreSQL を起動し、`testsupport/testdata/schema.sql` / `seed.sql` を適用、`SUPABASE_URL` / `DB_SSLMODE=disable` / `TEST_*`（固定UUID）を設定してから接続する。テスト終了時にプールをクローズしコンテナを破棄する。
- `.env` ファイルの自動読み込みは行わない（古い `.env.test` がコンテナ経路を意図せず乗っ取るのを防ぐため）。
- モジュールルート（`go.mod`）を上方向探索で特定するため、旧来のパッケージ階層依存の相対パス（`../../` と `../../../` の混在）は解消済み。

#### Service/Handler層セットアップ（DB接続なし）

```go
func SetupTest(t *testing.T) {
    // 環境変数の読み込み（.env.test）
    godotenv.Load("../../../.env.test")

    // ログ設定の初期化
    logger.InitLogger()
}
```

- DB接続は不要（モックで代替）
- 環境変数 `JWT_SECRET_KEY` の設定が必要（config パッケージの init で参照）
- `TEST_MODE=true` 設定時はログ出力が抑制される

### 4.3 ログ制御

テスト実行中のログ出力は `TEST_MODE` 環境変数で制御する。

```go
// logger/logger.go
func InitLogger() {
    if os.Getenv("TEST_MODE") == "true" {
        InfoLog.SetOutput(io.Discard)
        ErrorLog.SetOutput(io.Discard)
        WarnLog.SetOutput(io.Discard)
        DebugLog.SetOutput(io.Discard)
    }
}
```

---

## 5. モックパターン

### 5.1 Repository モック

各ドメインの Repository にはモック実装が用意されている。`testify/mock` パッケージを使用する。

| モックファイル | 実装インターフェース |
|---------------|---------------------|
| `repositories/blogs/blogs_mock.go` | `BlogRepository` |
| `repositories/blog_users/blog_users_mock.go` | `BlogUsersRepository` |
| `repositories/blog_likes/blog_likes_mock.go` | `BlogLikeRepository` |
| `repositories/blog_comments/blog_comments_mock.go` | `CommentRepository` |

#### モック構造体の例

```go
type MockBlogRepository struct {
    mock.Mock
}

func (m *MockBlogRepository) FetchBlogs() ([]models.BlogData, error) {
    args := m.Called()
    if args.Get(0) != nil {
        return args.Get(0).([]models.BlogData), args.Error(1)
    }
    return nil, args.Error(1)
}
```

#### モックの使用例（Service テスト）

```go
mockBlogRepository := new(repositories_blogs.MockBlogRepository)
blogService := services_blogs.NewBlogService(mockBlogRepository)

mockBlogRepository.On("FetchBlogs").Return(mockBlogData, nil)

blogs, err := blogService.FetchBlogs()

assert.NoError(t, err)
mockBlogRepository.AssertExpectations(t)
```

### 5.2 Service モック

各ドメインの Service にはモック実装が用意されている。Handler テストで使用する。

| モックファイル | 実装インターフェース |
|---------------|---------------------|
| `services/blogs/blogs_mock.go` | `BlogService` |
| `services/blog_users/blog_users_mock.go` | `UserService` |
| `services/blog_likes/blog_likes_mock.go` | `BlogLikeService` |
| `services/blog_comments/blog_comments_mock.go` | `CommentService` |
| `services/auth/auth_mock.go` | `AuthService` |

### 5.3 Cookie モック

Handler テストでは Cookie ユーティリティもモック化する。

| モックファイル | 実装インターフェース |
|---------------|---------------------|
| `utils/cookie/cookie_mock.go` | `CookieUtils` |

#### Cookie モックヘルパー

認証が必要なハンドラのテストでは、`SetMockBlogCookies` ヘルパー関数を使用してCookieを設定する。

```go
// handlers/blogs/blogs_cookie_mock.go
func SetMockBlogCookies(c echo.Context, req *http.Request, mockCookieUtils *utils_cookie.MockCookieUtils) {
    token := "mocked-token"
    validUserId := "valid-user-id"

    cookie := &http.Cookie{
        Name:  "token",
        Value: token,
        Path:  "/",
    }
    req.AddCookie(cookie)

    mockCookieUtils.On("GetAuthCookieValue", c, "token").Return(token, nil)
    mockCookieUtils.On("GetUserIdFromToken", c, token).Return(validUserId, nil)
}
```

---

## 6. テストカテゴリ別詳細

### 6.1 Repository テスト（インテグレーション / IT）

#### 概要

- 既定では `testcontainers` が使い捨ての PostgreSQL を起動し、決定的シードに対してテストを実行する（実 Supabase には接続しない）。
- `SUPABASE_URL` を環境変数で指定した場合のみ、その実 DB に接続する（`TEST_*` も併せて指定）。
- 3分類（正常系 / 準正常系 / 異常系）を検証する。

#### 環境変数要件

testcontainers 経路では `testsupport` が下記を自動設定するため、実行時に指定するものは無い（Docker 稼働のみ前提）。実 DB 経路では呼び出し側が環境変数で与える。

| 環境変数 | 説明 | 用途 |
|---------|------|------|
| `SUPABASE_URL` | 接続URL。**未設定なら testcontainers を起動** | DB接続の切替 |
| `DB_SSLMODE` | SSLモード（コンテナ経路では `disable`） | SSL設定 |
| `TEST_USER_ID` / `TEST_USER_NAME` / `TEST_USER_EMAIL` / `TEST_USER_PASSWD` | シード済みユーザーに対応する値 | ユーザー関連テスト |
| `TEST_BLOG_ID` | シード済みブログID (UUID) | ブログ取得テスト |

> Repository テストは `config` パッケージを import しないため `JWT_SECRET_KEY` は不要。

#### 分類とテスト対象

- 正常系: データ取得・作成・更新・削除（シードの既知値に対する具体値アサーション）。
- 準正常系: 存在しないID・無効なUUID形式・空文字パラメータ。
- 異常系: コネクションプールがクローズ済みでもクエリが panic せず error を返すこと（`zz_connection_closed_test.go`）。
- パイプラインテスト: CRUD操作の一連の流れ（Create → Read → Update → Delete）。

#### テストの特徴（脆いアサーションの緩和）

エラーメッセージ全文一致ではなく、SQLSTATE コードで検証し、PostgreSQL/ドライバのバージョン差異に強くしている。

```go
// 正常系: シードの既知値を具体的に検証
id := os.Getenv("TEST_BLOG_ID")
blog, err := repo.FetchBlogById(id)
assert.NoError(t, err)
assert.Equal(t, "Test Blog Title", blog.Title)

// 準正常系: 無効なUUIDは SQLSTATE 22P02 を検証（メッセージ全文には依存しない）
blog, err = repo.FetchBlogById("2")
assert.Error(t, err)
assert.Nil(t, blog)
var pgErr *pgconn.PgError
if assert.ErrorAs(t, err, &pgErr) {
    assert.Equal(t, "22P02", pgErr.Code)
}
```

### 6.2 Service テスト（単体テスト）

#### 概要

- Repository モックを使用してService層のビジネスロジックを検証する
- バリデーションロジックの正常系・異常系をテストする
- DB接続は不要

#### テスト対象機能

- 入力バリデーション（空文字チェック、メールフォーマットチェック）
- Repository呼び出しの正常系・異常系
- データ変換ロジック（タグのカンマ区切り分割・重複除去・ソート）
- 空リスト返却の動作確認

#### テストパターン

| パターン | 説明 |
|---------|------|
| 正常系 | モックが正常データを返す場合の動作確認 |
| 空リスト | モックが空リストを返す場合の動作確認 |
| バリデーションエラー | 空文字や不正形式の入力に対するエラー確認 |
| Repositoryエラー | モックがエラーを返す場合の動作確認 |

### 6.3 Handler テスト（単体テスト）

#### 概要

- `net/http/httptest` を使用してHTTPリクエスト/レスポンスを検証する
- Service モックとCookie モックを使用する
- ステータスコード、レスポンスボディ、エラーメッセージを検証する

#### テスト構造

```go
func TestHandler_FetchBlogs(t *testing.T) {
    // 1. Echoコンテキストのセットアップ
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/api/blogs", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    // 2. モックの初期化
    mockCookieUtils := new(utils_cookie.MockCookieUtils)
    mockService := new(service_blogs.MockBlogService)
    handler := handlers_blogs.NewBlogHandler(mockService, mockCookieUtils)

    // 3. モックの振る舞い設定
    mockService.On("FetchBlogs").Return(mockBlogData, nil)

    // 4. ハンドラ実行
    err := handler.FetchBlogs(c)
    assert.NoError(t, err)

    // 5. レスポンス検証
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Contains(t, rec.Body.String(), "title2")

    // 6. モック期待値の検証
    mockService.AssertExpectations(t)
}
```

#### テストパターン

| パターン | 検証内容 |
|---------|---------|
| 正常系 | ステータス200、レスポンスボディにデータ含有 |
| エラー系 | ステータス400/401/404/500、エラーメッセージの確認 |
| 空リスト | ステータス200、空配列 `[]` の確認 |
| 認証エラー | Cookie未設定時の401レスポンス |
| バリデーションエラー | 不正パラメータ時の400レスポンス |

---

## 7. テスト実行方法

> Go アプリは `backend/` 配下に集約されている（モノレポ構成）。以下の `go` コマンドはすべて `backend/` ディレクトリ内で実行する（例: `cd backend`）。

### 全テスト実行

```bash
# UT（DB不要）
go test ./... -v

# IT（Repository層 / testcontainers。Docker 稼働が前提）
go test -tags=integration ./repositories/... -v
```

> `go test ./...`（タグなし）に IT は含まれない（`//go:build integration` で分離）。

### 特定パッケージのテスト実行

```bash
# Repository テスト（IT / testcontainers）
go test -tags=integration ./repositories/blogs/test/... -v

# Service テスト（UT）
go test ./services/blogs/test/... -v

# Handler テスト（UT）
go test ./handlers/blogs/test/... -v
```

### 特定テスト関数の実行

```bash
go test ./handlers/blogs/test/... -v -run TestHandler_FetchBlogs
```

### テストカバレッジの計測

```bash
# UT
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 8. カバレッジ目標

| テストレベル | 対象 | 目標カバレッジ |
|-------------|------|---------------|
| Repository テスト | 全CRUD操作 | 各メソッドの正常系・異常系 |
| Service テスト | バリデーション、ビジネスロジック | 各メソッドの全分岐 |
| Handler テスト | HTTPレスポンス検証 | 各エンドポイントの全ステータスコード |

### ドメイン別テストカバレッジ

| ドメイン | Repository | Service | Handler |
|---------|-----------|---------|---------|
| Blogs | 11ファイル | 10ファイル | 10ファイル |
| Blog Users | 3ファイル | 3ファイル | 3ファイル |
| Blog Likes | 1ファイル | 4ファイル | 5ファイル |
| Blog Comments | 2ファイル | 2ファイル | 2ファイル |
| Auth | - | 1ファイル | 5ファイル |

---

## 9. テスト環境構成

### 環境変数ファイル

テスト用の環境変数は `.env.test` ファイルに定義する。このファイルはリポジトリのルートディレクトリに配置する。

#### 必要な環境変数

```env
# DB接続（Repository テストで必要）
SUPABASE_URL=postgresql://user:password@host:port/dbname

# JWT設定（全テストで必要）
JWT_SECRET_KEY=your-secret-key

# テストモード
TEST_MODE=true
ENV=test

# テストデータID（Repository テストで必要）
TEST_USER_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
TEST_BLOG_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

# CORS（config初期化で参照される可能性あり）
ALLOWED_ORIGINS=http://localhost:3000
```

### 注意事項

- Repository テストは実際のデータベースに対して実行されるため、テスト用のSupabaseプロジェクトまたは専用のテスト環境が必要
- テストデータの事前作成が必要（`TEST_USER_ID`、`TEST_BLOG_ID` に対応するレコード）
- パイプラインテスト（Create→Read→Update→Delete）はテストデータの作成と削除を一連の流れで行う
- `JWT_SECRET_KEY` が未設定の場合、`config` パッケージの `init()` で `log.Fatal` が発生しテストが中断される
