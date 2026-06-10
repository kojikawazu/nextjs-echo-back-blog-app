# API仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo）の全APIエンドポイントを定義するAPI仕様書です。コードベースのリバースエンジニアリングに基づいて作成しています。

## 目次

- [共通仕様](#共通仕様)
  - [ベースURL](#ベースurl)
  - [認証方式](#認証方式)
  - [共通エラーレスポンス形式](#共通エラーレスポンス形式)
  - [共通レスポンスヘッダー](#共通レスポンスヘッダー)
- [1. ヘルスチェック](#1-ヘルスチェック)
  - [GET /](#get-)
- [2. 認証（Auth）](#2-認証auth)
  - [POST /api/users/login](#post-apiuserslogin)
  - [GET /api/users/auth-check](#get-apiusersauth-check)
  - [POST /api/users/logout](#post-apiuserslogout)
- [3. ブログユーザー（Blog Users）](#3-ブログユーザーblog-users)
  - [GET /api/users/detail](#get-apiusersdetail)
  - [PUT /api/users/update](#put-apiusersupdate)
- [4. ブログ（Blogs）](#4-ブログblogs)
  - [GET /api/blogs](#get-apiblogs)
  - [GET /api/blogs/user/:userId](#get-apiblogsuseruserid)
  - [GET /api/blogs/detail/:id](#get-apiblogsdetailid)
  - [GET /api/blogs/categories](#get-apiblogscategories)
  - [GET /api/blogs/tags](#get-apiblogstags)
  - [GET /api/blogs/popular/:count](#get-apiblogspopularcount)
  - [POST /api/blogs/create](#post-apiblogscreate)
  - [PUT /api/blogs/update/:id](#put-apiblogsupdateid)
  - [DELETE /api/blogs/delete/:id](#delete-apiblogsdeleteid)
- [5. ブログいいね（Blog Likes）](#5-ブログいいねblog-likes)
  - [GET /api/blog-likes](#get-apiblog-likes)
  - [GET /api/blog-likes/generate-visit-id](#get-apiblog-likesgenerate-visit-id)
  - [GET /api/blog-likes/is-liked/:blogId](#get-apiblog-likesis-likedblogid)
  - [POST /api/blog-likes/create/:blogId](#post-apiblog-likescreateblogid)
  - [DELETE /api/blog-likes/delete/:blogId](#delete-apiblog-likesdeleteblogid)
- [6. コメント（Comments）](#6-コメントcomments)
  - [GET /api/comments/blog/:blogId](#get-apicommentsblogblogid)
  - [POST /api/comments/create](#post-apicommentscreate)
- [7. データモデル](#7-データモデル)
  - [BlogData](#blogdata)
  - [BlogUsersData](#blogusersdata)
  - [BlogLikesData](#bloglikesdata)
  - [BlogCommentsData](#blogcommentsdata)
  - [Claims（JWTペイロード）](#claimsjwtペイロード)
  - [ClaimsVisitId（訪問者JWTペイロード）](#claimsvisitid訪問者jwtペイロード)

---

## 共通仕様

### ベースURL

```
/api
```

### 認証方式

- JWT（JSON Web Token）をHTTP-Only Cookieで管理
- Cookie名: `token`（認証用）、`visit-id-token`（訪問者ID用）
- トークン有効期限: 1時間
- 署名アルゴリズム: HS256

### 共通エラーレスポンス形式

```json
{
  "error": "エラーメッセージ"
}
```

### 共通レスポンスヘッダー

- Content-Type: `application/json`
- CORS: `ALLOWED_ORIGINS` 環境変数で制御

---

## 1. ヘルスチェック

### GET /

サービスの稼働状態を確認するエンドポイント。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | text/plain |

#### レスポンス

| ステータス | 説明 |
|-----------|------|
| 200 OK | サービス稼働中 |

```
Service is running
```

---

## 2. 認証（Auth）

### POST /api/users/login

メールアドレスとパスワードでログインし、JWTトークンをCookieに設定する。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### リクエストボディ

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| email | string | はい | メールアドレス（有効な形式であること） |
| password | string | はい | パスワード |

#### レスポンス

**成功時（200 OK）**

```json
{
  "message": "Login successful"
}
```

- `Set-Cookie: token=<JWT>; HttpOnly; Path=/; Expires=<1時間後>`

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `Email and password are required` | email または password が空 |
| 400 Bad Request | `Invalid email format` | メールアドレスの形式が不正 |
| 401 Unauthorized | `Invalid credentials` | `FetchUserByEmailAndPassword` がエラーを返した場合（ユーザー未存在、DB取得失敗を含む全てのエラー） |
| 500 Internal Server Error | `An error occurred` | AuthService.Loginバリデーションで未定義のエラーが発生（現行実装では到達しない） |
| 500 Internal Server Error | `Could not create token` | JWTトークン生成に失敗 |

> **注意**: Handler層（`handlers/auth/auth_impl.go` の `Login`）は `FetchUserByEmailAndPassword` からのエラーを原因に関わらず一律401 `Invalid credentials` として返す。ユーザー未存在もDB接続エラーも区別せず認証失敗（401）に丸める実装のため、DB障害時にも401が返される。

---

### GET /api/users/auth-check

Cookieに含まれるJWTトークンの有効性を検証する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | application/json |

#### リクエストパラメータ

なし（CookieからJWTトークンを取得）

#### レスポンス

**成功時（200 OK）**

```json
{
  "message": "Authenticated",
  "user_id": "uuid-string",
  "username": "ユーザー名",
  "email": "user@example.com"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 401 Unauthorized | `Token not found` | Cookieにトークンが存在しない |
| 401 Unauthorized | `Invalid token` | トークンの解析失敗または無効 |

---

### POST /api/users/logout

JWTトークンのCookieを削除してログアウトする。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### リクエストパラメータ

なし

#### レスポンス

**成功時（200 OK）**

```json
{
  "message": "Logout successful"
}
```

- `Set-Cookie: token=; HttpOnly; Path=/; Expires=<過去の日時>`

---

## 3. ブログユーザー（Blog Users）

### GET /api/users/detail

JWTトークンからユーザーIDを取得し、そのユーザーの詳細情報を返す。パスワードフィールドは空文字に置換される。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | application/json |

#### リクエストパラメータ

なし（CookieのJWTからユーザーIDを抽出）

#### レスポンス

**成功時（200 OK）**

```json
{
  "id": "uuid-string",
  "name": "ユーザー名",
  "email": "user@example.com",
  "password": "",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Id is required` | IDが空 |
| 401 Unauthorized | `Error getting cookie` | Cookieの取得に失敗 |
| 401 Unauthorized | `Error getting userId from token` | トークンからユーザーID取得に失敗 |
| 500 Internal Server Error | `Failed to fetch user` | DB取得エラー（ユーザー未存在の場合も含む） |

> **注意**: Handler層には `404 User not found` の分岐が存在するが、Service層がRepository層のエラーを一律 `"failed to fetch user"` に変換するため、現行実装では404が返ることはない。

---

### PUT /api/users/update

JWTトークンで認証されたユーザーの情報を更新する。更新後、新しいJWTトークンを再発行する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | application/json |

#### リクエストボディ

```json
{
  "name": "新しいユーザー名",
  "email": "new_email@example.com",
  "password": "現在のパスワード",
  "newPassword": "新しいパスワード"
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| name | string | はい | 新しいユーザー名 |
| email | string | はい | 新しいメールアドレス |
| password | string | はい | 現在のパスワード（確認用） |
| newPassword | string | はい | 新しいパスワード |

#### レスポンス

**成功時（200 OK）**

```json
{
  "id": "uuid-string",
  "name": "新しいユーザー名",
  "email": "new_email@example.com",
  "password": "",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z"
}
```

> **注意**: UPDATE SQLのRETURNING句に `password` が含まれていないため、レスポンスの `password` フィールドは常に空文字となる。

- JWTトークンがCookieに再設定される

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `id is required` | IDが空 |
| 400 Bad Request | `Name is required` | name が空 |
| 400 Bad Request | `Email is required` | email が空 |
| 400 Bad Request | `Password is required` | password が空 |
| 400 Bad Request | `New password is required` | newPassword が空 |
| 400 Bad Request | `Invalid email format` | メールアドレスの形式が不正 |
| 401 Unauthorized | `Error getting cookie` | Cookieの取得に失敗 |
| 401 Unauthorized | `Error getting userId from token` | トークンからユーザーID取得に失敗 |
| 500 Internal Server Error | `Failed to update user` | DB更新エラー、ユーザー検証失敗、現パスワード不一致 |
| 500 Internal Server Error | `Failed to create token` | トークン再作成に失敗 |

> **注意**: Handler層には `404 User not found` の分岐が存在するが、Service層は `"failed to validate user"` や `"invalid current password"` を返すため、これらのエラーはHandler側の `"user not found"` 条件に一致せず、全てdefaultの500に落ちる。

---

## 4. ブログ（Blogs）

### GET /api/blogs

全ブログデータを取得する。いいね数とコメント数を集計して含む。作成日時の降順でソート。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### レスポンス

**成功時（200 OK）**

```json
[
  {
    "id": "uuid-string",
    "blog_user_id": "uuid-string",
    "title": "ブログタイトル",
    "description": "ブログの説明",
    "github_url": "https://github.com/user/repo",
    "category": "カテゴリ名",
    "tags": "tag1,tag2,tag3",
    "likes": 5,
    "comment_cnt": 3,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Error fetching blogs` | DB取得エラー |

---

### GET /api/blogs/user/:userId

指定されたユーザーIDに紐づくブログデータを取得する。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| userId | string (UUID) | はい | ブログユーザーID |

#### レスポンス

**成功時（200 OK）**

```json
[
  {
    "id": "uuid-string",
    "blog_user_id": "uuid-string",
    "title": "ブログタイトル",
    "description": "ブログの説明",
    "github_url": "https://github.com/user/repo",
    "category": "カテゴリ名",
    "tags": "tag1,tag2",
    "likes": 2,
    "comment_cnt": 1,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid userId` | userId が空 |
| 500 Internal Server Error | `Error fetching blogs` | DB取得エラー |

> **注意**: 該当ユーザーのブログが0件の場合、エラーではなく `200 OK` で空配列 `[]` が返される。Handler層には `404 Blog not found` の分岐が存在するが、Service層が返すエラー文字列 `"blogs not found"` と一致しないため（単数形/複数形の不一致）、この分岐は到達不能である。

---

### GET /api/blogs/detail/:id

指定されたIDのブログ詳細データを取得する。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| id | string (UUID) | はい | ブログID |

#### レスポンス

**成功時（200 OK）**

```json
{
  "id": "uuid-string",
  "blog_user_id": "uuid-string",
  "title": "ブログタイトル",
  "description": "ブログの説明",
  "github_url": "https://github.com/user/repo",
  "category": "カテゴリ名",
  "tags": "tag1,tag2",
  "likes": 5,
  "comment_cnt": 3,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid id` | id が空 |
| 404 Not Found | `Blog not found` | 該当ブログが存在しない |
| 500 Internal Server Error | `Error fetching blog` | DB取得エラー |

---

### GET /api/blogs/categories

登録されている全ブログのカテゴリ一覧を取得する（DISTINCT）。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### レスポンス

**成功時（200 OK）**

```json
["カテゴリ1", "カテゴリ2", "カテゴリ3"]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Error fetching categories` | DB取得エラー |

---

### GET /api/blogs/tags

登録されている全ブログのタグ一覧を取得する。カンマ区切りのタグを分割し、重複を除去してアルファベット順にソートして返す。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### レスポンス

**成功時（200 OK）**

```json
["Go", "React", "TypeScript", "Web"]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Error fetching tags` | DB取得エラー |

---

### GET /api/blogs/popular/:count

いいね数の多い順にブログを取得する（上位N件）。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| count | integer | はい | 取得件数 |

#### レスポンス

**成功時（200 OK）**

```json
[
  {
    "id": "uuid-string",
    "blog_user_id": "uuid-string",
    "title": "人気ブログタイトル",
    "description": "",
    "github_url": "",
    "category": "",
    "tags": "",
    "likes": 10,
    "comment_cnt": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

注意: popular エンドポイントでは `description`、`github_url`、`category`、`tags` は空文字で返される。

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid count` | count がint変換不可（非数値文字列） |
| 404 Not Found | `blog not found` | Service層が `"blog not found"` を含むエラーを返した場合（Handlerが `strings.Contains` で判定） |
| 500 Internal Server Error | `Error fetching popular blogs` | 上記以外のDB取得エラー、またはcount<=0（Service層が `"invalid count"` を返すが "blog not found" を含まないためdefault分岐で500となる） |

> **注意**: 該当ブログが0件の場合、エラーではなく `200 OK` で空配列 `[]` が返される。count<=0のバリデーションエラーはService層で検出されるが、Handler層にこのケース用の400分岐がないため500として返される。

---

### POST /api/blogs/create

新しいブログ記事を作成する。認証が必要。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | application/json |

#### リクエストボディ

```json
{
  "title": "ブログタイトル",
  "githubUrl": "https://github.com/user/repo",
  "category": "カテゴリ名",
  "description": "ブログの説明",
  "tags": "tag1,tag2,tag3"
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| title | string | はい | ブログタイトル |
| githubUrl | string | はい | GitHubリポジトリURL |
| category | string | はい | カテゴリ |
| description | string | はい | ブログの説明 |
| tags | string | はい | タグ（カンマ区切り） |

#### レスポンス

**成功時（201 Created）**

```json
{
  "id": "uuid-string",
  "blog_user_id": "uuid-string",
  "title": "ブログタイトル",
  "description": "ブログの説明",
  "github_url": "https://github.com/user/repo",
  "category": "カテゴリ名",
  "tags": "tag1,tag2,tag3",
  "likes": 0,
  "comment_cnt": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `Invalid userId` | userId が空 |
| 400 Bad Request | `Invalid title` | title が空 |
| 400 Bad Request | `Invalid githubUrl` | githubUrl が空 |
| 400 Bad Request | `Invalid category` | category が空 |
| 400 Bad Request | `Invalid description` | description が空 |
| 400 Bad Request | `Invalid tags` | tags が空 |
| 401 Unauthorized | `Error getting cookie` | Cookieの取得に失敗 |
| 401 Unauthorized | `Error getting userId from token` | トークンからユーザーID取得に失敗 |
| 500 Internal Server Error | `Failed to create blog` | DB作成エラー |

---

### PUT /api/blogs/update/:id

指定されたIDのブログ記事を更新する。認証が必要。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| id | string (UUID) | はい | 更新対象のブログID |

#### リクエストボディ

```json
{
  "title": "更新後タイトル",
  "githubUrl": "https://github.com/user/repo",
  "category": "更新後カテゴリ",
  "description": "更新後の説明",
  "tags": "tag1,tag2"
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| title | string | はい | ブログタイトル |
| githubUrl | string | はい | GitHubリポジトリURL |
| category | string | はい | カテゴリ |
| description | string | はい | ブログの説明 |
| tags | string | はい | タグ（カンマ区切り） |

#### レスポンス

**成功時（200 OK）**

```json
{
  "id": "uuid-string",
  "blog_user_id": "uuid-string",
  "title": "更新後タイトル",
  "description": "更新後の説明",
  "github_url": "https://github.com/user/repo",
  "category": "更新後カテゴリ",
  "tags": "tag1,tag2",
  "likes": 5,
  "comment_cnt": 3,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `Invalid id` | id が空 |
| 400 Bad Request | `Invalid title` | title が空 |
| 400 Bad Request | `Invalid githubUrl` | githubUrl が空 |
| 400 Bad Request | `Invalid category` | category が空 |
| 400 Bad Request | `Invalid description` | description が空 |
| 400 Bad Request | `Invalid tags` | tags が空 |
| 401 Unauthorized | `Error getting cookie` | Cookieの取得に失敗 |
| 401 Unauthorized | `Error getting userId from token` | トークンからユーザーID取得に失敗 |
| 500 Internal Server Error | `Failed to update blog` | DB更新エラー |

---

### DELETE /api/blogs/delete/:id

指定されたIDのブログ記事を削除する。認証が必要。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（token Cookie） |
| Content-Type | なし（No Content） |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| id | string (UUID) | はい | 削除対象のブログID |

#### レスポンス

**成功時（204 No Content）**

レスポンスボディなし。

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid id` | id が空文字 |
| 401 Unauthorized | `Error getting cookie` | Cookieの取得に失敗 |
| 401 Unauthorized | `Error getting userId from token` | トークンからユーザーID取得に失敗 |
| 500 Internal Server Error | `Failed to delete blog` | DB削除エラー、または非空の不正UUID形式 |

> **注意**: Service層は `id == ""` のみを400として検出する。非空だがUUID形式でない文字列は、Service層のバリデーションを通過しRepository層の `uuid.Parse` で弾かれるが、Service層がそのエラーを `"failed to delete blog"` に変換するため、Handlerでは500として返される。

---

## 5. ブログいいね（Blog Likes）

### GET /api/blog-likes

訪問者IDに紐づくいいねデータ一覧を取得する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（visit-id-token Cookie） |
| Content-Type | application/json |
| Cache-Control | no-store |

#### レスポンス

**成功時（200 OK）**

```json
[
  {
    "id": "uuid-string",
    "blog_id": "uuid-string",
    "visit_id": "uuid-string",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Failed to get visit id token` | visit-id-token Cookieの取得に失敗 |
| 500 Internal Server Error | `Failed to get visit id` | 訪問者IDの取得に失敗 |
| 500 Internal Server Error | `Error fetching blog likes by visit id` | DB取得エラー |

---

### GET /api/blog-likes/generate-visit-id

訪問者IDを生成してCookieに設定する。既にCookieに有効な訪問者IDが存在する場合はスキップする。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### レスポンス

**成功時（200 OK）** - 新規生成

```json
{
  "message": "Visitor id generated successfully"
}
```

- `Set-Cookie: visit-id-token=<JWT>; HttpOnly; Path=/; Expires=<1時間後>`

**成功時（200 OK）** - 既存

```json
{
  "message": "Visitor id already exists"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Failed to create visitor token` | トークン生成に失敗 |

---

### GET /api/blog-likes/is-liked/:blogId

指定したブログに対して訪問者がいいねしているかを確認する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（visit-id-token Cookie） |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| blogId | string (UUID) | はい | ブログID |

#### レスポンス

**成功時（200 OK）**

```json
{
  "isLiked": true
}
```

または

```json
{
  "isLiked": false
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Failed to get visit id token` | visit-id-token Cookieの取得に失敗 |
| 500 Internal Server Error | `Failed to get visit id` | 訪問者IDの取得に失敗 |

> **注意**: いいね判定（`IsBlogLiked`）でService層がエラーを返した場合も、500にはならず `200 OK` で `{"isLiked": false}` を返す（いいね未登録を正常系として扱う実装のため、DB障害も同様に握り潰される）。500はvisit-id-token Cookieの取得・トークン解析に失敗した場合のみ。

---

### POST /api/blog-likes/create/:blogId

指定したブログにいいねを追加する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（visit-id-token Cookie） |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| blogId | string (UUID) | はい | ブログID |

#### レスポンス

**成功時（200 OK）**

```json
{
  "id": "uuid-string",
  "blog_id": "uuid-string",
  "visit_id": "uuid-string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `BlogId or VisitId is empty` | blogId または visitId が空 |
| 400 Bad Request | `Blog is already liked` | 既にいいね済み |
| 500 Internal Server Error | `Failed to get visit id token` | visit-id-token Cookieの取得に失敗 |
| 500 Internal Server Error | `Failed to get visit id` | 訪問者IDの取得に失敗 |
| 500 Internal Server Error | `Error creating blog like` | DB作成エラー |

---

### DELETE /api/blog-likes/delete/:blogId

指定したブログのいいねを削除する。

| 項目 | 内容 |
|------|------|
| 認証 | 必要（visit-id-token Cookie） |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| blogId | string (UUID) | はい | ブログID |

#### レスポンス

**成功時（200 OK）**

```json
{
  "message": "Blog like deleted successfully"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 500 Internal Server Error | `Failed to get visit id token` | visit-id-token Cookieの取得に失敗 |
| 500 Internal Server Error | `Failed to get visit id` | 訪問者IDの取得に失敗 |
| 500 Internal Server Error | `Error deleting blog like` | DB削除エラー |

---

## 6. コメント（Comments）

### GET /api/comments/blog/:blogId

指定したブログIDに紐づくコメント一覧を取得する。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### パスパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|------|------|------|
| blogId | string (UUID) | はい | ブログID |

#### レスポンス

**成功時（200 OK）**

```json
[
  {
    "id": "uuid-string",
    "blog_id": "uuid-string",
    "guest_user": "ゲストユーザー名",
    "comment": "コメント内容",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Invalid blogId` | blogId が空 |
| 404 Not Found | `Comments not found` | DBクエリエラー発生時（Service層がRepository層のエラーを一律 `"comments not found"` に変換するため、Handlerの404分岐に到達する） |

> **注意**: 該当ブログにコメントが0件の場合、エラーではなく `200 OK` で空配列 `[]` が返される。DBクエリ自体が失敗した場合は、Service層が `"comments not found"` を返すためHandlerの404分岐に到達し、500ではなく404が返る。Handlerのdefault 500分岐（`Error fetching comments`）は、現行のService実装では到達しない。

---

### POST /api/comments/create

新しいコメントを作成する。

| 項目 | 内容 |
|------|------|
| 認証 | 不要 |
| Content-Type | application/json |

#### リクエストボディ

```json
{
  "blogId": "uuid-string",
  "guestUser": "ゲストユーザー名",
  "comment": "コメント内容"
}
```

| フィールド | 型 | 必須 | 説明 |
|-----------|------|------|------|
| blogId | string (UUID) | はい | コメント対象のブログID |
| guestUser | string | はい | ゲストユーザー名 |
| comment | string | はい | コメント内容 |

#### レスポンス

**成功時（201 Created）**

```json
{
  "id": "uuid-string",
  "blog_id": "uuid-string",
  "guest_user": "ゲストユーザー名",
  "comment": "コメント内容",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**エラー時**

| ステータス | エラーメッセージ | 条件 |
|-----------|----------------|------|
| 400 Bad Request | `Error binding request` | リクエストボディのパースに失敗 |
| 400 Bad Request | `Invalid blogId` | blogId が空 |
| 400 Bad Request | `Invalid guestUser` | guestUser が空 |
| 400 Bad Request | `Invalid comment` | comment が空 |
| 500 Internal Server Error | `Failed to create comment` | DB作成エラー |

---

## 7. データモデル

### BlogData

| フィールド | JSON名 | DB名 | 型 | 説明 |
|-----------|--------|------|----|------|
| ID | id | id | string (UUID) | ブログID |
| BlogUserId | blog_user_id | blog_user_id | string (UUID) | ブログユーザーID |
| Title | title | title | string | タイトル |
| Description | description | description | string | 説明 |
| GithubUrl | github_url | github_url | string | GitHubリポジトリURL |
| Category | category | category | string | カテゴリ |
| Tags | tags | tags | string | タグ（カンマ区切り） |
| Likes | likes | likes | int8 | いいね数（集計値） |
| CommentCnt | comment_cnt | comment_cnt | int8 | コメント数（集計値） |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |
| UpdatedAt | updated_at | updated_at | time.Time | 更新日時 |

### BlogUsersData

| フィールド | JSON名 | DB名 | 型 | 説明 |
|-----------|--------|------|----|------|
| ID | id | id | string (UUID) | ユーザーID |
| Name | name | name | string | ユーザー名 |
| Email | email | email | string | メールアドレス |
| Password | password | password | string | パスワード |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |
| UpdatedAt | updated_at | updated_at | time.Time | 更新日時 |

### BlogLikesData

| フィールド | JSON名 | DB名 | 型 | 説明 |
|-----------|--------|------|----|------|
| ID | id | id | string (UUID) | いいねID |
| BlogId | blog_id | blog_id | string (UUID) | ブログID |
| VisitId | visit_id | visit_id | string (UUID) | 訪問者ID |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |
| UpdatedAt | updated_at | updated_at | time.Time | 更新日時 |

### BlogCommentsData

| フィールド | JSON名 | DB名 | 型 | 説明 |
|-----------|--------|------|----|------|
| ID | id | id | string (UUID) | コメントID |
| BlogId | blog_id | blog_id | string (UUID) | ブログID |
| GuestUser | guest_user | guest_user | string | ゲストユーザー名 |
| Comment | comment | comment | string | コメント内容 |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |

### Claims（JWTペイロード）

| フィールド | JSON名 | 型 | 説明 |
|-----------|--------|------|------|
| UserID | user_id | string | ユーザーID |
| Email | email | string | メールアドレス |
| Username | username | string | ユーザー名 |

### ClaimsVisitId（訪問者JWTペイロード）

| フィールド | JSON名 | 型 | 説明 |
|-----------|--------|------|------|
| VisitId | visit_id | string | 訪問者ID |
