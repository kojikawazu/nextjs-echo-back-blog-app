# データベース（PostgreSQL / Supabase・pgx 直叩き）

本プロジェクトは ORM を使わず `jackc/pgx` で SQL を直接実行する。ORM が自動で担保する不変条件は、**DB スキーマ側で担保する**（アプリケーションコードの規律に頼らない）。SQL の書き方そのものは [`security.md`](security.md)（インジェクション対策）に従う。

## 監査列

全テーブルに監査列を持たせ、値は **DB 側で自動設定する**。アプリケーションコードで日時を組み立てない。

| 列 | 型 | 既定 | 用途 |
|---|---|---|---|
| `created_at` | `timestamptz` | `DEFAULT now() NOT NULL` | 作成日時 |
| `updated_at` | `timestamptz` | `DEFAULT now() NOT NULL` | 更新日時（UPDATE 時はトリガで自動更新） |
| `deleted_at` | `timestamptz NULL` | なし | 論理削除日時（要件がある場合のみ） |

- **手動代入を禁止**する。リポジトリ層の INSERT / UPDATE 文に `created_at` / `updated_at` を書かない（`SET updated_at = NOW()` も含む）。`RETURNING` で読み出すのは可。
- **タイムゾーン付き（`timestamptz`）で統一**する。`timestamp`（タイムゾーンなし）を混在させない。Cloud Run / Supabase / ローカルでコンテナの TZ が異なるため、`timestamp` は暗黙変換で値がずれる。
- `created_at` は **UPDATE で書き換えない**。更新文の SET 句に含めない。
- `updated_at` の自動更新は、**共通トリガ関数 + 各テーブルの `BEFORE UPDATE` トリガ**で実現する。テーブルごとに関数を書き写さない。

```sql
-- 共通トリガ関数（1 つだけ定義する）
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 監査列を持つ各テーブルに 1 つずつ張る
CREATE TRIGGER trg_blogs_set_updated_at
    BEFORE UPDATE ON blogs
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

- **例外**: シードデータ・テストで日時を固定したい場合のみ明示指定を許容する。本番コードパス（`repositories/**`）に持ち込まない。
- 操作ユーザー（`created_by` / `updated_by`）は現時点で要件がないため持たない。導入する場合は、認証ユーザーを `context.Context` で伝搬したうえでリポジトリ層の共通ヘルパーに集約し、ユースケースごとに詰めない。

## スキーマの同期

スキーマ定義は複数箇所に存在するため、**変更時は全てを同一 PR 内で揃える**（片方だけ直すと IT が本番と乖離し、検知能力を失う）。

| ファイル | 役割 |
|---|---|
| `backend/testsupport/testdata/schema.sql` | IT / E2E（testcontainers）用スキーマ |
| `backend/testsupport/testdata/seed.sql` | IT / E2E 用シードデータ |
| `work/backup.sql` | Supabase 本番のダンプ（参照用） |
| `docs/05-data-specification.md` | テーブル定義のドキュメント |

- テーブル・列を追加変更したら、`schema.sql` を本番スキーマに一致させる。列の**型まで**一致させる（型の差異は IT をすり抜ける典型的な乖離）。
- 本番側の変更は Supabase 上のマイグレーションとして適用し、手順を [`docs/supabase-migration-guide.md`](../../docs/supabase-migration-guide.md) に沿って残す。
