# ドキュメント索引

`nextjs-echo-back-blog-app`（Go + Echo ブログバックエンド API）の仕様・設計ドキュメント一覧。プロジェクト概要はリポジトリ直下の [`../README.md`](../README.md) を参照。

ドキュメントは 4 層で構成している。

- **標準仕様書（`01`〜`11`）** — 仕様の正準。番号順に読むと全体像をつかめる。
- **移行ドキュメント** — インフラ・データ移行の手順とレポート。
- **[`tech-reports/`](./tech-reports/)** — 利用技術の理解向上レポート（pgx / JWT / Hono / クリーンアーキテクチャ）。
- **[`test-design/`](./test-design/)** — テスト設計の詳細ケース表。

## 読み進め順（おすすめ）

`01 要求 → 02 要件 → 03 機能 → 05 データ → 06 セキュリティ → 07 API → 08 テスト → 09 アーキテクチャ`。
04・10・11 は随時参照。レイヤ構成・DI・インフラ・CI/CD を把握したい場合は [`09-architecture-specification.md`](./09-architecture-specification.md) から。

## 標準仕様書

| # | ドキュメント | 概要 |
|---|---|---|
| 01 | [要求仕様書](./01-business-requirements.md) | 背景・目的・スコープ・ステークホルダー・制約・データモデル概要 |
| 02 | [要件仕様書](./02-requirements-specification.md) | 機能要件・受け入れ基準（AC）・非機能要件 |
| 03 | [機能仕様書](./03-functional-specification.md) | 各機能の処理フロー・バリデーション・挙動 |
| 04 | [非機能仕様書](./04-non-functional-specification.md) | 性能・可用性・運用などの非機能要件 |
| 05 | [データ仕様書](./05-data-specification.md) | テーブル定義・モデル構造体・リレーション |
| 06 | [セキュリティ仕様書](./06-security-specification.md) | 認証・認可・JWT・Cookie・CORS・シークレット管理 |
| 07 | [API 仕様書](./07-api-specification.md) | エンドポイント・リクエスト/レスポンス・認証・ステータスコード |
| 08 | [テスト仕様書](./08-test-specification.md) | テスト分類・テストファイル一覧・カバレッジ |
| 09 | [アーキテクチャ仕様書](./09-architecture-specification.md) | 技術スタック・レイヤ構成・DI・ディレクトリ構成・インフラ・CI/CD |
| 10 | [その他仕様書](./10-miscellaneous-specification.md) | 用語集・外部参照・補足事項 |
| 11 | [タスク](./11-tasks.md) | 完了済み実績・改善タスク（技術的負債） |

## 移行ドキュメント

| ドキュメント | 概要 |
|---|---|
| [Cloud Run 移行レポート](./cloud-run-migration-report.md) | AWS App Runner → Google Cloud Run 移行レポート |
| [Supabase 移行手順書](./supabase-migration-guide.md) | Supabase データ移行手順書 |

## tech-reports/ — 技術理解レポート

利用技術の理解を深めるための調査レポート。

| ドキュメント | 対象 |
|---|---|
| [pgx v4](./tech-reports/01-pgx-v4.md) | PostgreSQL ドライバ（pgxpool・Simple Protocol） |
| [JWT / HTTP-Only Cookie](./tech-reports/02-jwt-http-only-cookie.md) | JWT (HS256) と HTTP-Only Cookie による認証 |
| [Hono](./tech-reports/03-hono.md) | フロント側 API Route（Hono） |
| [クリーンアーキテクチャ](./tech-reports/04-clean-architecture.md) | Handler / Service / Repository の 3 層構成 |

## test-design/ — テスト設計

| ドキュメント | 対象 |
|---|---|
| [WebAPI テスト設計](./test-design/webapi-test-design.md) | WebAPI 全ドメイン強化のテストケース設計 |

## 関連

- 開発ルール: [`../CLAUDE.md`](../CLAUDE.md) と [`../.claude/rules/`](../.claude/rules/)
- ドキュメント更新の影響マップ: [`../.claude/rules/documentation.md`](../.claude/rules/documentation.md)
