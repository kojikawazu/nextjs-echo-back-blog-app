# API仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo）の全APIエンドポイントを定義するAPI仕様書です。コードベースのリバースエンジニアリングに基づいて作成しています。

## 目次

- [共通仕様](#共通仕様)
  - [ベースURL](#ベースurl)
  - [認証方式](#認証方式)
  - [共通エラーレスポンス形式](#共通エラーレスポンス形式)
  - [レートリミット](#レートリミット)
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

エラーレスポンスは全エンドポイントで次の形式に統一する（`dto.ErrorResponse`、#123）。

```json
{
  "code": "BLOG_NOT_FOUND",
  "error": "Blog not found"
}
```

| フィールド | 型 | 説明 |
|-----------|------|------|
| code | string | クライアントが機械的に判定するための安定した識別子（一度公開したら変更しない） |
| error | string | 人が読むためのメッセージ |

- ハンドラは HTTP ステータスや本文を組み立てず、エラー（`apperrors` の sentinel）を返す。ステータス・本文への変換は `middlewares.HTTPErrorHandler` に集約している。サービス層のエラーは文字列比較ではなく `errors.Is` で判定する。
- **500 は内部情報を出さない**。原因に関わらず常に `{"code": "INTERNAL_ERROR", "error": "Internal server error"}` を返し、エラーの詳細はアクセスログ（`error` 属性）にのみ記録する。
- Echo 自身が返すエラー（存在しないルート・許可されていないメソッド等）も同じ形式で返す。コードは HTTP ステータス名から生成する（例: `NOT_FOUND`・`METHOD_NOT_ALLOWED`）。メッセージはステータス名（`Not Found` 等）で、Echo の内部メッセージは返さない。
- フロントエンド（BFF）はエラー本文を解釈せず、ステータスコード（`response.ok` / `status`）で判定している（#123 の調査結果）。本文の形式変更はフロントエンドに影響しないが、ステータスコードの変更は影響する。

#### エラーコード一覧

| コード | ステータス | 意味 |
|--------|-----------|------|
| `VALIDATION_ERROR` | 400 | 入力値の不正（どの項目かはメッセージで示す。例: `Invalid title`） |
| `INVALID_REQUEST_BODY` | 400 | リクエストボディのパース（バインド）に失敗 |
| `INVALID_CURRENT_PASSWORD` | 400 | ユーザー更新時、現在のパスワードが一致しない |
| `AUTH_TOKEN_MISSING` | 401 | 認証トークン（token Cookie）が存在しない |
| `AUTH_TOKEN_INVALID` | 401 | 認証トークンが不正・期限切れ・種類違い（訪問者トークン等） |
| `VISITOR_TOKEN_MISSING` | 401 | 訪問者トークン（visit-id-token Cookie）が存在しない（`GET /api/blog-likes/generate-visit-id` で再発行できる） |
| `VISITOR_TOKEN_INVALID` | 401 | 訪問者トークンが不正・期限切れ・種類違い（認証トークン等）（同上） |
| `INVALID_CREDENTIALS` | 401 | ログイン失敗（ユーザー不在・パスワード不一致を区別しない） |
| `BLOG_NOT_FOUND` | 404 | 該当するブログが存在しない |
| `USER_NOT_FOUND` | 404 | 該当するユーザーが存在しない |
| `ALREADY_LIKED` | 409 | 既にいいね済み |
| `TOO_MANY_REQUESTS` | 429 | レートリミット超過 |
| `INTERNAL_ERROR` | 500 | サーバー内部のエラー（想定外のエラーを含む） |
| `NOT_FOUND` / `METHOD_NOT_ALLOWED` 等 | 404 / 405 等 | Echo 由来のエラー（ステータス名から生成） |

### レートリミット

クライアント IP 単位で上限を設け、超過時は **429 Too Many Requests** を返す（#124）。

```json
{
  "code": "TOO_MANY_REQUESTS",
  "error": "Too many requests"
}
```

| 対象 | 上限 |
|------|------|
| 全エンドポイント（`GET /` を除く） | 毎秒 20 件・瞬間 60 件まで |
| `POST /api/users/login`・`PUT /api/users/update`（2 つで共有） | 5 回まで、以降 12 秒ごとに 1 回（5 回/分） |

- クライアント IP は、BFF が付与する `X-Client-IP` を使う（本番のみ。`TRUST_CLIENT_IP_HEADER=true`）。ヘッダがない・解釈できない場合や、ローカル・テストでは `X-Forwarded-For` の右端（Cloud Run のフロントエンドが付与した値）から取る。詳細は [06-security-specification.md](06-security-specification.md) 8.2（#146）。

### 呼び出し元の制限

- 本番のバックエンドは BFF（フロントエンドの Cloud Run のサービスアカウント）からのみ呼べる（Cloud Run IAM。#146）。BFF は ID トークンを `Authorization: Bearer` で付ける。ID トークンのないリクエストは、アプリに届く前に Cloud Run が 401 で拒否する（このエラーは統一エラーレスポンスの形式ではない）。
- BFF が付与するヘッダ:

| ヘッダ | 内容 |
|--------|------|
| `Authorization` | `Bearer <ID トークン>`（audience はバックエンドの URL）。Cloud Run が検証する |
| `X-Client-IP` | クライアント IP（レートリミットのカウントに使う） |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `INVALID_REQUEST_BODY` | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `VALIDATION_ERROR` | `Email and password are required` | email または password が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid email format` | メールアドレスの形式が不正 |
| 401 Unauthorized | `INVALID_CREDENTIALS` | `Invalid credentials` | ユーザーが存在しない、またはパスワードが一致しない（どちらかは区別しない） |
| 429 Too Many Requests | `TOO_MANY_REQUESTS` | `Too many requests` | 認証系レートリミット超過（`PUT /api/users/update` と共有。5 回/分） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー、JWT トークン生成に失敗 |

> DB 障害は認証失敗（401）に丸めず 500 を返す（#123 で変更。以前は原因に関わらず 401 だった）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |

> #123 以前は本文のキーが `message`（`{"message": "Token not found"}` 等）だったが、共通エラーレスポンス形式に統一した。

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
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

> レスポンスはレスポンス DTO（`dto.UserResponse`）で返し、`password` キーは含まない（#121）。

**エラー時**

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Id is required` | トークンから取得したユーザー ID が空 |
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |
| 404 Not Found | `USER_NOT_FOUND` | `User not found` | トークンのユーザーが存在しない（#123 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

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
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T00:00:00Z"
}
```

> レスポンスはレスポンス DTO（`dto.UserResponse`）で返し、`password` キーは含まない（#121）。

- JWTトークンがCookieに再設定される

**エラー時**

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `INVALID_REQUEST_BODY` | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `VALIDATION_ERROR` | `Id is required` | トークンから取得したユーザー ID が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Name is required` | name が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Email is required` | email が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Password is required` | password が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `New password is required` | newPassword が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid email format` | メールアドレスの形式が不正 |
| 400 Bad Request | `VALIDATION_ERROR` | `Password is too long` | newPassword が 72 バイトを超える（bcrypt の上限） |
| 400 Bad Request | `INVALID_CURRENT_PASSWORD` | `Invalid current password` | 現在のパスワードが一致しない（#123 で 500 から変更） |
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |
| 404 Not Found | `USER_NOT_FOUND` | `User not found` | トークンのユーザーが存在しない（#123 で 500 から変更） |
| 429 Too Many Requests | `TOO_MANY_REQUESTS` | `Too many requests` | 認証系レートリミット超過（`POST /api/users/login` と共有。5 回/分） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得・更新エラー、パスワードのハッシュ化失敗、トークン再作成に失敗 |

> 新しいパスワードは bcrypt でハッシュ化して保存する。現在のパスワードは保存済みのハッシュと照合する（#120）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid userId` | userId が空、または UUID 形式でない（#123 で UUID 形式の検証を追加。以前は DB エラーとして 500） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

> **注意**: 該当ユーザーのブログが 0 件の場合はエラーではなく `200 OK` を返す（404 にはしない。本文は「7. データモデル」の 0 件時の扱いを参照）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid id` | id が空、または UUID 形式でない（#123 で 404 から変更） |
| 404 Not Found | `BLOG_NOT_FOUND` | `Blog not found` | 該当ブログが存在しない |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー（#123 で 404 から変更。以前は DB 障害も 404 だった） |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid count` | count が数値でない、または 0 以下（0 以下は #123 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

> **注意**: 該当ブログが 0 件の場合はエラーではなく `200 OK` を返す（404 にはしない）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `INVALID_REQUEST_BODY` | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid userId` | トークンから取得したユーザー ID が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid title` | title が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid githubUrl` | githubUrl が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid category` | category が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid description` | description が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid tags` | tags が空 |
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 作成エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `INVALID_REQUEST_BODY` | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid id` | id が空、または UUID 形式でない（#123 で UUID 形式の検証を追加。以前は DB エラーとして 500） |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid title` | title が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid githubUrl` | githubUrl が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid category` | category が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid description` | description が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid tags` | tags が空 |
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |
| 404 Not Found | `BLOG_NOT_FOUND` | `Blog not found` | 該当ブログが存在しない（#123 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 更新エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid id` | id が空、または UUID 形式でない（#123 で UUID 形式の検証を Service 層に移動。以前は 500） |
| 401 Unauthorized | `AUTH_TOKEN_MISSING` | `Authentication token is missing` | token Cookie が存在しない |
| 401 Unauthorized | `AUTH_TOKEN_INVALID` | `Authentication token is invalid` | トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（訪問者トークン等。#156）・`user_id` が空 |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 削除エラー |

> **注意**: 存在しない（UUID 形式は正しい）ID の削除はエラーにせず `204 No Content` を返す（冪等な削除として扱う）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `VisitId is empty` | トークンから取得した訪問者 ID が空 |
| 401 Unauthorized | `VISITOR_TOKEN_MISSING` | `Visitor token is missing` | visit-id-token Cookie が存在しない（#148 で 500 から変更） |
| 401 Unauthorized | `VISITOR_TOKEN_INVALID` | `Visitor token is invalid` | 訪問者トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（認証トークン等。#156）・`visit_id` が空（#148 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | トークン生成に失敗 |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `BlogId or VisitId is empty` | blogId または訪問者 ID が空 |
| 401 Unauthorized | `VISITOR_TOKEN_MISSING` | `Visitor token is missing` | visit-id-token Cookie が存在しない（#148 で 500 から変更） |
| 401 Unauthorized | `VISITOR_TOKEN_INVALID` | `Visitor token is invalid` | 訪問者トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（認証トークン等。#156）・`visit_id` が空（#148 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 確認エラー（#123 で変更。以前は DB 障害も `200 OK` の `{"isLiked": false}` に握り潰していた） |

> いいねが存在しない場合はエラーではなく `200 OK` で `{"isLiked": false}` を返す（Repository 層が「該当なし」をエラーではなく `false` として返す）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `BlogId or VisitId is empty` | blogId または訪問者 ID が空（#123 で 500 から変更） |
| 401 Unauthorized | `VISITOR_TOKEN_MISSING` | `Visitor token is missing` | visit-id-token Cookie が存在しない（#148 で 500 から変更） |
| 401 Unauthorized | `VISITOR_TOKEN_INVALID` | `Visitor token is invalid` | 訪問者トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（認証トークン等。#156）・`visit_id` が空（#148 で 500 から変更） |
| 409 Conflict | `ALREADY_LIKED` | `Blog is already liked` | 既にいいね済み（#123 で 400 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | いいね済みかの確認、または DB 作成エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 401 Unauthorized | `VISITOR_TOKEN_MISSING` | `Visitor token is missing` | visit-id-token Cookie が存在しない（#148 で 500 から変更） |
| 401 Unauthorized | `VISITOR_TOKEN_INVALID` | `Visitor token is invalid` | 訪問者トークンの解析失敗・署名不正・期限切れ・`exp` なし・種類違い（認証トークン等。#156）・`visit_id` が空（#148 で 500 から変更） |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 削除エラー |

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid blogId` | blogId が空 |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 取得エラー（#123 で 404 から変更。以前は DB 障害を 404 `Comments not found` で返していた） |

> **注意**: 該当ブログにコメントが 0 件の場合はエラーではなく `200 OK` を返す（404 にはしない）。

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

| ステータス | コード | エラーメッセージ | 条件 |
|-----------|--------|----------------|------|
| 400 Bad Request | `INVALID_REQUEST_BODY` | `Invalid request body` | リクエストボディのパースに失敗 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid blogId` | blogId が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid guestUser` | guestUser が空 |
| 400 Bad Request | `VALIDATION_ERROR` | `Invalid comment` | comment が空 |
| 500 Internal Server Error | `INTERNAL_ERROR` | `Internal server error` | DB 作成エラー |

---

## 7. データモデル

API のレスポンスはすべてレスポンス DTO（`backend/dto/`）で返し、モデル（`backend/models/`）を直接返さない（`.claude/rules/api.md`）。モデルは JSON タグを持たない。

- 一覧のレスポンスは、該当が 0 件の場合 `[]` ではなく **`null`** を返す（DTO 化以前からの挙動。#139 で互換性のため維持）
- `GET /api/blogs/categories`・`GET /api/blogs/tags` は文字列の配列（`[]string`）をそのまま返す

### BlogResponse（`dto.BlogResponse`）

ブログ記事のレスポンス DTO。

| フィールド | JSON名 | 元のモデル項目（DB名） | 型 | 説明 |
|-----------|--------|----------------------|----|------|
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

### UserResponse（`dto.UserResponse`）

ユーザー情報のレスポンス DTO。`models.BlogUsersData` から公開してよいフィールドのみを写す。パスワードは含めない。

| フィールド | JSON名 | 元のモデル項目（DB名） | 型 | 説明 |
|-----------|--------|----------------------|----|------|
| ID | id | id | string (UUID) | ユーザーID |
| Name | name | name | string | ユーザー名 |
| Email | email | email | string | メールアドレス |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |
| UpdatedAt | updated_at | updated_at | time.Time | 更新日時 |

> `models.BlogUsersData` は API に直接返さない。JSON タグを持たず、`Password` は `json:"-"` のため誤ってシリアライズされても出力されない。

### BlogLikeResponse（`dto.BlogLikeResponse`）

いいねのレスポンス DTO。フロントエンドの `BlogLike` 型が `visit_id` を参照するため含める。

| フィールド | JSON名 | 元のモデル項目（DB名） | 型 | 説明 |
|-----------|--------|----------------------|----|------|
| ID | id | id | string (UUID) | いいねID |
| BlogId | blog_id | blog_id | string (UUID) | ブログID |
| VisitId | visit_id | visit_id | string (UUID) | 訪問者ID |
| CreatedAt | created_at | created_at | time.Time | 作成日時 |
| UpdatedAt | updated_at | updated_at | time.Time | 更新日時 |

### CommentResponse（`dto.CommentResponse`）

コメントのレスポンス DTO。

| フィールド | JSON名 | 元のモデル項目（DB名） | 型 | 説明 |
|-----------|--------|----------------------|----|------|
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
