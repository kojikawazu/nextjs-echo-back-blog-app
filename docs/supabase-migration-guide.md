# Supabase データ移行手順書

旧Supabaseプロジェクトから新Supabaseプロジェクトへデータを移行する手順。

## 目次

- [前提条件](#前提条件)
- [対象テーブル](#対象テーブル)
- [接続情報](#接続情報)
- [手順](#手順)
  - [Step 1: 移行元からデータをエクスポート](#step-1-移行元からデータをエクスポート)
  - [Step 2: ダンプファイルの確認](#step-2-ダンプファイルの確認)
  - [Step 3: 移行先へデータをインポート](#step-3-移行先へデータをインポート)
  - [Step 4: 移行結果の確認](#step-4-移行結果の確認)
  - [Step 5: アプリケーションの接続先を切り替え](#step-5-アプリケーションの接続先を切り替え)
- [トラブルシューティング](#トラブルシューティング)
  - [既存テーブルがある場合の対処](#既存テーブルがある場合の対処)
- [注意事項](#注意事項)
- [付録 A: パスワードのハッシュ化（#120）](#付録-a-パスワードのハッシュ化120)
- [付録 B: updated_at の自動更新と型の統一（#117）](#付録-b-updated_at-の自動更新と型の統一117)

---

## 前提条件

- `pg_dump` / `psql` コマンドが使用可能であること（PostgreSQL クライアントツール）
- 移行元・移行先の Supabase プロジェクトの接続情報を把握していること
- 移行元のデータベースが稼働中であること

## 対象テーブル

| テーブル名 | 説明 |
|------------|------|
| `blog_users` | ブログユーザー情報（ID, 名前, メール, パスワード） |
| `blogs` | ブログ記事（タイトル, 説明, GitHub URL, カテゴリ, タグ） |
| `blog_likes` | ブログいいね情報（ブログID, 訪問者ID） |
| `blog_comments` | ブログコメント情報（ブログID, ゲストユーザー名, コメント） |

## 接続情報

| 項目 | 値 |
|------|------|
| 移行元ホスト | `aws-0-ap-northeast-1.pooler.supabase.com` |
| 移行元ポート | `6543` |
| 移行先ホスト | `aws-1-ap-northeast-1.pooler.supabase.com` |
| 移行先ポート | `5432` |
| データベース名 | `postgres` |
| スキーマ | `public` |

> **注意**: ホスト名・ポートはプロジェクトにより異なる場合があります。
> Supabase ダッシュボードの `Settings > Database` で確認してください。

## 手順

### Step 1: 移行元からデータをエクスポート

`pg_dump` で移行元の `public` スキーマをダンプファイルに出力する。

```bash
pg_dump "postgresql://postgres.[移行元プロジェクトRef]:[パスワード]@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres" \
  --schema=public \
  -f backup.sql
```

**パラメータ説明:**

| パラメータ | 説明 |
|-----------|------|
| `postgres.[移行元プロジェクトRef]` | Supabase のユーザー名（プロジェクト参照ID付き） |
| `[パスワード]` | Supabase データベースパスワード |
| `--schema=public` | `public` スキーマのみをダンプ対象にする |
| `-f backup.sql` | 出力ファイル名 |

### Step 2: ダンプファイルの確認

エクスポートしたファイルの中身を確認し、対象テーブルが含まれていることを検証する。

```bash
# ファイルサイズの確認
ls -lh backup.sql

# テーブル定義が含まれているか確認
grep "CREATE TABLE" backup.sql
```

期待される出力（テーブル名）:
- `blog_users`
- `blogs`
- `blog_likes`
- `blog_comments`

### Step 3: 移行先へデータをインポート

`psql` で移行先のデータベースにダンプファイルを適用する。

```bash
psql "postgresql://postgres.[移行先プロジェクトRef]:[パスワード]@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres" \
  -f backup.sql
```

### Step 4: 移行結果の確認

移行先に接続し、データが正しくインポートされたことを確認する。

```bash
psql "postgresql://postgres.[移行先プロジェクトRef]:[パスワード]@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres"
```

```sql
-- テーブル一覧の確認
\dt public.*

-- 各テーブルのレコード数を確認
SELECT 'blog_users' AS table_name, COUNT(*) FROM blog_users
UNION ALL
SELECT 'blogs', COUNT(*) FROM blogs
UNION ALL
SELECT 'blog_likes', COUNT(*) FROM blog_likes
UNION ALL
SELECT 'blog_comments', COUNT(*) FROM blog_comments;
```

### Step 5: アプリケーションの接続先を切り替え

`.env` ファイルの `SUPABASE_URL` を移行先のURLに変更する。

```bash
# .env を編集
SUPABASE_URL=postgresql://postgres.[移行先プロジェクトRef]:[パスワード]@aws-1-ap-northeast-1.pooler.supabase.com:5432/postgres
```

アプリケーションを再起動して接続を確認する（Go アプリは `backend/` 配下。`backend/` 内で実行する）。

```bash
cd backend
go run main.go
```

起動ログで以下が表示されれば成功:
```
Connected to Supabase successfully
Test query successful
```

## トラブルシューティング

| 症状 | 原因 | 対処 |
|------|------|------|
| `pg_dump: error: connection to server failed` | 接続情報の誤り / IP制限 | Supabase ダッシュボードで接続文字列を再確認 |
| `ERROR: relation "xxx" already exists` | 移行先に既存テーブルがある | 移行先のテーブルを事前に削除するか、`--clean` オプションを付ける |
| `ERROR: permission denied` | 権限不足 | Supabase のデータベースパスワードが正しいか確認 |
| レコード数が一致しない | ダンプ中にデータが更新された | 移行元のアプリを停止してから再度エクスポート |

### 既存テーブルがある場合の対処

移行先に既にテーブルが存在する場合、`--clean` オプションを付けてエクスポートする。

```bash
pg_dump "postgresql://postgres.[移行元プロジェクトRef]:[パスワード]@aws-0-ap-northeast-1.pooler.supabase.com:6543/postgres" \
  --schema=public \
  --clean \
  -f backup.sql
```

## 注意事項

- 移行作業中はアプリケーションを停止し、データの整合性を確保すること
- パスワードや接続文字列をGitにコミットしないこと
- 移行前に移行先プロジェクトのバックアップを取得しておくことを推奨
- `backup.sql` は作業完了後に削除すること（認証情報がSQLコメントに含まれる場合がある）

---

## 付録 A: パスワードのハッシュ化（#120）

`blog_users.password` に平文で保存されていた既存ユーザーのパスワードを、bcrypt ハッシュに一括変換する。アプリは #120 以降、bcrypt ハッシュのみを照合する（平文の保存値では認証が通らない）。

### 前提

- 本番の Supabase プロジェクトは**他アプリと共用**のため、`public.blog_users` 以外に触れない
- `pgcrypto` 拡張が `extensions` スキーマに導入済みであること（`extensions.crypt` / `extensions.gen_salt` を使う）
- pgcrypto の `gen_salt('bf', 10)` は `$2a$10$` 形式のハッシュを生成し、Go の `golang.org/x/crypto/bcrypt` で照合できる（コストもアプリの `bcrypt.DefaultCost` と同じ 10）

### 手順

1. #120 の PR をマージし、Cloud Run へのデプロイ完了を確認する（**この時点から移行完了までログインできない**。閲覧系の機能には影響しない）
2. 直ちに移行 SQL を適用する（Supabase のマイグレーションとして記録する）
3. 移行後の確認クエリで、平文が残っていないことを確認する
4. 本番でログインできることを確認する
5. 移行前の本番ダンプ（`work/backup.sql` 等）には平文が残っているため、削除または再取得する

### 移行前の確認（値は取得しない）

```sql
SELECT
  count(*)                                                   AS users,
  count(*) FILTER (WHERE password ~ '^\$2[aby]\$[0-9]{2}\$') AS bcrypt_like,
  count(*) FILTER (WHERE password IS NULL OR password = '')  AS empty_password
FROM public.blog_users;
```

### 移行 SQL

```sql
-- 平文の行のみを対象にする（bcrypt 形式の行はスキップするため、再実行しても安全）
UPDATE public.blog_users
SET password = extensions.crypt(password, extensions.gen_salt('bf', 10))
WHERE password IS NOT NULL
  AND password <> ''
  AND password !~ '^\$2[aby]\$[0-9]{2}\$';
```

- 平文は DB の外に出さずに DB 内でハッシュ化する（通信・ログに平文が流れない）
- 空・NULL のパスワードは変換しない（該当行は認証できない状態のまま。移行前の確認で件数を把握する）

### 移行後の確認

```sql
-- users = bcrypt_like であること（平文が残っていない）
SELECT
  count(*)                                                   AS users,
  count(*) FILTER (WHERE password ~ '^\$2[aby]\$10\$')       AS bcrypt_like
FROM public.blog_users;
```

### ロールバック

- ハッシュ化は不可逆のため、平文には戻せない。アプリを #120 より前に戻すとログインできなくなる
- 問題が起きた場合はアプリを戻さず、該当ユーザーのパスワードを新しいハッシュで再設定する:

```sql
-- '<new-password>' は実行時に指定する（SQL ファイル・履歴に残さない）
UPDATE public.blog_users
SET password = extensions.crypt('<new-password>', extensions.gen_salt('bf', 10))
WHERE id = '<user-id>';
```

---

## 付録 B: updated_at の自動更新と型の統一（#117）

`blog_users` / `blogs` / `blog_likes` の `updated_at` を、UPDATE 時に DB のトリガで自動更新するようにし、型を `timestamptz NOT NULL` に揃える（[`.claude/rules/database.md`](../.claude/rules/database.md)）。

### 前提

- 本番の Supabase プロジェクトは**他アプリと共用**のため、`public.blog_users` / `public.blogs` / `public.blog_likes` と、新規の関数 `public.set_updated_at()` 以外に触れない
- 関数は `CREATE OR REPLACE` ではなく `CREATE` で作る。同名の関数が既にある場合は上書きせずにエラーで止まる（他アプリの関数を壊さない）
- 既存の `updated_at`（`timestamp`）は UTC として書き込まれている（DB の `TimeZone` は UTC。2026-10-12 の確認で、全行の `updated_at` を UTC とみなすと `created_at` と一致した）。変換は `AT TIME ZONE 'UTC'` で明示し、実行するセッションのタイムゾーンに依存させない
- Go のコードは変更しない（`timestamp` / `timestamptz` のどちらでも `time.Time` に読み込める）。**マイグレーションとアプリのデプロイはどちらが先でもよい**
- 型の変更でテーブルが書き換わり、その間は排他ロックがかかる。件数が少ない（数十行）ため一瞬で終わる

### 手順

1. 移行前の確認クエリを実行し、結果が想定どおりであることを確認する
2. SQL Editor で移行 SQL を実行する（1 トランザクション。途中で失敗した場合は全体が取り消される）
3. 移行後の確認クエリで、型・トリガ・関数を確認する
4. 本番でブログを 1 件編集し、`updated_at` が編集時刻に進むことを確認する（任意）

### 移行前の確認

```sql
BEGIN READ ONLY;
-- ① updated_at は timestamp（without time zone）・NULL 許可であること
SELECT table_name, column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema = 'public' AND table_name IN ('blog_users', 'blogs', 'blog_likes')
  AND column_name IN ('created_at', 'updated_at')
ORDER BY table_name, column_name;

-- ② 同名の関数・トリガが無いこと（0 行）
SELECT 'function' AS kind, p.proname::text AS name
FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname = 'public' AND p.proname = 'set_updated_at'
UNION ALL
SELECT 'trigger', trigger_name::text FROM information_schema.triggers
WHERE event_object_schema = 'public' AND event_object_table IN ('blog_users', 'blogs', 'blog_likes');

-- ③ updated_at が NULL の行が無いこと（すべて 0。NOT NULL にするため）
SELECT
  (SELECT count(*) FROM public.blog_users WHERE updated_at IS NULL) AS blog_users,
  (SELECT count(*) FROM public.blogs      WHERE updated_at IS NULL) AS blogs,
  (SELECT count(*) FROM public.blog_likes WHERE updated_at IS NULL) AS blog_likes;
ROLLBACK;
```

### 移行 SQL

```sql
BEGIN;

-- 共通トリガ関数（共用 DB のため OR REPLACE を使わない）
CREATE FUNCTION public.set_updated_at() RETURNS trigger
LANGUAGE plpgsql SET search_path = '' AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

-- blog_users
ALTER TABLE public.blog_users
    ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at SET DEFAULT now(),
    ALTER COLUMN updated_at SET NOT NULL;
CREATE TRIGGER trg_blog_users_set_updated_at
    BEFORE UPDATE ON public.blog_users FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- blogs
ALTER TABLE public.blogs
    ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at SET DEFAULT now(),
    ALTER COLUMN updated_at SET NOT NULL;
CREATE TRIGGER trg_blogs_set_updated_at
    BEFORE UPDATE ON public.blogs FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

-- blog_likes
ALTER TABLE public.blog_likes
    ALTER COLUMN updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    ALTER COLUMN updated_at SET DEFAULT now(),
    ALTER COLUMN updated_at SET NOT NULL;
CREATE TRIGGER trg_blog_likes_set_updated_at
    BEFORE UPDATE ON public.blog_likes FOR EACH ROW EXECUTE FUNCTION public.set_updated_at();

COMMIT;
```

- 既存行の `updated_at` は作成時刻のまま（過去の更新時刻は記録されていないため復元しない）。以降の UPDATE から正しく進む
- `SET DEFAULT now()` は既に設定済みだが、型の変更後も既定値が残ることを明示するために含める

### 移行後の確認

```sql
BEGIN READ ONLY;
-- ① updated_at が timestamptz・NOT NULL・DEFAULT now() であること
SELECT table_name, column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema = 'public' AND table_name IN ('blog_users', 'blogs', 'blog_likes')
  AND column_name = 'updated_at'
ORDER BY table_name;

-- ② 3 テーブルに BEFORE UPDATE のトリガがあること（3 行）
SELECT event_object_table, trigger_name, action_timing, event_manipulation, action_statement
FROM information_schema.triggers
WHERE event_object_schema = 'public' AND event_object_table IN ('blog_users', 'blogs', 'blog_likes')
ORDER BY event_object_table;

-- ③ 既存の値が変換でずれていないこと（updated_at = created_at の行数 = 全件。差は 1 秒以内）
SELECT t, count(*) AS total,
       count(*) FILTER (WHERE abs(extract(epoch FROM u - c)) <= 1) AS unchanged
FROM (
  SELECT 'blog_users' AS t, created_at AS c, updated_at AS u FROM public.blog_users
  UNION ALL SELECT 'blogs', created_at, updated_at FROM public.blogs
  UNION ALL SELECT 'blog_likes', created_at, updated_at FROM public.blog_likes
) x GROUP BY t ORDER BY t;
ROLLBACK;
```

③ は、移行の直後（アプリがまだ行を更新していない時点）に実行する。

### ロールバック

アプリは `updated_at` の型に依存しないため、ロールバックしてもアプリの変更は不要。

```sql
BEGIN;
DROP TRIGGER trg_blog_users_set_updated_at ON public.blog_users;
DROP TRIGGER trg_blogs_set_updated_at      ON public.blogs;
DROP TRIGGER trg_blog_likes_set_updated_at ON public.blog_likes;
DROP FUNCTION public.set_updated_at();

ALTER TABLE public.blog_users
    ALTER COLUMN updated_at DROP NOT NULL,
    ALTER COLUMN updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE public.blogs
    ALTER COLUMN updated_at DROP NOT NULL,
    ALTER COLUMN updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC';
ALTER TABLE public.blog_likes
    ALTER COLUMN updated_at DROP NOT NULL,
    ALTER COLUMN updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC';
COMMIT;
```
