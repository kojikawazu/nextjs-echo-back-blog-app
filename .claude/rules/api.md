---
description: Go + Echo バックエンド API 設計・レイヤー分離
globs: "backend/**"
---

# API ルール（Go + Echo）

## アーキテクチャ

- **小規模**: レイヤードアーキテクチャ（Handler → Usecase → Repository）
- **中〜大規模**: クリーンアーキテクチャ（Domain → Usecase → Interface → Infrastructure）
- **ハンドラーは薄く保つ**。責務はリクエストのバインド/バリデーション・usecase/query-service 呼び出し・レスポンス DTO 返却のみ。ビジネスロジックをハンドラに書かない。
- **書き込み（Command）は usecase**、**読み込み（Query）は query-service** に分離する（読み取りは usecase を経由せず読み取り専用に最適化）。
- リポジトリは**インターフェースを定義**し、実装を差し替え可能にする（usecase から利用）。

## ディレクトリ構成（レイヤード）

```
cmd/api/main.go        # エントリポイント・DI
internal/
├── handler/           # リクエストハンドラー（薄い。委譲のみ）
├── usecase/           # 書き込み（Command）ビジネスロジック
├── query/             # 読み込み（Query / query-service）
├── repository/        # データアクセス（interface + 実装）
├── domain/            # ドメインモデル・値オブジェクト
├── dto/               # リクエスト/レスポンス構造体
├── config/            # 設定
└── middleware/         # ミドルウェア
```

## レスポンス DTO（DB をそのまま返さない）

- ドメインモデル・DB 構造体をそのまま `c.JSON()` しない。**レスポンス DTO 構造体**にマッピングし、公開してよいフィールドだけを厳選して返す（内部フィールド・機密の漏洩防止、API 契約と内部スキーマの疎結合化）。
- JSON タグはレスポンス DTO 側で管理する（ドメイン構造体に API 都合の JSON タグを持たせない）。

## 共通方針

- RESTful 設計（リソース指向エンドポイント）
- レスポンス形式: JSON（`c.JSON()`）
- エラーは必ずハンドリングする（`_` で無視しない）。関数は最後の戻り値で error を返す。
- 依存注入は `main.go` でコンストラクタインジェクション。
- HTTP ステータス: 400（不正リクエスト）、401（未認証）、403（権限不足）、404（不在）、500（サーバーエラー）
- センシティブ情報をログ・レスポンスに含めない。
