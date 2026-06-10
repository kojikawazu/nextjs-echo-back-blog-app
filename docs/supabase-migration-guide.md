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

アプリケーションを再起動して接続を確認する。

```bash
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
