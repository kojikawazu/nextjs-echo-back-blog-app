# データ仕様書

本ドキュメントは、ブログWebアプリケーションバックエンドのデータ構造、テーブル定義、リレーション、およびデータフローを、実装コードから逆引きして整理したものである。

## 目次

- [1. ER図（エンティティ関連図）](#1-er図エンティティ関連図)
  - [リレーション一覧](#リレーション一覧)
- [2. テーブル定義](#2-テーブル定義)
  - [2.1 blog_users（ブログユーザー）](#21-blog_usersブログユーザー)
  - [2.2 blogs（ブログ記事）](#22-blogsブログ記事)
  - [2.3 blog_likes（ブログいいね）](#23-blog_likesブログいいね)
  - [2.4 blog_comments（ブログコメント）](#24-blog_commentsブログコメント)
- [3. JWTペイロード構造](#3-jwtペイロード構造)
  - [3.1 認証トークン（Claims）](#31-認証トークンclaims)
  - [3.2 訪問者IDトークン（ClaimsVisitId）](#32-訪問者idトークンclaimsvisitid)
- [4. データフロー](#4-データフロー)
  - [4.1 標準的なリクエスト処理フロー](#41-標準的なリクエスト処理フロー)
  - [4.2 主要クエリパターン](#42-主要クエリパターン)
- [5. UUID使用方針](#5-uuid使用方針)
- [6. Supabase接続設定](#6-supabase接続設定)
  - [6.1 接続パラメータ](#61-接続パラメータ)
  - [6.2 接続ライフサイクル](#62-接続ライフサイクル)
  - [6.3 本番環境での秘密情報管理](#63-本番環境での秘密情報管理)
- [7. 監査列方針](#7-監査列方針)
  - [7.1 現状との乖離（未対応）](#71-現状との乖離未対応)

---

## 1. ER図（エンティティ関連図）

以下にテキストベースのER図を示す。

```
+------------------+       +------------------+
|   blog_users     |       |     blogs        |
+------------------+       +------------------+
| id (PK, UUID)   |──1:N──| id (PK, UUID)    |
| name             |       | blog_user_id (FK)|
| email            |       | title            |
| password         |       | description      |
| created_at       |       | github_url       |
| updated_at       |       | category         |
+------------------+       | tags             |
                           | likes            |
                           | comment_cnt      |
                           | created_at       |
                           | updated_at       |
                           +------------------+
                                  |
                    +-------------+-------------+
                    |                           |
               1:N  |                      1:N  |
                    v                           v
          +------------------+       +------------------+
          |   blog_likes     |       |  blog_comments   |
          +------------------+       +------------------+
          | id (PK, UUID)   |       | id (PK, UUID)    |
          | blog_id (FK)    |       | blog_id (FK)     |
          | visit_id        |       | guest_user       |
          | created_at      |       | comment          |
          | updated_at      |       | created_at       |
          +------------------+       +------------------+
```

### リレーション一覧

| 親テーブル | 子テーブル | カーディナリティ | 外部キー |
|---|---|---|---|
| blog_users | blogs | 1:N | blogs.blog_user_id → blog_users.id |
| blogs | blog_likes | 1:N | blog_likes.blog_id → blogs.id |
| blogs | blog_comments | 1:N | blog_comments.blog_id → blogs.id |

## 2. テーブル定義

### 2.1 blog_users（ブログユーザー）

ブログの管理者ユーザー情報を格納するテーブル。

| カラム名 | データ型 | NULL | 制約 | 説明 |
|---|---|---|---|---|
| id | UUID | NOT NULL | PRIMARY KEY | ユーザーID（自動生成） |
| name | TEXT | NOT NULL | - | ユーザー名 |
| email | TEXT | NOT NULL | - | メールアドレス（ログイン認証に使用） |
| password | TEXT | NOT NULL | - | パスワード（平文で格納） |
| created_at | TIMESTAMP | NOT NULL | DEFAULT now() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | DEFAULT now() | 更新日時 |

対応するGoモデル: `models.BlogUsersData`

```go
type BlogUsersData struct {
    ID        string    `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Email     string    `json:"email" db:"email"`
    Password  string    `json:"password" db:"password"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**備考:**
- ユーザー情報取得時（`FetchBlogUsers`）にはレスポンス返却前にパスワードフィールドを空文字列に上書きしてから返す。
- email/password照合時のSQLクエリでは、パスワードを平文比較している。

### 2.2 blogs（ブログ記事）

ブログ記事の情報を格納するテーブル。

| カラム名 | データ型 | NULL | 制約 | 説明 |
|---|---|---|---|---|
| id | UUID | NOT NULL | PRIMARY KEY | ブログID（自動生成） |
| blog_user_id | UUID | NOT NULL | FOREIGN KEY → blog_users.id | 作成者のユーザーID |
| title | TEXT | NOT NULL | - | ブログタイトル |
| description | TEXT | NOT NULL | - | ブログ本文・説明 |
| github_url | TEXT | NOT NULL | - | 関連するGitHubリポジトリのURL |
| category | TEXT | NOT NULL | - | カテゴリ |
| tags | TEXT | NOT NULL | - | タグ（文字列として格納） |
| likes | INT8 | NOT NULL | DEFAULT 0 | いいね数（集計値を保持するカラム） |
| comment_cnt | INT8 | NOT NULL | DEFAULT 0 | コメント数（集計値を保持するカラム） |
| created_at | TIMESTAMP | NOT NULL | DEFAULT now() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | DEFAULT now() | 更新日時 |

対応するGoモデル: `models.BlogData`

```go
type BlogData struct {
    ID          string    `json:"id" db:"id"`
    BlogUserId  string    `json:"blog_user_id" db:"blog_user_id"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    GithubUrl   string    `json:"github_url" db:"github_url"`
    Category    string    `json:"category" db:"category"`
    Tags        string    `json:"tags" db:"tags"`
    Likes       int8      `json:"likes" db:"likes"`
    CommentCnt  int8      `json:"comment_cnt" db:"comment_cnt"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

**備考:**
- `likes` と `comment_cnt` はテーブル上のカラムとして存在するが、読み取り時には `blog_likes` テーブルおよび `blog_comments` テーブルとのLEFT JOINによるCOUNT集計値が使用される。
- Go側では `int8` 型（-128～127）で保持されるため、いいね数やコメント数が127を超えるとオーバーフローする。
- INSERT時は `RETURNING` 句でテーブル上の `likes`/`comment_cnt` カラム値をそのまま返却し、SELECT時はJOIN集計値を使用する。
- ブログ作成時に `blog_user_id` のUUID形式バリデーションをリポジトリ層で実施する。
- ブログ一覧は `created_at DESC` で降順ソートされる。

### 2.3 blog_likes（ブログいいね）

匿名訪問者によるブログ記事へのいいね情報を格納するテーブル。

| カラム名 | データ型 | NULL | 制約 | 説明 |
|---|---|---|---|---|
| id | UUID | NOT NULL | PRIMARY KEY | いいねID（自動生成） |
| blog_id | UUID | NOT NULL | FOREIGN KEY → blogs.id | 対象ブログのID |
| visit_id | TEXT | NOT NULL | - | 訪問者のID（JWTから抽出されたUUID文字列） |
| created_at | TIMESTAMP | NOT NULL | DEFAULT now() | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | DEFAULT now() | 更新日時 |

対応するGoモデル: `models.BlogLikesData`

```go
type BlogLikesData struct {
    ID        string    `json:"id" db:"id"`
    BlogId    string    `json:"blog_id" db:"blog_id"`
    VisitId   string    `json:"visit_id" db:"visit_id"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**備考:**
- `blog_id` と `visit_id` の組み合わせでいいねの一意性を判定する（同一訪問者が同一ブログに複数回いいねすることを防ぐロジックがService層に存在）。
- 削除時は `blog_id` と `visit_id` の両方を条件として指定する。

### 2.4 blog_comments（ブログコメント）

ブログ記事へのゲストコメント情報を格納するテーブル。

| カラム名 | データ型 | NULL | 制約 | 説明 |
|---|---|---|---|---|
| id | UUID | NOT NULL | PRIMARY KEY | コメントID（自動生成） |
| blog_id | UUID | NOT NULL | FOREIGN KEY → blogs.id | 対象ブログのID |
| guest_user | TEXT | NOT NULL | - | コメント投稿者のゲスト名 |
| comment | TEXT | NOT NULL | - | コメント本文 |
| created_at | TIMESTAMP | NOT NULL | DEFAULT now() | 作成日時 |

対応するGoモデル: `models.BlogCommentsData`

```go
type BlogCommentsData struct {
    ID        string    `json:"id" db:"id"`
    BlogId    string    `json:"blog_id" db:"blog_id"`
    GuestUser string    `json:"guest_user" db:"guest_user"`
    Comment   string    `json:"comment" db:"comment"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

**備考:**
- `blog_likes` と異なり、`updated_at` カラムは存在しない（コメントは作成のみで更新・削除機能なし）。
- 認証不要でゲストユーザーが投稿可能。

## 3. JWTペイロード構造

### 3.1 認証トークン（Claims）

```go
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    jwt.StandardClaims  // ExpiresAt を含む
}
```

Cookie名: `token`

### 3.2 訪問者IDトークン（ClaimsVisitId）

```go
type ClaimsVisitId struct {
    VisitId string `json:"visit_id"`
    jwt.StandardClaims  // ExpiresAt を含む
}
```

Cookie名: `visit-id-token`

## 4. データフロー

### 4.1 標準的なリクエスト処理フロー

```
クライアント
  │
  ├── HTTPリクエスト（JSON body / Path params / Cookies）
  │
  v
Echo Framework
  │
  ├── Middleware（Logger → Recover → CORS）
  │
  v
Handler層（handlers/）
  │
  ├── リクエストバインド（c.Bind, c.Param）
  ├── Cookie認証チェック（書き込み操作時）
  ├── エラーハンドリング → HTTPステータスコード変換
  │
  v
Service層（services/）
  │
  ├── 入力バリデーション
  ├── ビジネスロジック
  │
  v
Repository層（repositories/）
  │
  ├── SQLクエリ構築
  ├── パラメータバインド（$1, $2, ...）
  │
  v
pgxpool（コネクションプール）
  │
  ├── Pool.Query() / Pool.QueryRow() / Pool.Exec()
  │
  v
Supabase（PostgreSQL）
```

### 4.2 主要クエリパターン

| 操作 | メソッド | クエリパターン |
|---|---|---|
| 一覧取得 | `Pool.Query()` | SELECT + LEFT JOIN（集計）+ ORDER BY |
| 単一取得 | `Pool.QueryRow()` | SELECT + WHERE id = $1 |
| 作成 | `Pool.QueryRow()` | INSERT + RETURNING |
| 更新 | `Pool.QueryRow()` | WITH ... UPDATE + RETURNING + LEFT JOIN |
| 削除 | `Pool.Exec()` | DELETE WHERE id = $1 |

- ブログ取得系のクエリでは、`blog_likes` と `blog_comments` を LEFT JOIN + サブクエリで集計し、`COALESCE` でNULLを0に変換する。
- パラメータ化クエリ（`$1`, `$2`, ...）により、SQLインジェクションを防止している。

## 5. UUID使用方針

- すべてのテーブルの主キー（`id`）にはUUID型を採用している。
- `blog_user_id`、`blog_id` などの外部キーもUUID型。
- ブログ作成・削除時にはリポジトリ層で `uuid.Parse()` によるUUID形式のバリデーションを実施する。
- 訪問者ID（`visit_id`）は `uuid.New().String()` で生成されるUUID文字列。
- GoモデルではUUIDを `string` 型として保持する。

## 6. Supabase接続設定

### 6.1 接続パラメータ

| パラメータ | 値/設定元 | 説明 |
|---|---|---|
| 接続URL | 環境変数 `SUPABASE_URL` | PostgreSQL接続文字列 |
| SSLモード | `sslmode=require` | URLに付与される |
| 最大接続数 | 10 | `config.MaxConns` |
| アイドルタイムアウト | 30秒 | `config.MaxConnIdleTime` |
| Simple Protocol | 有効 | `config.ConnConfig.PreferSimpleProtocol = true` |

### 6.2 接続ライフサイクル

1. **初期化**（`InitSupabase()`）: 接続URL解析 → プール設定 → 接続確立 → Ping確認
2. **テストクエリ**（`TestQuery()`）: `SELECT 1` を実行して動作確認
3. **利用**: 各リポジトリが `supabase.Pool` と `supabase.Ctx`（`context.Background()`）を直接参照
4. **終了**（`ClosePool()`）: グレースフルシャットダウン時にプールをクローズ

### 6.3 本番環境での秘密情報管理

- `SUPABASE_URL` はTerraformの `google_secret_manager_secret` リソースとして管理される。
- Cloud Runのコンテナ環境変数として、Secret Managerから自動注入される。

## 7. 監査列方針

本プロジェクトはORMを使わず `pgx` でSQLを直接実行するため、監査列（`created_at` / `updated_at`）の値は**DB側で自動設定する**方針とする。詳細なルールは [`.claude/rules/database.md`](../.claude/rules/database.md) を参照。

| 列 | あるべき型 | あるべき既定 | 用途 |
|---|---|---|---|
| `created_at` | `timestamptz` | `DEFAULT now() NOT NULL` | 作成日時。UPDATE では書き換えない |
| `updated_at` | `timestamptz` | `DEFAULT now() NOT NULL` + `BEFORE UPDATE` トリガ | 更新日時 |

- リポジトリ層のINSERT / UPDATE文に監査列を書かない（`SET updated_at = NOW()` も含む）。`RETURNING` での読み出しは可。
- タイムゾーン付き（`timestamptz`）で統一する。Cloud Run / Supabase / ローカルでコンテナのTZが異なるため、`timestamp`（タイムゾーンなし）を混在させると暗黙変換で値がずれる。

### 7.1 現状との乖離（未対応）

現行スキーマは上記方針を満たしておらず、以下の既知の問題がある。

| 問題 | 影響 | 該当箇所 |
|---|---|---|
| `updated_at` を更新する `BEFORE UPDATE` トリガが存在しない。Goコード側にも `updated_at` への代入がない | **`updated_at` が作成時刻のまま更新されない** | `blog_users` / `blogs` / `blog_likes` |
| `created_at` は `timestamptz`、`updated_at` は `timestamp`（TZなし）で型が不整合 | TZ差異による値のずれ | 同上 |

- 該当は `work/backup.sql`（Supabase本番ダンプ）と `backend/testsupport/testdata/schema.sql`（IT/E2E用）の双方。
- `blog_comments` は `updated_at` を持たない（更新機能がないため現状は問題にならない）。
- スキーマ修正は影響範囲が異なるため別タスク（issue #117）として扱う。
