---
description: ドキュメント更新・設計書管理ルール（影響マップ + opt-out の完了条件）
globs:
---

# ドキュメント

コード変更がドキュメント（CLAUDE.md / README.md / docs/）と乖離しないことを構造的に担保する。

## 完了条件（opt-out）

変更は、下記「影響マップ」の対応ドキュメントを**同一 PR 内で更新する**ことを完了条件とする。

- 更新不要と判断した場合は、**PR 説明にその理由を明記する**（省略＝未対応とみなす）。
- この乖離チェックは `/self-review` と `/pr-create` の確認対象に含まれる。

## 影響マップ（変更種別 → 更新必須ドキュメント）

「どのドキュメントだっけ？」を考えさせないための逆引き表。Go + Echo バックエンド API の実態に合わせている。

| 変更種別 | 更新必須ドキュメント |
|---|---|
| ハンドラ追加・変更、ルーティング（エンドポイント）の追加・変更、リクエスト/レスポンス形式の変更 | docs/07-api-specification.md |
| モデル定義・DB スキーマ・マイグレーション・テーブル構造の変更 | docs/05-data-specification.md |
| 認証・認可・トークン・入力検証・CORS 等セキュリティ仕様の変更 | docs/06-security-specification.md |
| ディレクトリ構成・レイヤ構成・ミドルウェア・依存関係など構造的変更 | docs/09-architecture-specification.md |
| 機能の追加・挙動変更（ユーザー視点の振る舞い） | docs/03-functional-specification.md |
| 要件・スコープの変更 | docs/01-business-requirements.md / docs/02-requirements-specification.md |
| 性能・可用性・運用など非機能要件の変更 | docs/04-non-functional-specification.md |
| テスト方針・テストケースの追加・変更 | docs/08-test-specification.md / docs/test-design/ |
| デプロイ・インフラ・環境構成（Cloud Run / Supabase 等）の変更 | docs/10-miscellaneous-specification.md / docs/cloud-run-migration-report.md / docs/supabase-migration-guide.md |
| タスク進捗・実装計画の変更 | docs/11-tasks.md |
| ルール（.claude/rules/）の追加・変更、利用技術・セットアップ手順の変更 | CLAUDE.md / README.md |

該当する変更がない場合はスキップする。

## 補足

- **設計書の管理**: タスクごとに設計書を新規作成しない。既存の仕様書ドキュメント（docs/01〜11-*.md, docs/test-design/）に追記・更新する。
