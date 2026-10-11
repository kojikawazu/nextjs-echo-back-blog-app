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

