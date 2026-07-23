---
description: テスト分類・原則（スタック非依存）
globs: 
---

# テストルール

## テスト分類

| 分類 | 定義 |
|------|------|
| 正常系（Normal） | 期待通りの入力 → 正しい結果 |
| 準正常系（Semi-Normal） | 想定内の異常入力 → 適切なハンドリング |
| 異常系（Abnormal） | 想定外のエラー → 安全な失敗 |

## 原則

- テストは仕様の証明。テストが失敗したら実装を修正する（テストを実装に合わせない）。
- 正常系 1 : 異常系（準正常系 + 異常系）2 以上の比率を目安とする。
- ビジネスロジックをモックしない。モックは外部 I/O（HTTP通信、DB接続、ファイルシステム）のみ。
- `toBeTruthy()` 等の曖昧なアサーションを避け、具体的な値で検証する。

## テストツール

| テスト種別 | ツール |
|-----------|--------|
| ユニットテスト（UT） | Go testing（テーブルドリブン）。外部I/Oのみモック |
| インテグレーションテスト（IT） | testcontainers（実DB）。`//go:build integration` |
| E2E（API-E2E） | httptest + testcontainers。`//go:build e2e`。実サーバー起動 → HTTP フロー検証 |
| スモークテスト | scripts/smoke-test.sh（未実装） |

共有ヘルパー `backend/testsupport`（build tag `integration || e2e`）が testcontainers の PostgreSQL 起動・スキーマ/シード適用・接続を担う。IT/E2E とも Docker 稼働が前提で、`SUPABASE_URL` を指定した場合のみ実DBへ接続する。

## テストファイル配置

Go 標準どおりユニット / インテグレーションテストを**対象と同じパッケージにコロケートする**:

- **ユニット / インテグレーションテスト**: 対象と同階層に `_test.go`（例: `handler.go` → `handler_test.go`）。インテグレーションはビルドタグ（例: `//go:build integration`）で分離してよい
- **E2E テスト**: `e2e/` に集約（`//go:build e2e`）。`config` が起動時に `JWT_SECRET_KEY` を要求するため、実行時に env で渡す
