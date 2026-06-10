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
| Repository テスト | Repository層 | 結合テスト | 必要（Supabase） | なし |
| Service テスト | Service層 | 単体テスト | 不要 | Repository モック |
| Handler テスト | Handler層 | 単体テスト | 不要 | Service モック + Cookie モック |

### テストピラミッド

```
        /  Handler テスト  \        ← HTTPリクエスト/レスポンスの検証
       / Service テスト      \      ← ビジネスロジックの検証
      / Repository テスト      \    ← データアクセスの検証（実DB）
     /_________________________\
```

---

## 2. テストツール・ライブラリ

| ツール | バージョン | 用途 |
|--------|----------|------|
| testing（標準パッケージ） | Go 1.20 | テストフレームワーク |
| github.com/stretchr/testify | v1.9.0 | アサーション（assert）・モック（mock） |
| net/http/httptest（標準パッケージ） | Go 1.20 | HTTPリクエスト/レスポンスのテスト |
| github.com/joho/godotenv | v1.5.1 | テスト用環境変数の読み込み |
| github.com/labstack/echo/v4 | v4.12.0 | Echoコンテキストの生成 |

---

## 3. テストファイル構成

### 3.1 テストファイル一覧（全63ファイル）

#### Repository層テスト（結合テスト）

| ファイルパス | テスト対象 |
|-------------|-----------|
| `repositories/blogs/test/main_test.go` | TestMain（セットアップ） |
| `repositories/blogs/test/blogs_FetchBlogs_test.go` | 全ブログ取得 |
| `repositories/blogs/test/blogs_FetchBlogById_test.go` | ブログID指定取得 |
| `repositories/blogs/test/blogs_FetchBlogsByUserId_test.go` | ユーザーID指定取得 |
| `repositories/blogs/test/blogs_FetchBlogCategories_test.go` | カテゴリ取得 |
| `repositories/blogs/test/blogs_FetchBlogTags_test.go` | タグ取得 |
| `repositories/blogs/test/blogs_FetchBlogPopular_test.go` | 人気ブログ取得 |
| `repositories/blogs/test/blogs_CreateBlog_test.go` | ブログ作成 |
| `repositories/blogs/test/blogs_UpdateBlog_test.go` | ブログ更新 |
| `repositories/blogs/test/blogs_DeleteBlog_test.go` | ブログ削除 |
| `repositories/blogs/test/blogs_pipeline_test.go` | ブログCRUDパイプライン |
| `repositories/blog_users/blog_users_FetchUserById_test.go` | ユーザーID指定取得 |
| `repositories/blog_users/blog_users_FetchUserByEmailAndPassword_test.go` | メール・パスワード指定取得 |
| `repositories/blog_users/blog_users_UpdateBlogUsers_test.go` | ユーザー更新 |
| `repositories/blog_comments/blog_comments_FetchCommentsByBlogId_test.go` | ブログID指定コメント取得 |
| `repositories/blog_comments/blog_comments_CreateComment_test.go` | コメント作成 |
| `repositories/blog_likes/blog_likes_pipeline_test.go` | いいねCRUDパイプライン |

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

#### Repository層セットアップ（DB接続あり）

```go
func SetupTest(t *testing.T) {
    // 環境変数の読み込み（.env.test）
    godotenv.Load("../../../.env.test")

    // ログ設定の初期化
    logger.InitLogger()

    // Supabaseクライアントの初期化
    supabase.InitSupabase()
}
```

- 環境変数ファイル `.env.test` からDB接続情報を読み込む
- Supabase（PostgreSQL）への接続プールを初期化する
- パスは相対パスで指定（テストファイルの配置場所に依存）

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

### 6.1 Repository テスト（結合テスト）

#### 概要

- 実際のSupabase（PostgreSQL）データベースに接続してテストを実行する
- テストデータは事前に環境変数で指定されたIDを使用する
- CRUD操作の正常系・異常系を検証する

#### 環境変数要件

| 環境変数 | 説明 | 用途 |
|---------|------|------|
| `SUPABASE_URL` | Supabase接続URL | DB接続 |
| `JWT_SECRET_KEY` | JWT署名キー | config初期化 |
| `TEST_USER_ID` | テスト用ユーザーID (UUID) | ユーザー関連テスト |
| `TEST_BLOG_ID` | テスト用ブログID (UUID) | ブログ取得テスト |
| `TEST_MODE` | テストモードフラグ (`true`) | ログ抑制 |

#### テスト対象機能

- 正常系: データ取得、作成、更新、削除
- 異常系: 存在しないID、無効なUUID形式、空文字パラメータ
- パイプラインテスト: CRUD操作の一連の流れ（Create → Read → Update → Delete）

#### テストの特徴

```go
// 正常系: 環境変数からテストデータIDを取得
id := os.Getenv("TEST_BLOG_ID")
blog, err := repo.FetchBlogById(id)
assert.NoError(t, err)
assert.NotNil(t, blog)

// 異常系: 無効なUUID形式のエラーを検証
blog, err := repo.FetchBlogById("2")
assert.Error(t, err)
assert.Nil(t, blog)
assert.Equal(t, "ERROR: invalid input syntax for type uuid: \"2\" (SQLSTATE 22P02)", err.Error())
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

### 全テスト実行

```bash
go test ./... -v
```

### 特定パッケージのテスト実行

```bash
# Repository テスト（DB接続が必要）
go test ./repositories/blogs/test/... -v

# Service テスト
go test ./services/blogs/test/... -v

# Handler テスト
go test ./handlers/blogs/test/... -v
```

### 特定テスト関数の実行

```bash
go test ./handlers/blogs/test/... -v -run TestHandler_FetchBlogs
```

### テストカバレッジの計測

```bash
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
