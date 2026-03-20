# 機能仕様書 (Functional Specification)

本ドキュメントは、ブログWebアプリケーションバックエンド（nextjs-echo-back-blog-app）の機能仕様を定義する。各機能の処理フロー、ビジネスロジック、バリデーションルールを既存のソースコードからリバースエンジニアリングにより詳細に記述する。

---

## 1. 認証機能

### 1.1 認証フロー全体像

```
[フロントエンド]
    |
    | POST /api/blog-users/login (email, password)
    v
[AuthHandler.Login]
    |
    | バリデーション（AuthService.Login）
    | ユーザー取得（UserService.FetchUserByEmailAndPassword）
    | JWTトークン生成
    | HTTP-Only Cookie設定
    v
[レスポンス: 200 OK + Set-Cookie: token]
    |
    | GET /api/blog-users/auth-check (Cookie: token)
    v
[AuthHandler.CheckAuth]
    |
    | Cookie取得 → JWT解析 → 有効性検証
    v
[レスポンス: 200 OK + user_id, username, email]
    |
    | POST /api/blog-users/logout
    v
[AuthHandler.Logout]
    |
    | Cookie削除（有効期限を過去に設定）
    v
[レスポンス: 200 OK]
```

### 1.2 ログイン処理詳細

#### 処理フロー

1. リクエストボディから `email` と `password` をJSON形式で取得
2. `AuthService.Login` でバリデーションを実行:
   - `email` と `password` が空でないことを確認
   - `email` が `net/mail.ParseAddress` による有効なメール形式であることを確認
3. `UserService.FetchUserByEmailAndPassword` でデータベースからユーザーを検索:
   - `email` と `password` の両方が一致するユーザーをSELECT
   - パスワードはDBに平文保存されているため、直接比較
4. JWTトークン生成:
   - ペイロード: `user_id`, `email`, `username`
   - 有効期限: 現在時刻 + 1時間
   - 署名アルゴリズム: HS256
   - 署名鍵: 環境変数 `JWT_SECRET_KEY`
5. Cookie設定:
   - Cookie名: `token`
   - HttpOnly: `true`
   - Path: `/`
   - 本番環境: `Secure=true`, `SameSite=None`
   - 開発環境: `Secure=false`, `SameSite=Lax`

#### エラーハンドリング

| 条件 | HTTPステータス | エラーメッセージ |
|------|---------------|-----------------|
| リクエストボディ不正 | 400 | Invalid request body |
| email/password空 | 400 | Email and password are required |
| メール形式不正 | 400 | Invalid email format |
| ユーザー未発見 | 404 | User not found |
| トークン生成失敗 | 500 | Could not create token |

### 1.3 認証確認処理詳細

#### 処理フロー

1. `token` Cookieから値を取得
2. `jwt.ParseWithClaims` でトークンを解析（Claims構造体にデコード）
3. トークンの有効性を検証（`token.Valid`）
4. 認証成功時、Claimsから `user_id`, `username`, `email` を返却

#### エラーハンドリング

| 条件 | HTTPステータス | エラーメッセージ |
|------|---------------|-----------------|
| Cookieなし | 401 | Token not found |
| トークン解析失敗 | 401 | Invalid token |
| トークン無効 | 401 | Invalid token |

### 1.4 ログアウト処理詳細

#### 処理フロー

1. `token` Cookieの値を空文字に設定
2. 有効期限を `time.Unix(0, 0)`（1970-01-01）に設定
3. Cookie属性（HttpOnly, Path, Secure, SameSite）は認証Cookieと同一設定

---

## 2. ブログ管理機能

### 2.1 ブログ作成フロー

```
[フロントエンド]
    |
    | POST /api/blogs/create (title, githubUrl, category, description, tags)
    | Cookie: token
    v
[BlogHandler.CreateBlog]
    |
    | 1. Cookie "token" からJWTトークン取得
    | 2. JWTトークン解析 → userId取得
    | 3. リクエストボディのバインド
    | 4. BlogService.CreateBlog 呼び出し
    v
[BlogService.CreateBlog]
    |
    | バリデーション（全フィールド非空チェック）
    | BlogRepository.CreateBlog 呼び出し
    v
[BlogRepository.CreateBlog]
    |
    | UUID形式バリデーション（userId）
    | INSERT INTO blogs ... RETURNING
    v
[レスポンス: 201 Created + ブログデータ]
```

#### バリデーションルール（Service層）

| フィールド | ルール | エラーメッセージ |
|------------|--------|-----------------|
| userId | 空文字でないこと | invalid userId |
| title | 空文字でないこと | invalid title |
| githubUrl | 空文字でないこと | invalid githubUrl |
| category | 空文字でないこと | invalid category |
| description | 空文字でないこと | invalid description |
| tags | 空文字でないこと | invalid tags |

#### バリデーションルール（Repository層）

| フィールド | ルール | エラーメッセージ |
|------------|--------|-----------------|
| userId | 空文字でないこと | user_id cannot be empty |
| userId | UUID形式であること | invalid user_id format: must be a valid UUID |

#### SQLクエリ

```sql
INSERT INTO blogs (blog_user_id, title, github_url, category, description, tags)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, blog_user_id, title, description, github_url, category, tags, likes, comment_cnt, created_at, updated_at
```

### 2.2 ブログ更新フロー

```
[フロントエンド]
    |
    | PUT /api/blogs/update/:id (title, githubUrl, category, description, tags)
    | Cookie: token
    v
[BlogHandler.UpdateBlog]
    |
    | 1. Cookie "token" からJWTトークン取得
    | 2. JWTトークン解析 → userId取得（認証確認のみ、値は未使用）
    | 3. パスパラメータから id 取得
    | 4. リクエストボディのバインド
    | 5. BlogService.UpdateBlog 呼び出し
    v
[BlogService.UpdateBlog]
    |
    | バリデーション（全フィールド非空チェック）
    | BlogRepository.UpdateBlog 呼び出し
    v
[BlogRepository.UpdateBlog]
    |
    | CTEを使用したUPDATE + いいね数・コメント数集計
    v
[レスポンス: 200 OK + 更新済みブログデータ]
```

#### SQLクエリ

```sql
WITH updated_blog AS (
    UPDATE blogs
    SET title = $2, github_url = $3, category = $4, description = $5, tags = $6
    WHERE id = $1
    RETURNING id, blog_user_id, title, description, github_url, category, tags, created_at, updated_at
)
SELECT ub.id, ub.blog_user_id, ub.title, ub.description, ub.github_url, ub.category, ub.tags,
       COALESCE(l.like_count, 0) AS likes,
       COALESCE(c.comment_count, 0) AS comment_cnt,
       ub.created_at, ub.updated_at
FROM updated_blog ub
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS like_count
    FROM blog_likes
    GROUP BY blog_id
) l ON ub.id = l.blog_id
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS comment_count
    FROM blog_comments
    GROUP BY blog_id
) c ON ub.id = c.blog_id
```

**設計ポイント**: CTEを使用し、UPDATE実行と同時にいいね数・コメント数の集計結果を返却する。

### 2.3 ブログ削除フロー

```
[フロントエンド]
    |
    | DELETE /api/blogs/delete/:id
    | Cookie: token
    v
[BlogHandler.DeleteBlog]
    |
    | 1. Cookie "token" からJWTトークン取得
    | 2. JWTトークン解析 → userId取得（認証確認のみ）
    | 3. パスパラメータから id 取得
    | 4. BlogService.DeleteBlog 呼び出し
    v
[BlogService.DeleteBlog]
    |
    | バリデーション（id非空チェック）
    | BlogRepository.DeleteBlog 呼び出し
    v
[BlogRepository.DeleteBlog]
    |
    | UUID形式バリデーション（id）
    | DELETE FROM blogs WHERE id = $1
    v
[レスポンス: 204 No Content]
```

### 2.4 ブログ一覧取得

#### 全ブログ取得

認証不要。全てのブログを `created_at` 降順で取得する。

**SQLクエリ**: `blogs` テーブルに対して `blog_likes` および `blog_comments` をLEFT JOINし、いいね数とコメント数をサブクエリで集計する。

```sql
SELECT b.id, b.blog_user_id, b.title, b.description, b.github_url, b.category, b.tags,
       COALESCE(l.like_count, 0) AS likes,
       COALESCE(c.comment_count, 0) AS comment_cnt,
       b.created_at, b.updated_at
FROM blogs b
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS like_count FROM blog_likes GROUP BY blog_id
) l ON b.id = l.blog_id
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS comment_count FROM blog_comments GROUP BY blog_id
) c ON b.id = c.blog_id
ORDER BY b.created_at DESC
```

**設計ポイント**: いいね数とコメント数はblogsテーブルの `likes`, `comment_cnt` カラムではなく、`blog_likes` および `blog_comments` テーブルからリアルタイムに集計される。

#### 人気ブログ取得

認証不要。いいね数の多い順に指定件数のブログを取得する。

**バリデーション（Service層）**: `count` が0以下の場合は `invalid count` エラー

**SQLクエリ**: 返却フィールドは `id`, `blog_user_id`, `title`, `likes`, `created_at`, `updated_at` のみ（軽量レスポンス）。

```sql
SELECT b.id, b.blog_user_id, b.title,
       COALESCE(l.like_count, 0) AS likes,
       b.created_at, b.updated_at
FROM blogs b
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS like_count FROM blog_likes GROUP BY blog_id
) l ON b.id = l.blog_id
ORDER BY likes DESC
LIMIT $1
```

### 2.5 カテゴリ・タグ取得

#### カテゴリ取得

DBから `DISTINCT category` を `ORDER BY category` で取得し、そのまま返却する。

#### タグ取得（サービス層のビジネスロジック）

1. DBから `DISTINCT tags` を取得（各行はカンマ区切りの文字列）
2. サービス層で以下の加工処理を実行:
   - 各行のタグ文字列をカンマ（`,`）で分割
   - 各タグの前後空白を `strings.TrimSpace` で除去
   - 空文字タグを除外
   - `map[string]struct{}` を使用して重複排除
   - `sort.Strings` でアルファベット順にソート
3. ユニークなタグのスライスを返却

---

## 3. いいねシステム

### 3.1 訪問者ID管理フロー

```
[訪問者（フロントエンド）]
    |
    | GET /api/blog-likes/generate-visit-id
    v
[BlogLikeHandler.GenerateVisitorId]
    |
    |-- Cookie "visit-id-token" 存在チェック
    |   |
    |   |-- 有効なトークンあり → 200 OK "Visitor id already exists"（スキップ）
    |   |
    |   |-- トークンなし or 無効 → 新規生成
    |       |
    |       | 1. uuid.New() で訪問者ID生成
    |       | 2. JWTトークン作成（VisitId をペイロードに含む）
    |       |    - 有効期限: 1時間
    |       |    - 署名: HS256 + JWT_SECRET_KEY
    |       | 3. "visit-id-token" Cookieに設定
    |       v
    |       200 OK "Visitor id generated successfully"
```

#### Cookie設定（visit-id-token）

| 属性 | 本番環境 | 開発環境 |
|------|----------|----------|
| Name | visit-id-token | visit-id-token |
| HttpOnly | true | true |
| Path | / | / |
| Secure | true | false |
| SameSite | None | Lax |
| Expires | 1時間後 | 1時間後 |

### 3.2 いいね追加フロー

```
[訪問者（フロントエンド）]
    |
    | POST /api/blog-likes/create/:blogId
    | Cookie: visit-id-token
    v
[BlogLikeHandler.CreateBlogLike]
    |
    | 1. Cookie "visit-id-token" からJWTトークン取得
    | 2. JWTトークン解析 → visitId取得
    | 3. パスパラメータから blogId 取得
    | 4. BlogLikeService.CreateBlogLike 呼び出し
    v
[BlogLikeService.CreateBlogLike]
    |
    | 1. バリデーション: blogId, visitId が空でないこと
    | 2. 重複チェック: IsBlogLiked(blogId, visitId)
    |    |
    |    |-- 既にいいね済み → エラー "blog is already liked"
    |    |-- 未いいね → 続行
    | 3. BlogLikeRepository.CreateBlogLike 呼び出し
    v
[BlogLikeRepository.CreateBlogLike]
    |
    | INSERT INTO blog_likes (blog_id, visit_id)
    | VALUES ($1, $2)
    | RETURNING id, created_at, updated_at
    v
[レスポンス: 200 OK + いいねデータ]
```

#### エラーハンドリング

| 条件 | HTTPステータス | エラーメッセージ |
|------|---------------|-----------------|
| visit-id-token Cookie取得失敗 | 500 | Failed to get visit id token |
| visitId解析失敗 | 500 | Failed to get visit id |
| blogId/visitId空 | 400 | BlogId or VisitId is empty |
| 重複いいね | 400 | Blog is already liked |
| DB挿入失敗 | 500 | Error creating blog like |

### 3.3 いいね削除フロー

```
[訪問者（フロントエンド）]
    |
    | DELETE /api/blog-likes/delete/:blogId
    | Cookie: visit-id-token
    v
[BlogLikeHandler.DeleteBlogLike]
    |
    | 1. Cookie "visit-id-token" → visitId取得
    | 2. パスパラメータから blogId 取得
    | 3. BlogLikeService.DeleteBlogLike 呼び出し
    v
[BlogLikeRepository.DeleteBlogLike]
    |
    | DELETE FROM blog_likes WHERE blog_id = $1 AND visit_id = $2
    v
[レスポンス: 200 OK]
```

### 3.4 いいね状態確認フロー

```
[訪問者（フロントエンド）]
    |
    | GET /api/blog-likes/is-liked/:blogId
    | Cookie: visit-id-token
    v
[BlogLikeHandler.IsBlogLiked]
    |
    | 1. Cookie "visit-id-token" → visitId取得
    | 2. パスパラメータから blogId 取得
    | 3. BlogLikeService.IsBlogLiked 呼び出し
    v
[BlogLikeRepository.IsBlogLiked]
    |
    | SELECT id FROM blog_likes WHERE blog_id = $1 AND visit_id = $2
    |
    |-- レコード存在 → {"isLiked": true}
    |-- レコード不在 → {"isLiked": false}
    v
[レスポンス: 200 OK + {"isLiked": bool}]
```

**設計ポイント**: エラー発生時（レコード未発見）でもステータスコードは `200 OK` で `{"isLiked": false}` を返す。

---

## 4. コメントシステム

### 4.1 コメント投稿フロー

```
[ゲストユーザー（フロントエンド）]
    |
    | POST /api/comments/create
    | Body: { blogId, guestUser, comment }
    v
[CommentHandler.CreateComment]
    |
    | 1. リクエストボディのバインド
    | 2. CommentService.CreateComment 呼び出し
    v
[CommentService.CreateComment]
    |
    | バリデーション:
    |   - blogId: 空文字でないこと
    |   - guestUser: 空文字でないこと
    |   - comment: 空文字でないこと
    | CommentRepository.CreateComment 呼び出し
    v
[CommentRepository.CreateComment]
    |
    | INSERT INTO blog_comments (blog_id, guest_user, comment)
    | VALUES ($1, $2, $3)
    | RETURNING id, blog_id, guest_user, comment, created_at
    v
[レスポンス: 201 Created + コメントデータ]
```

#### バリデーションルール

| フィールド | ルール | エラーメッセージ | HTTPステータス |
|------------|--------|-----------------|---------------|
| blogId | 空文字でないこと | Invalid blogId | 400 |
| guestUser | 空文字でないこと | Invalid guestUser | 400 |
| comment | 空文字でないこと | Invalid comment | 400 |

### 4.2 コメント一覧取得フロー

```
[訪問者（フロントエンド）]
    |
    | GET /api/comments/blog/:blogId
    v
[CommentHandler.FetchCommentsByBlogId]
    |
    | 1. パスパラメータから blogId 取得
    | 2. CommentService.FetchCommentsByBlogId 呼び出し
    v
[CommentService.FetchCommentsByBlogId]
    |
    | バリデーション: blogId 空文字チェック
    | CommentRepository.FetchCommentsByBlogId 呼び出し
    v
[CommentRepository.FetchCommentsByBlogId]
    |
    | SELECT id, blog_id, guest_user, comment, created_at
    | FROM blog_comments WHERE blog_id = $1
    v
[レスポンス: 200 OK + コメントデータ配列]
```

---

## 5. ユーザープロフィール管理

### 5.1 プロフィール取得フロー

```
[ブログオーナー（フロントエンド）]
    |
    | GET /api/blog-users/detail
    | Cookie: token
    v
[BlogUsersHandler.FetchBlogUsers]
    |
    | 1. Cookie "token" からJWTトークン取得
    | 2. JWTトークン解析 → userId取得
    | 3. UserService.FetchUserById 呼び出し
    | 4. レスポンスから password フィールドを空文字に設定
    v
[レスポンス: 200 OK + ユーザーデータ（password除外）]
```

**設計ポイント**: パスワードはDB取得後にハンドラー層で `user.Password = ""` と設定してからレスポンスに含める。

### 5.2 プロフィール更新フロー

```
[ブログオーナー（フロントエンド）]
    |
    | PUT /api/blog-users/update
    | Cookie: token
    | Body: { name, email, password, newPassword }
    v
[BlogUsersHandler.UpdateBlogUsers]
    |
    | 1. Cookie "token" → userId取得
    | 2. リクエストボディのバインド
    | 3. UserService.UpdateUser 呼び出し
    v
[UserService.UpdateUser]
    |
    | 1. バリデーション:
    |    - id: 空文字でないこと
    |    - name: 空文字でないこと
    |    - email: 空文字でないこと
    |    - password: 空文字でないこと
    |    - newPassword: 空文字でないこと
    |    - email: net/mail.ParseAddress による形式チェック
    | 2. 現パスワード検証:
    |    - FetchBlogUsersById(id) で現在のユーザー情報取得
    |    - currentUser.Password != password → エラー
    | 3. UserRepository.UpdateBlogUsers(id, name, email, newPassword) 呼び出し
    v
[BlogUsersHandler.UpdateBlogUsers（続き）]
    |
    | 4. 更新成功後:
    |    - 新しいJWTトークン生成（更新後のユーザー情報で）
    |    - Cookie "token" を新しいトークンで更新
    v
[レスポンス: 200 OK + 更新済みユーザーデータ + Set-Cookie: token（更新）]
```

#### バリデーションルール

| 順序 | フィールド | ルール | エラーメッセージ |
|------|------------|--------|-----------------|
| 1 | id | 空文字でないこと | id is required |
| 2 | name | 空文字でないこと | name is required |
| 3 | email | 空文字でないこと | email is required |
| 4 | password | 空文字でないこと | password is required |
| 5 | newPassword | 空文字でないこと | new password is required |
| 6 | email | net/mail.ParseAddress形式 | invalid email format |
| 7 | password | DB上の現パスワードと一致 | invalid current password |

**設計ポイント**: プロフィール更新時にJWTトークンを再発行する。これにより、ユーザー名やメールアドレスの変更がトークンのペイロードに即座に反映される。

---

## 6. 共通仕様

### 6.1 APIエンドポイント一覧

| メソッド | パス | 機能 | 認証 |
|----------|------|------|------|
| GET | `/` | ヘルスチェック | 不要 |
| POST | `/api/blog-users/login` | ログイン | 不要 |
| GET | `/api/blog-users/auth-check` | 認証確認 | 要（token） |
| POST | `/api/blog-users/logout` | ログアウト | 不要 |
| GET | `/api/blog-users/detail` | プロフィール取得 | 要（token） |
| PUT | `/api/blog-users/update` | プロフィール更新 | 要（token） |
| GET | `/api/blogs` | 全ブログ取得 | 不要 |
| GET | `/api/blogs/user/:userId` | ユーザー別ブログ取得 | 不要 |
| GET | `/api/blogs/detail/:id` | ブログ詳細取得 | 不要 |
| GET | `/api/blogs/categories` | カテゴリ一覧取得 | 不要 |
| GET | `/api/blogs/tags` | タグ一覧取得 | 不要 |
| GET | `/api/blogs/popular/:count` | 人気ブログ取得 | 不要 |
| POST | `/api/blogs/create` | ブログ作成 | 要（token） |
| PUT | `/api/blogs/update/:id` | ブログ更新 | 要（token） |
| DELETE | `/api/blogs/delete/:id` | ブログ削除 | 要（token） |
| GET | `/api/blog-likes` | いいね一覧取得 | 要（visit-id-token） |
| GET | `/api/blog-likes/generate-visit-id` | 訪問者ID生成 | 不要 |
| GET | `/api/blog-likes/is-liked/:blogId` | いいね状態確認 | 要（visit-id-token） |
| POST | `/api/blog-likes/create/:blogId` | いいね追加 | 要（visit-id-token） |
| DELETE | `/api/blog-likes/delete/:blogId` | いいね削除 | 要（visit-id-token） |
| GET | `/api/comments/blog/:blogId` | コメント一覧取得 | 不要 |
| POST | `/api/comments/create` | コメント投稿 | 不要 |

### 6.2 認証方式の分類

本システムには2種類のCookieベース認証がある。

#### 認証用Cookie（token）

- **目的**: ブログオーナーの認証
- **Cookie名**: `token`
- **用途**: ブログCRUD、プロフィール管理
- **ペイロード**: `user_id`, `email`, `username`
- **有効期限**: 1時間

#### 訪問者ID Cookie（visit-id-token）

- **目的**: 匿名訪問者の識別
- **Cookie名**: `visit-id-token`
- **用途**: いいね機能
- **ペイロード**: `visit_id`（UUID）
- **有効期限**: 1時間

### 6.3 レスポンス形式

#### 成功レスポンス

- 単一リソース取得: `200 OK` + JSONオブジェクト
- リスト取得: `200 OK` + JSON配列
- リソース作成: `201 Created` + 作成されたJSONオブジェクト
- リソース削除: `204 No Content`（ブログ削除のみ）、`200 OK` + メッセージ（いいね削除）

#### エラーレスポンス

全てのエラーは以下の形式で返される:

```json
{
  "error": "エラーメッセージ"
}
```

または認証確認系:

```json
{
  "message": "エラーメッセージ"
}
```

### 6.4 CORS設定

| 項目 | 設定値 |
|------|--------|
| AllowOrigins | 環境変数 `ALLOWED_ORIGINS`（カンマ区切り） |
| AllowMethods | GET, POST, PUT, DELETE |
| AllowHeaders | Origin, Content-Type, Authorization, Access-Control-Allow-Credentials |
| AllowCredentials | true |

### 6.5 依存性注入（DI）構成

```
routes.SetupRoutes
├── CookieUtils (utils_cookie.NewCookieUtils)
├── Repositories
│   ├── BlogUsersRepository (NewBlogUsersRepository)
│   ├── BlogRepository (NewBlogRepository)
│   ├── BlogLikeRepository (NewBlogLikeRepository)
│   └── CommentRepository (NewCommentRepository)
├── Services
│   ├── AuthService (NewAuthService)
│   ├── UserService (NewUserService ← BlogUsersRepository)
│   ├── BlogService (NewBlogService ← BlogRepository)
│   ├── BlogLikeService (NewBlogLikeService ← BlogLikeRepository)
│   └── CommentService (NewCommentService ← CommentRepository)
└── Handlers
    ├── AuthHandler (NewAuthHandler ← UserService, AuthService)
    ├── BlogUsersHandler (NewBlogUsersHandler ← UserService, CookieUtils)
    ├── BlogHandler (NewBlogHandler ← BlogService, CookieUtils)
    ├── BlogLikeHandler (NewBlogLikeHandler ← BlogLikeService, CookieUtils)
    └── CommentHandler (NewCommentHandler ← CommentService)
```

全てのレイヤー間はインターフェースで定義されており、テスト時にはモック実装に差し替えが可能。各ドメイン（auth, blogs, blog_users, blog_likes, blog_comments）にはモックファイルが用意されている。

### 6.6 サーバー起動・停止

#### 起動シーケンス

1. `.env` ファイルの読み込み（存在しなくてもエラーにならない）
2. ログ設定の初期化
3. Supabaseクライアントの初期化（コネクションプール設定）
4. テストクエリ実行（`SELECT 1` で接続確認）
5. Echoインスタンス生成
6. ミドルウェア設定（Logger, Recover, CORS）
7. ルーティング設定（DI含む）
8. サーバー起動（デフォルトポート: 8080）

#### 停止シーケンス（グレースフルシャットダウン）

1. SIGINT/SIGTERM シグナルの受信
2. Echoサーバーのクローズ
3. Supabaseコネクションプールのクローズ
