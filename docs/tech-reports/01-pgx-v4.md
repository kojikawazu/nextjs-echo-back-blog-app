# pgx v4 理解向上レポート

本レポートは、Go言語のPostgreSQLドライバー **pgx v4** について、本プロジェクト（nextjs-echo-back-blog-app）の実装を題材に技術理解を深めるためのドキュメントである。

---

## 目次

- [1. pgx とは何か](#1-pgx-とは何か)
  - [1.1 概要](#11-概要)
  - [1.2 pgx を選ぶ理由](#12-pgx-を選ぶ理由)
  - [1.3 pgx v4 と v5 の違い](#13-pgx-v4-と-v5-の違い)
- [2. コネクションプール（pgxpool）](#2-コネクションプールpgxpool)
  - [2.1 プールの概念](#21-プールの概念)
  - [2.2 本プロジェクトの設定](#22-本プロジェクトの設定)
  - [2.3 各設定項目の詳細解説](#23-各設定項目の詳細解説)
- [3. クエリ実行パターン](#3-クエリ実行パターン)
  - [3.1 Pool.Query() — 複数行の取得](#31-poolquery--複数行の取得)
  - [3.2 Pool.QueryRow() — 単一行の取得](#32-poolqueryrow--単一行の取得)
  - [3.3 Pool.Exec() — 結果を返さない操作](#33-poolexec--結果を返さない操作)
- [4. パラメータ化クエリ（SQLインジェクション対策）](#4-パラメータ化クエリsqlインジェクション対策)
  - [4.1 仕組み](#41-仕組み)
  - [4.2 本プロジェクトでの使用例](#42-本プロジェクトでの使用例)
- [5. 高度なSQLパターン](#5-高度なsqlパターン)
  - [5.1 CTE（Common Table Expression）— WITH句](#51-ctecommon-table-expression-with句)
  - [5.2 COALESCE — NULLの安全な処理](#52-coalesce--nullの安全な処理)
  - [5.3 サブクエリによる集計](#53-サブクエリによる集計)
- [6. コンテキスト（Context）の役割](#6-コンテキストcontextの役割)
  - [6.1 本プロジェクトの実装](#61-本プロジェクトの実装)
  - [6.2 コンテキストの本来の用途](#62-コンテキストの本来の用途)
  - [6.3 改善の可能性](#63-改善の可能性)
- [7. グレースフルシャットダウン](#7-グレースフルシャットダウン)
  - [7.1 本プロジェクトの実装](#71-本プロジェクトの実装)
- [8. エラーハンドリング](#8-エラーハンドリング)
  - [8.1 本プロジェクトのパターン](#81-本プロジェクトのパターン)
  - [8.2 pgx が返す代表的なエラー](#82-pgx-が返す代表的なエラー)
- [9. まとめ：pgx v4 の学習ポイント](#9-まとめpgx-v4-の学習ポイント)

---

## 1. pgx とは何か

### 1.1 概要

pgx（ピージーエックス）は、Go言語向けに設計された **PostgreSQL専用ドライバー** である。標準ライブラリの `database/sql` インターフェースを経由せず、PostgreSQLのプロトコルに直接アクセスすることで、高いパフォーマンスとPostgreSQL固有の機能への完全なアクセスを提供する。

本プロジェクトでは **pgx v4.18.3**（`github.com/jackc/pgx/v4`）を使用し、Supabase（PostgreSQL）への接続に利用している。

### 1.2 pgx を選ぶ理由

| 比較項目 | database/sql + pq | pgx |
|----------|-------------------|-----|
| PostgreSQL専用機能 | 制限あり | フルサポート |
| パフォーマンス | 標準的 | 高い（中間レイヤーなし） |
| コネクションプール | 別ライブラリ必要 | `pgxpool` 内蔵 |
| 型マッピング | 基本型のみ | PostgreSQL型を直接サポート |
| バッチ処理 | 非サポート | `Batch` APIサポート |
| LISTEN/NOTIFY | 別途実装 | ネイティブサポート |

### 1.3 pgx v4 と v5 の違い

本プロジェクトは **v4** を使用しているが、現在は **v5** が最新バージョンである。

| 項目 | pgx v4 | pgx v5 |
|------|--------|--------|
| Go最低バージョン | Go 1.17 | Go 1.19 |
| モジュールパス | `github.com/jackc/pgx/v4` | `github.com/jackc/pgx/v5` |
| コンテキスト | 全メソッドでContext必須 | 同様 |
| 型システム | `pgtype` v1 | `pgtype` v2（大幅刷新） |
| エラーハンドリング | 基本的 | `pgconn.PgError` が改善 |
| Generics | 非対応 | `pgx.CollectRows` 等のGenerics対応 |

---

## 2. コネクションプール（pgxpool）

### 2.1 プールの概念

データベース接続は作成コストが高い（TCP接続確立、認証ハンドシェイク、SSL/TLSネゴシエーション）。コネクションプールは、接続を再利用することでこのオーバーヘッドを削減する仕組みである。

```
アプリケーション
  ├── goroutine A ──→ 接続1（借用）──→ Supabase
  ├── goroutine B ──→ 接続2（借用）──→ Supabase
  ├── goroutine C ──→ 接続3（借用）──→ Supabase
  └── goroutine D ──→ 待機中（プールの上限に達した場合）
         │
    コネクションプール（pgxpool.Pool）
    ┌─────────────────────────────┐
    │ 最大10接続をプール管理       │
    │ アイドル30秒で接続を破棄     │
    │ 借用/返却を自動管理          │
    └─────────────────────────────┘
```

### 2.2 本プロジェクトの設定

`supabase/client.go` での設定内容：

```go
// グローバルなプールとコンテキスト
var Pool *pgxpool.Pool
var Ctx = context.Background()

func InitSupabase() {
    // 環境変数から接続URLを取得
    supabaseURL := os.Getenv("SUPABASE_URL")

    // SSLモードを強制
    connStr := supabaseURL + "?sslmode=require"

    // 接続URLをパース
    config, err := pgxpool.ParseConfig(connStr)

    // プール設定
    config.MaxConns = 10                              // 最大同時接続数
    config.MaxConnIdleTime = 30 * time.Second          // アイドル接続の生存時間
    config.ConnConfig.PreferSimpleProtocol = true       // Simple Protocolを優先

    // プール作成
    Pool, err = pgxpool.ConnectConfig(Ctx, config)

    // 接続確認
    Pool.Ping(Ctx)
}
```

### 2.3 各設定項目の詳細解説

#### MaxConns = 10

同時にプールから借用できる接続の上限数。

**なぜ10か？**
- Supabase無料枠はデータベース接続数に制限がある
- 個人ブログ規模のトラフィックには十分
- Cloud Runのインスタンスが複数起動した場合でも、各インスタンスが10接続で収まる

**上限に達した場合の挙動：**
```
goroutine → Pool.Query() → 空き接続なし → ブロック（待機）
                                    ↓
                              接続が返却されたら自動的に借用
```

#### MaxConnIdleTime = 30秒

使われていない接続を30秒後に自動的に閉じる。

**なぜ重要か？**
- Cloud Runはリクエストがない時にインスタンスをスケールダウンする
- 不要な接続を保持し続けるとSupabase側の接続枠を消費する
- 30秒は「頻繁なリクエスト時は接続を維持しつつ、アイドル時は解放する」バランス

#### PreferSimpleProtocol = true

PostgreSQLには2つの通信プロトコルがある：

| プロトコル | 特徴 |
|-----------|------|
| **Extended Protocol**（デフォルト） | Prepared Statement使用、パラメータ型を事前に解決、同一クエリの繰り返し実行に有利 |
| **Simple Protocol** | クエリを1回のリクエストで送信、Prepared Statementなし、接続ごとの状態管理不要 |

**本プロジェクトでSimple Protocolを選ぶ理由：**
- Supabaseの接続プーラー（PgBouncer）との互換性
- PgBouncerのトランザクションモードではPrepared Statementが接続間で共有されないため競合が発生する
- Simple Protocolならこの問題を回避できる

```
Extended Protocol:
  App → PARSE(SQL) → BIND(params) → EXECUTE → DB
        ↑ Prepared Statementが接続に紐づく（PgBouncer問題の原因）

Simple Protocol:
  App → QUERY(SQL + params) → DB
        ↑ ステートレス（PgBouncer問題なし）
```

#### sslmode=require

Supabaseへの通信を必ずSSL/TLSで暗号化する。

---

## 3. クエリ実行パターン

### 3.1 Pool.Query() — 複数行の取得

`rows` イテレータを返し、`for rows.Next()` でループ処理する。

```go
// 本プロジェクトの実装例：全ブログ取得
func (r *BlogRepositoryImpl) FetchBlogs() ([]models.BlogData, error) {
    query := `
        SELECT b.id, b.blog_user_id, b.title, b.description,
               b.github_url, b.category, b.tags,
               COALESCE(l.like_count, 0) AS likes,
               COALESCE(c.comment_count, 0) AS comment_cnt,
               b.created_at, b.updated_at
        FROM blogs b
        LEFT JOIN (
            SELECT blog_id, COUNT(*) AS like_count
            FROM blog_likes GROUP BY blog_id
        ) l ON b.id = l.blog_id
        LEFT JOIN (
            SELECT blog_id, COUNT(*) AS comment_count
            FROM blog_comments GROUP BY blog_id
        ) c ON b.id = c.blog_id
        ORDER BY b.created_at DESC`

    // (1) クエリ実行 → rowsイテレータを取得
    rows, err := supabase.Pool.Query(supabase.Ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  // ← 必ずdeferで閉じる（接続リーク防止）

    var blogs []models.BlogData

    // (2) 行を1つずつ読み取る
    for rows.Next() {
        var blog models.BlogData
        var likeCount, commentCount int

        // (3) 各カラムをGoの変数にスキャン
        err := rows.Scan(
            &blog.ID, &blog.BlogUserId, &blog.Title,
            &blog.Description, &blog.GithubUrl, &blog.Category,
            &blog.Tags, &likeCount, &commentCount,
            &blog.CreatedAt, &blog.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        // (4) 型変換（int → int8）
        blog.Likes = int8(likeCount)
        blog.CommentCnt = int8(commentCount)
        blogs = append(blogs, blog)
    }

    // (5) ループ終了後のエラーチェック（重要！）
    if rows.Err() != nil {
        return nil, rows.Err()
    }

    return blogs, nil
}
```

**重要ポイント：**

1. **`defer rows.Close()`** — 忘れると接続がプールに返却されず、接続リークが発生する
2. **`rows.Err()`** — ループ中にネットワークエラーが発生した場合、`rows.Next()` は `false` を返すだけでエラーは保持される。ループ後に必ず確認する
3. **`Scan` の引数順** — SELECT句のカラム順と完全に一致させる必要がある

### 3.2 Pool.QueryRow() — 単一行の取得

1行だけの結果を取得する場合に使用する。`rows.Next()` ループが不要でコードが簡潔になる。

```go
// 本プロジェクトの実装例：ブログ作成（INSERT + RETURNING）
func (r *BlogRepositoryImpl) CreateBlog(
    userId, title, githubUrl, category, description, tags string,
) (*models.BlogData, error) {
    query := `
        INSERT INTO blogs (blog_user_id, title, github_url, category, description, tags)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, blog_user_id, title, description, github_url,
                  category, tags, likes, comment_cnt, created_at, updated_at`

    var blog models.BlogData

    // QueryRow → 即座にScan
    err := supabase.Pool.QueryRow(
        supabase.Ctx, query,
        userId, title, githubUrl, category, description, tags,
    ).Scan(
        &blog.ID, &blog.BlogUserId, &blog.Title,
        &blog.Description, &blog.GithubUrl, &blog.Category,
        &blog.Tags, &blog.Likes, &blog.CommentCnt,
        &blog.CreatedAt, &blog.UpdatedAt,
    )

    if err != nil {
        return nil, err
    }
    return &blog, nil
}
```

**QueryRow vs Query の使い分け：**

| 用途 | メソッド |
|------|---------|
| 0〜N行の結果 | `Pool.Query()` |
| 必ず1行の結果（主キー検索、INSERT RETURNING） | `Pool.QueryRow()` |
| 結果不要（DELETE, UPDATE without RETURNING） | `Pool.Exec()` |

### 3.3 Pool.Exec() — 結果を返さない操作

DELETE文のように、影響行数のみが重要で結果データが不要な場合に使用する。

```go
// 本プロジェクトの実装例：ブログ削除
func (r *BlogRepositoryImpl) DeleteBlog(id string) error {
    query := `DELETE FROM blogs WHERE id = $1`

    _, err := supabase.Pool.Exec(supabase.Ctx, query, id)
    if err != nil {
        return err
    }
    return nil
}
```

---

## 4. パラメータ化クエリ（SQLインジェクション対策）

### 4.1 仕組み

pgxでは `$1`, `$2`, `$3` ... のプレースホルダーを使ってパラメータを渡す。

```go
// 安全 ✅ — パラメータ化クエリ
query := "SELECT * FROM blogs WHERE id = $1"
row := Pool.QueryRow(ctx, query, userInput)

// 危険 ❌ — 文字列結合（SQLインジェクションの脆弱性）
query := "SELECT * FROM blogs WHERE id = '" + userInput + "'"
row := Pool.QueryRow(ctx, query)
```

**動作原理：**
```
パラメータ化クエリの場合：
  App: QUERY("SELECT * FROM blogs WHERE id = $1", ["攻撃文字列"])
  DB:  $1を文字列リテラルとして扱う → SQLの一部として解釈されない

文字列結合の場合：
  App: QUERY("SELECT * FROM blogs WHERE id = ''; DROP TABLE blogs; --'")
  DB:  全体をSQLとして解釈 → テーブル削除される可能性
```

### 4.2 本プロジェクトでの使用例

```go
// ユーザー認証クエリ — $1, $2 で安全にパラメータを渡す
query := `
    SELECT id, name, email, password, created_at, updated_at
    FROM blog_users
    WHERE email = $1 AND password = $2`

row := supabase.Pool.QueryRow(supabase.Ctx, query, email, password)
```

---

## 5. 高度なSQLパターン

### 5.1 CTE（Common Table Expression）— WITH句

本プロジェクトのブログ更新処理では、CTEを使って「UPDATEと集計結果の取得」を1つのクエリで実行している。

```sql
-- ブログ更新クエリ（repositories/blogs/blogs_impl.go）
WITH updated_blog AS (
    -- (1) まずUPDATEを実行
    UPDATE blogs
    SET title = $2, github_url = $3, category = $4,
        description = $5, tags = $6
    WHERE id = $1
    RETURNING id, blog_user_id, title, description,
              github_url, category, tags, created_at, updated_at
)
-- (2) UPDATEの結果と集計値をJOINして返す
SELECT ub.id, ub.blog_user_id, ub.title, ub.description,
       ub.github_url, ub.category, ub.tags,
       COALESCE(l.like_count, 0) AS likes,
       COALESCE(c.comment_count, 0) AS comment_cnt,
       ub.created_at, ub.updated_at
FROM updated_blog ub
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS like_count
    FROM blog_likes GROUP BY blog_id
) l ON ub.id = l.blog_id
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS comment_count
    FROM blog_comments GROUP BY blog_id
) c ON ub.id = c.blog_id
```

**CTEの利点：**
```
CTEなし: UPDATE実行 → 別のSELECTクエリで集計取得 → 2往復
CTEあり: UPDATE + SELECT を1つのクエリで実行 → 1往復
```
ネットワーク往復を削減し、データの一貫性も保証される。

### 5.2 COALESCE — NULLの安全な処理

LEFT JOINでは結合先にデータがない場合NULLが返るが、Goの `int` 型はNULLをスキャンできない。`COALESCE` で0に変換する。

```sql
COALESCE(l.like_count, 0) AS likes
-- l.like_countがNULLなら0を返す
```

### 5.3 サブクエリによる集計

```sql
LEFT JOIN (
    SELECT blog_id, COUNT(*) AS like_count
    FROM blog_likes
    GROUP BY blog_id
) l ON b.id = l.blog_id
```

**なぜサブクエリを使うのか？**
- 直接JOINすると、ブログ×いいね×コメントの直積で行数が爆発する
- サブクエリで先に集計してからJOINすることで、正確なカウント値を取得できる

---

## 6. コンテキスト（Context）の役割

### 6.1 本プロジェクトの実装

```go
var Ctx = context.Background()
```

本プロジェクトでは `context.Background()` をグローバルに定義し、全クエリで共有している。

### 6.2 コンテキストの本来の用途

コンテキストは以下の目的で使用される：

| 用途 | 説明 | 本プロジェクト |
|------|------|---------------|
| キャンセル | リクエストが中断された時にクエリも中止 | 未使用 |
| タイムアウト | クエリに時間制限を設ける | 未使用 |
| 値の伝播 | リクエストIDなどをクエリに紐づける | 未使用 |

### 6.3 改善の可能性

```go
// 現在の実装
rows, err := Pool.Query(supabase.Ctx, query)

// 改善案：リクエストスコープのコンテキスト + タイムアウト
func (r *BlogRepositoryImpl) FetchBlogs(ctx context.Context) ([]models.BlogData, error) {
    queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    rows, err := supabase.Pool.Query(queryCtx, query)
    // ...
}
```

---

## 7. グレースフルシャットダウン

### 7.1 本プロジェクトの実装

```go
// main.go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-quit  // シグナルを待機
    log.Println("Shutting down server...")
    e.Close()             // HTTPサーバー停止
    supabase.ClosePool()  // コネクションプール解放
}()
```

**なぜプールのクローズが重要か？**
- Cloud Runがインスタンスを停止する際、SIGTERMが送信される
- プールをクローズしないと、Supabase側に「幽霊接続」が残り接続枠を浪費する
- `Pool.Close()` は全アクティブ接続を安全に切断する

---

## 8. エラーハンドリング

### 8.1 本プロジェクトのパターン

```go
// Repository層：エラーをログ出力して上位層に返す
rows, err := supabase.Pool.Query(supabase.Ctx, query)
if err != nil {
    logger.ErrorLog.Printf("Failed to fetch blogs: %v", err)
    return nil, err
}
```

### 8.2 pgx が返す代表的なエラー

| エラー | 発生条件 | 例 |
|--------|---------|-----|
| `pgconn.PgError` | SQLエラー（制約違反、構文エラー等） | `ERROR: invalid input syntax for type uuid` |
| `pgx.ErrNoRows` | QueryRowで行が見つからない | 存在しないIDで検索 |
| `context.DeadlineExceeded` | タイムアウト | クエリが時間内に完了しない |
| 接続エラー | ネットワーク障害 | Supabaseダウン時 |

---

## 9. まとめ：pgx v4 の学習ポイント

| # | ポイント | 本プロジェクトでの実践 |
|---|---------|---------------------|
| 1 | コネクションプールで接続を効率管理 | `pgxpool.Pool` で最大10接続をプール |
| 2 | Simple Protocol でPgBouncer互換性確保 | `PreferSimpleProtocol = true` |
| 3 | `$1, $2` でSQLインジェクション防止 | 全クエリでパラメータ化 |
| 4 | `defer rows.Close()` で接続リーク防止 | 全Query呼び出しで使用 |
| 5 | `rows.Err()` でイテレーション中エラー検出 | ループ後にチェック |
| 6 | `RETURNING` で INSERT/UPDATE結果を1往復で取得 | ブログ作成・更新で使用 |
| 7 | CTE + LEFT JOIN で複雑な集計を1クエリで実行 | ブログ更新時のいいね・コメント数集計 |
| 8 | グレースフルシャットダウンでプールを安全に解放 | SIGTERMハンドリングで `Pool.Close()` |
