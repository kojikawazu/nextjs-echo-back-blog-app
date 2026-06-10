# クリーンアーキテクチャ 理解向上レポート

本レポートは、**クリーンアーキテクチャ（Clean Architecture）** について、本プロジェクト（nextjs-echo-back-blog-app）の実装を題材に技術理解を深めるためのドキュメントである。

---

## 目次

- [1. クリーンアーキテクチャとは何か](#1-クリーンアーキテクチャとは何か)
  - [1.1 概要](#11-概要)
  - [1.2 伝統的なクリーンアーキテクチャの同心円](#12-伝統的なクリーンアーキテクチャの同心円)
  - [1.3 核心原則](#13-核心原則)
- [2. 本プロジェクトの3層構成](#2-本プロジェクトの3層構成)
  - [2.1 層の対応関係](#21-層の対応関係)
  - [2.2 各層の責務](#22-各層の責務)
- [3. インターフェースによる依存性逆転](#3-インターフェースによる依存性逆転)
  - [3.1 依存性逆転原則（DIP）とは](#31-依存性逆転原則dipとは)
  - [3.2 本プロジェクトのインターフェース定義](#32-本プロジェクトのインターフェース定義)
  - [3.3 ファイル構成パターン](#33-ファイル構成パターン)
- [4. 依存性注入（DI）の組み立て](#4-依存性注入diの組み立て)
  - [4.1 組み立て場所](#41-組み立て場所)
  - [4.2 依存関係の可視化](#42-依存関係の可視化)
  - [4.3 コンストラクタインジェクション](#43-コンストラクタインジェクション)
- [5. リクエスト処理の完全フロー](#5-リクエスト処理の完全フロー)
  - [5.1 ブログ詳細取得の例](#51-ブログ詳細取得の例)
- [6. エラー伝播パターン](#6-エラー伝播パターン)
  - [6.1 各層でのエラー処理の責務](#61-各層でのエラー処理の責務)
  - [6.2 エラーメッセージベースの分岐](#62-エラーメッセージベースの分岐)
- [7. テスタビリティ](#7-テスタビリティ)
  - [7.1 モックによるレイヤー分離テスト](#71-モックによるレイヤー分離テスト)
  - [7.2 Service テストの例](#72-service-テストの例)
  - [7.3 Handler テストの例](#73-handler-テストの例)
  - [7.4 テストピラミッド](#74-テストピラミッド)
- [8. ドメイン駆動のディレクトリ構成](#8-ドメイン駆動のディレクトリ構成)
  - [8.1 本プロジェクトの構成](#81-本プロジェクトの構成)
  - [8.2 レイヤーごと vs ドメインごとの構成](#82-レイヤーごと-vs-ドメインごとの構成)
- [9. クロスカッティング関心事](#9-クロスカッティング関心事)
  - [9.1 横断的関心事とは](#91-横断的関心事とは)
  - [9.2 ミドルウェアの位置づけ](#92-ミドルウェアの位置づけ)
- [10. クリーンアーキテクチャのトレードオフ](#10-クリーンアーキテクチャのトレードオフ)
  - [10.1 メリットとデメリット](#101-メリットとデメリット)
  - [10.2 本プロジェクトでの適合性](#102-本プロジェクトでの適合性)
- [11. まとめ：学習ポイント](#11-まとめ学習ポイント)

---

## 1. クリーンアーキテクチャとは何か

### 1.1 概要

クリーンアーキテクチャは、Robert C. Martin（Uncle Bob）が2012年に提唱したソフトウェアアーキテクチャの設計原則である。核心は **関心の分離（Separation of Concerns）** と **依存性の方向の統一** にある。

### 1.2 伝統的なクリーンアーキテクチャの同心円

```
┌─────────────────────────────────────────────────────────┐
│                                                         │
│   ┌─────────────────────────────────────────────────┐   │
│   │                                                 │   │
│   │   ┌─────────────────────────────────────────┐   │   │
│   │   │                                         │   │   │
│   │   │   ┌─────────────────────────────────┐   │   │   │
│   │   │   │                                 │   │   │   │
│   │   │   │         Entities                │   │   │   │
│   │   │   │      （ビジネスルール）           │   │   │   │
│   │   │   │                                 │   │   │   │
│   │   │   └─────────────────────────────────┘   │   │   │
│   │   │                                         │   │   │
│   │   │         Use Cases                       │   │   │
│   │   │      （アプリケーションロジック）         │   │   │
│   │   │                                         │   │   │
│   │   └─────────────────────────────────────────┘   │   │
│   │                                                 │   │
│   │         Interface Adapters                      │   │
│   │      （コントローラ、プレゼンタ、ゲートウェイ）   │   │
│   │                                                 │   │
│   └─────────────────────────────────────────────────┘   │
│                                                         │
│         Frameworks & Drivers                            │
│      （Web, DB, UI, 外部サービス）                       │
│                                                         │
└─────────────────────────────────────────────────────────┘

依存の方向: 外側 → 内側（内側は外側を知らない）
```

### 1.3 核心原則

| 原則 | 説明 |
|------|------|
| **依存性ルール** | 依存は外側から内側へのみ。内側の層は外側の層を知らない |
| **関心の分離** | 各層は1つの責務だけを持つ |
| **テスタビリティ** | ビジネスロジックがUI/DB/フレームワークに依存しないため、単体テストが容易 |
| **フレームワーク非依存** | フレームワークを交換しても、ビジネスロジックは変更不要 |
| **DB非依存** | データベースを交換しても、ビジネスロジックは変更不要 |

---

## 2. 本プロジェクトの3層構成

### 2.1 層の対応関係

本プロジェクトではクリーンアーキテクチャを **3層構成** に簡略化して採用している。

```
伝統的クリーンアーキテクチャ          本プロジェクトの3層
───────────────────────            ─────────────────────

Frameworks & Drivers   ←───→     Handler層（プレゼンテーション）
Interface Adapters                  │ Echo, HTTP, Cookie, JSON
                                    │
Use Cases              ←───→     Service層（ビジネスロジック）
                                    │ バリデーション、変換、判定
                                    │
Entities               ←───→     Repository層（データアクセス）
                                    │ SQL, pgx, Supabase
                                    │
                                    ↓
                                 PostgreSQL (Supabase)
```

### 2.2 各層の責務

```
┌────────────────────────────────────────────────────────────────┐
│  Handler層（handlers/）                                        │
│                                                                │
│  責務：                                                         │
│  ・HTTPリクエストの受信とパース                                   │
│  ・パスパラメータ、リクエストボディの取得                          │
│  ・Cookie/JWT トークンの取得と認証チェック                        │
│  ・Service層の呼び出し                                          │
│  ・HTTPレスポンスの生成（ステータスコード、JSONボディ）             │
│  ・エラーメッセージからHTTPステータスコードへの変換                 │
│                                                                │
│  知っていること: echo.Context, HTTP, JSON, Cookie                │
│  知らないこと: SQL, データベース, テーブル構造                     │
└───────────────────────────┬────────────────────────────────────┘
                            │ Interface（依存性逆転）
                            ↓
┌────────────────────────────────────────────────────────────────┐
│  Service層（services/）                                        │
│                                                                │
│  責務：                                                         │
│  ・入力バリデーション（空文字チェック、メール形式チェック等）       │
│  ・ビジネスロジックの実行                                        │
│  ・データ変換（タグのカンマ分割、重複排除、ソート等）              │
│  ・Repository層の呼び出し                                       │
│  ・ドメインエラーの生成                                          │
│                                                                │
│  知っていること: ビジネスルール、バリデーションルール               │
│  知らないこと: HTTP, Cookie, SQL, テーブル構造                   │
└───────────────────────────┬────────────────────────────────────┘
                            │ Interface（依存性逆転）
                            ↓
┌────────────────────────────────────────────────────────────────┐
│  Repository層（repositories/）                                 │
│                                                                │
│  責務：                                                         │
│  ・SQLクエリの構築と実行                                         │
│  ・データベース接続の管理（pgxpool経由）                          │
│  ・クエリ結果のスキャンとモデルへのマッピング                     │
│  ・UUID形式のバリデーション                                      │
│  ・エラーのログ出力                                              │
│                                                                │
│  知っていること: SQL, pgx, テーブル構造, Supabase               │
│  知らないこと: HTTP, Cookie, ビジネスルール                      │
└────────────────────────────────────────────────────────────────┘
```

---

## 3. インターフェースによる依存性逆転

### 3.1 依存性逆転原則（DIP）とは

通常、上位層は下位層に依存する。しかし、インターフェースを使うことで依存の方向を逆転できる。

```
依存性逆転なし（密結合）:
  Handler → Service → Repository（具象クラス）
  ↑ 上位が下位の具体的な実装を直接参照
  ↑ Repository を変更するとService も変更が必要

依存性逆転あり（疎結合）= 本プロジェクト:
  Handler → ServiceInterface ← Service実装
                  ↑                  ↓
              インターフェース    RepositoryInterface ← Repository実装
  ↑ 上位はインターフェースだけに依存
  ↑ Repository の実装を変更してもService は変更不要
```

### 3.2 本プロジェクトのインターフェース定義

#### Repository インターフェース

```go
// repositories/blogs/blogs.go
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

#### Service インターフェース

```go
// services/blogs/blogs.go
type BlogService interface {
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

### 3.3 ファイル構成パターン

各ドメインは統一されたファイル構成を持つ：

```
blogs/
├── blogs.go              ← インターフェース定義 + コンストラクタ
├── blogs_impl.go         ← 実装（本番コード）
└── blogs_mock.go         ← モック実装（テスト用）
```

```go
// blogs.go — インターフェース + コンストラクタ
type BlogRepository interface {
    FetchBlogs() ([]models.BlogData, error)
    // ...
}

type BlogRepositoryImpl struct{}

func NewBlogRepository() BlogRepository {
    return &BlogRepositoryImpl{}
}

// blogs_impl.go — 実装
func (r *BlogRepositoryImpl) FetchBlogs() ([]models.BlogData, error) {
    // 実際のSQL実行
}

// blogs_mock.go — モック
type MockBlogRepository struct {
    mock.Mock
}

func (m *MockBlogRepository) FetchBlogs() ([]models.BlogData, error) {
    args := m.Called()
    // テスト用の振る舞い
}
```

---

## 4. 依存性注入（DI）の組み立て

### 4.1 組み立て場所

全ての依存関係は `routes/routes.go` の `SetupRoutes()` で一箇所にまとめて組み立てられる。

```go
// routes/routes.go
func SetupRoutes(e *echo.Echo) {
    // ===== ユーティリティ =====
    cookieUtils := utils_cookie.NewCookieUtils()

    // ===== Repository層（最下位）=====
    userRepository := repositories_blog_users.NewBlogUsersRepository()
    blogRepository := repositories_blogs.NewBlogRepository()
    blogLikeRepository := repositories_blogs_likes.NewBlogLikeRepository()
    commentRepository := repositories_comments.NewCommentRepository()

    // ===== Service層（中間）=====
    // Repository をインターフェース経由で注入
    authService := services_auth.NewAuthService()
    userService := services_blog_users.NewUserService(userRepository)
    blogService := services_blogs.NewBlogService(blogRepository)
    blogLikeService := services_blogs_likes.NewBlogLikeService(blogLikeRepository)
    commentService := services_comments.NewCommentService(commentRepository)

    // ===== Handler層（最上位）=====
    // Service をインターフェース経由で注入
    authHandler := handlers_auth.NewAuthHandler(userService, authService)
    blogUsersHandler := handlers_blog_users.NewBlogUsersHandler(userService, cookieUtils)
    blogHandler := handlers_blogs.NewBlogHandler(blogService, cookieUtils)
    blogLikeHandler := handlers_blogs_likes.NewBlogLikeHandler(blogLikeService, cookieUtils)
    commentHandler := handlers_comments.NewCommentHandler(commentService)

    // ===== ルーティング定義 =====
    api := e.Group("/api")
    // ... 各エンドポイントの定義
}
```

### 4.2 依存関係の可視化

```
                    SetupRoutes() で組み立て
                    ┌──────────────────────┐
                    │                      │
    ┌───────────────┼──────────────────────┼──────────────────┐
    │               │                      │                  │
    ↓               ↓                      ↓                  ↓

CookieUtils    UserRepo    BlogRepo    LikeRepo    CommentRepo
    │              │           │           │             │
    │              ↓           ↓           ↓             ↓
    │         UserService  BlogService LikeService CommentService
    │              │           │           │             │
    │              │           │           │             │
    ↓              ↓           ↓           ↓             ↓
AuthHandler  BlogUsersH   BlogHandler  LikeHandler  CommentHandler
```

### 4.3 コンストラクタインジェクション

Go言語にはコンストラクタ構文がないため、`NewXxx()` 関数でインターフェースを受け取る。

```go
// Service のコンストラクタ
func NewBlogService(blogRepository repositories_blogs.BlogRepository) BlogService {
    return &BlogServiceImpl{
        BlogRepository: blogRepository,  // インターフェース型で保持
    }
}

// 構造体の定義
type BlogServiceImpl struct {
    BlogRepository repositories_blogs.BlogRepository  // 具象型ではなくインターフェース
}
```

**なぜコンストラクタインジェクションか？**
- 依存関係が明示的（引数を見れば何に依存しているか一目瞭然）
- 初期化後に依存関係が変更されない（不変性）
- テスト時にモックを渡すだけで差し替え可能

---

## 5. リクエスト処理の完全フロー

### 5.1 ブログ詳細取得の例

```
HTTP GET /api/blogs/detail/550e8400-e29b-41d4-a716-446655440000
```

```
Step 1: Handler層（blogs_impl.go）
┌────────────────────────────────────────────────────────┐
│ func (h *BlogHandler) FetchBlogById(c echo.Context)    │
│                                                        │
│   // パスパラメータを取得                                │
│   id := c.Param("id")                                  │
│                                                        │
│   // Service層を呼び出す（インターフェース経由）          │
│   blog, err := h.BlogService.FetchBlogById(id)         │
│                                                        │
│   // エラーハンドリング                                  │
│   if err != nil {                                      │
│       switch err.Error() {                             │
│       case "invalid id":                               │
│           return c.JSON(400, ...)                       │
│       case "blog not found":                           │
│           return c.JSON(404, ...)                       │
│       default:                                         │
│           return c.JSON(500, ...)                       │
│       }                                                │
│   }                                                    │
│                                                        │
│   // 成功レスポンス                                      │
│   return c.JSON(200, blog)                             │
└───────────────────────────┬────────────────────────────┘
                            │
                            ↓
Step 2: Service層（blogs_impl.go）
┌────────────────────────────────────────────────────────┐
│ func (s *BlogServiceImpl) FetchBlogById(id string)     │
│                                                        │
│   // バリデーション                                      │
│   if id == "" {                                        │
│       return nil, errors.New("invalid id")             │
│   }                                                    │
│                                                        │
│   // Repository層を呼び出す（インターフェース経由）       │
│   blog, err := s.BlogRepository.FetchBlogById(id)      │
│                                                        │
│   if err != nil {                                      │
│       return nil, errors.New("blog not found")         │
│   }                                                    │
│                                                        │
│   return blog, nil                                     │
└───────────────────────────┬────────────────────────────┘
                            │
                            ↓
Step 3: Repository層（blogs_impl.go）
┌────────────────────────────────────────────────────────┐
│ func (r *BlogRepositoryImpl) FetchBlogById(id string)  │
│                                                        │
│   query := `SELECT ... FROM blogs b                    │
│              LEFT JOIN ... WHERE b.id = $1`            │
│                                                        │
│   row := supabase.Pool.QueryRow(supabase.Ctx, query,   │
│                                 id)                    │
│                                                        │
│   var blog models.BlogData                             │
│   err := row.Scan(&blog.ID, &blog.Title, ...)          │
│                                                        │
│   return &blog, err                                    │
└───────────────────────────┬────────────────────────────┘
                            │
                            ↓
Step 4: Supabase (PostgreSQL)
┌────────────────────────────────────────────────────────┐
│   SELECT b.id, b.title, b.description, ...             │
│   FROM blogs b                                         │
│   LEFT JOIN (SELECT blog_id, COUNT(*) ...)              │
│   WHERE b.id = '550e8400-...'                          │
│                                                        │
│   → 1行の結果を返却                                     │
└────────────────────────────────────────────────────────┘
```

---

## 6. エラー伝播パターン

### 6.1 各層でのエラー処理の責務

```
Repository層:                    Service層:                    Handler層:
  DBエラーをそのまま返す →         ドメインエラーに変換 →         HTTPステータスに変換
  "ERROR: invalid UUID"           "blog not found"              404 Not Found
  sql.ErrNoRows                   "invalid id"                  400 Bad Request
  接続エラー                       "failed to create blog"       500 Internal Error
```

### 6.2 エラーメッセージベースの分岐

本プロジェクトでは、エラーの種別判定にエラーメッセージの文字列比較を使用している。

```go
// Handler層のエラーハンドリング
switch err.Error() {
case "invalid id":
    return c.JSON(http.StatusBadRequest,
        map[string]string{"error": "Invalid id"})
case "blog not found":
    return c.JSON(http.StatusNotFound,
        map[string]string{"error": "Blog not found"})
default:
    return c.JSON(http.StatusInternalServerError,
        map[string]string{"error": "Error fetching blog"})
}
```

**この方式の特徴：**

| 長所 | 短所 |
|------|------|
| 実装がシンプル | エラーメッセージの変更がAPI挙動に影響 |
| 追加のエラー型定義不要 | typo（タイプミス）で分岐が壊れる |
| 理解しやすい | IDEの静的解析が効かない |

**改善案（カスタムエラー型）：**
```go
// 改善案
type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

var ErrBlogNotFound = &AppError{Code: 404, Message: "blog not found"}
var ErrInvalidID = &AppError{Code: 400, Message: "invalid id"}
```

---

## 7. テスタビリティ

### 7.1 モックによるレイヤー分離テスト

クリーンアーキテクチャの最大の利点の1つが、各層を独立してテストできること。

```
テスト対象          依存関係のモック化
────────          ─────────────────

Handler テスト:   MockService + MockCookieUtils → 実Service不要
  ↓                ↓
  ├── HTTPリクエストの生成（httptest）
  ├── ハンドラの実行
  └── HTTPレスポンスの検証（ステータスコード、ボディ）

Service テスト:   MockRepository → 実DB不要
  ↓                ↓
  ├── バリデーションロジックの検証
  ├── データ変換ロジックの検証
  └── Repository呼び出しの検証

Repository テスト: 実DB（Supabase）→ モックなし（結合テスト）
  ↓
  ├── SQL実行の検証
  ├── スキャン結果の検証
  └── エラーケースの検証
```

### 7.2 Service テストの例

```go
// services/blogs/test/blogs_FetchBlogById_test.go
func TestService_FetchBlogById_Success(t *testing.T) {
    // (1) モックリポジトリを作成
    mockRepo := new(repositories_blogs.MockBlogRepository)

    // (2) Service にモックを注入
    blogService := services_blogs.NewBlogService(mockRepo)

    // (3) モックの振る舞いを定義
    expectedBlog := &models.BlogData{
        ID:    "test-id-123",
        Title: "テストブログ",
    }
    mockRepo.On("FetchBlogById", "test-id-123").Return(expectedBlog, nil)

    // (4) Service メソッドを実行
    result, err := blogService.FetchBlogById("test-id-123")

    // (5) 結果を検証
    assert.NoError(t, err)
    assert.Equal(t, "テストブログ", result.Title)

    // (6) モックが期待通り呼ばれたか確認
    mockRepo.AssertExpectations(t)
}

func TestService_FetchBlogById_EmptyId(t *testing.T) {
    mockRepo := new(repositories_blogs.MockBlogRepository)
    blogService := services_blogs.NewBlogService(mockRepo)

    // 空IDでバリデーションエラーになることを確認
    result, err := blogService.FetchBlogById("")

    assert.Error(t, err)
    assert.Equal(t, "invalid id", err.Error())
    assert.Nil(t, result)

    // リポジトリは呼ばれないことを確認
    mockRepo.AssertNotCalled(t, "FetchBlogById")
}
```

### 7.3 Handler テストの例

```go
// handlers/blogs/test/blogs_FetchBlogById_test.go
func TestHandler_FetchBlogById_Success(t *testing.T) {
    // (1) Echo テストコンテキストを作成
    e := echo.New()
    req := httptest.NewRequest(http.MethodGet, "/api/blogs/detail/test-id", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetParamNames("id")
    c.SetParamValues("test-id")

    // (2) モックService と モックCookie を作成
    mockService := new(services_blogs.MockBlogService)
    mockCookie := new(utils_cookie.MockCookieUtils)
    handler := handlers_blogs.NewBlogHandler(mockService, mockCookie)

    // (3) モックの振る舞いを定義
    blog := &models.BlogData{ID: "test-id", Title: "テスト"}
    mockService.On("FetchBlogById", "test-id").Return(blog, nil)

    // (4) ハンドラを実行
    err := handler.FetchBlogById(c)

    // (5) HTTPレスポンスを検証
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
    assert.Contains(t, rec.Body.String(), "テスト")

    mockService.AssertExpectations(t)
}
```

### 7.4 テストピラミッド

```
テスト量と実行速度の関係：

        /  Handler テスト  \        ← 最速（モックのみ）
       /                    \         HTTPリクエスト/レスポンスの検証
      /  Service テスト      \      ← 高速（モックのみ）
     /                        \       ビジネスロジックの検証
    /  Repository テスト        \    ← 低速（実DB接続）
   /                            \     データアクセスの検証
  /______________________________\

  上に行くほど: テスト実行が速い、依存が少ない
  下に行くほど: テスト実行が遅い、実環境に近い
```

---

## 8. ドメイン駆動のディレクトリ構成

### 8.1 本プロジェクトの構成

```
nextjs-echo-back-blog-app/
│
├── handlers/                    ← Handler層
│   ├── auth/                   ← 認証ドメイン
│   ├── blogs/                  ← ブログドメイン
│   ├── blog_users/             ← ユーザードメイン
│   ├── blog_likes/             ← いいねドメイン
│   └── blog_comments/          ← コメントドメイン
│
├── services/                    ← Service層
│   ├── auth/
│   ├── blogs/
│   ├── blog_users/
│   ├── blog_likes/
│   └── blog_comments/
│
├── repositories/                ← Repository層
│   ├── auth/
│   ├── blogs/
│   ├── blog_users/
│   ├── blog_likes/
│   └── blog_comments/
│
├── models/                      ← データモデル（層をまたぐ共有型）
│   ├── auth.go
│   ├── blogs.go
│   ├── blog_users.go
│   ├── blog_likes.go
│   └── blog_comments.go
│
├── routes/                      ← DI組み立て + ルーティング
├── middlewares/                  ← 横断的関心事
├── config/                      ← 設定
├── logger/                      ← ログ
├── supabase/                    ← DB接続
└── utils/                       ← ユーティリティ
```

### 8.2 レイヤーごと vs ドメインごとの構成

```
レイヤーごと（採用しなかった方式）:       ドメインごと（本プロジェクト）:

handlers/                              handlers/
├── auth_handler.go                    ├── auth/
├── blogs_handler.go                   │   ├── auth.go
├── users_handler.go                   │   ├── auth_impl.go
├── likes_handler.go                   │   └── auth_test.go
└── comments_handler.go                ├── blogs/
                                       │   ├── blogs.go
services/                              │   ├── blogs_impl.go
├── auth_service.go                    │   ├── blogs_mock.go
├── blogs_service.go                   │   └── test/
├── users_service.go                   │       ├── blogs_FetchBlogs_test.go
├── likes_service.go                   │       └── ...
└── comments_service.go                └── ...
```

**ドメインごとの構成の利点：**
- 関連ファイルが近くに配置され、ナビゲーションが容易
- 新しいドメインの追加がディレクトリ作成だけで完結
- ドメインの削除も安全（他のドメインに影響しない）

---

## 9. クロスカッティング関心事

### 9.1 横断的関心事とは

複数のドメインにまたがる共通機能のこと。クリーンアーキテクチャでは、これらを特定の層やドメインから分離する。

```
横断的関心事                配置場所
──────────                ────────
ログ出力                   logger/, utils/log/
Cookie管理                 utils/cookie/
CORS設定                   middlewares/
パニックリカバリ            middlewares/
設定管理                   config/
DB接続管理                 supabase/
データモデル               models/
```

### 9.2 ミドルウェアの位置づけ

```go
// middlewares/middlewares.go
func SetupMiddlewares(e *echo.Echo) {
    e.Use(middleware.Logger())   // 全リクエストのアクセスログ
    e.Use(middleware.Recover())  // パニックの自動リカバリ
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins:     strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowCredentials: true,
        // ...
    }))
}
```

ミドルウェアは Handler層の手前で動作し、全てのリクエストに共通の処理を適用する。

---

## 10. クリーンアーキテクチャのトレードオフ

### 10.1 メリットとデメリット

| メリット | デメリット |
|---------|----------|
| 各層を独立してテスト可能 | ファイル数が多い（インターフェース、実装、モック） |
| フレームワーク交換が容易 | 単純な機能でもボイラープレートが必要 |
| ビジネスロジックが純粋 | 学習コストが高い |
| チーム開発で責務が明確 | 層間のデータ変換コストがかかる場合がある |
| 長期的な保守性が高い | 小規模プロジェクトではオーバーエンジニアリング |

### 10.2 本プロジェクトでの適合性

本プロジェクトは個人ブログのバックエンドであり、規模としては小〜中規模である。クリーンアーキテクチャを採用することで：

- **学習目的**: アーキテクチャパターンの実践的な学習になる
- **将来の拡張性**: 新しいドメイン（例：タグ管理、通知機能）の追加が容易
- **テストの信頼性**: モックベースの単体テストで品質を担保

---

## 11. まとめ：学習ポイント

| # | ポイント | 本プロジェクトでの実践 |
|---|---------|---------------------|
| 1 | 依存は外側→内側の一方向のみ | Handler→Service→Repository（逆方向の依存なし） |
| 2 | インターフェースで層間を疎結合に | 全ドメインでインターフェース定義 |
| 3 | コンストラクタインジェクションでDI | `NewXxxService(repo)` パターン |
| 4 | DI組み立ては1箇所に集約 | `routes/routes.go` の `SetupRoutes()` |
| 5 | 各ドメインにインターフェース/実装/モックの3ファイル | `xxx.go` / `xxx_impl.go` / `xxx_mock.go` |
| 6 | Handler: HTTPの関心事のみ | リクエスト解析、レスポンス生成、ステータスコード |
| 7 | Service: ビジネスルールのみ | バリデーション、データ変換、判定ロジック |
| 8 | Repository: データアクセスのみ | SQL実行、結果スキャン、エラーログ |
| 9 | テストはモックで層を分離 | Service→MockRepo、Handler→MockService |
| 10 | ドメインごとのディレクトリ構成 | blogs/, blog_users/, blog_likes/, blog_comments/, auth/ |
