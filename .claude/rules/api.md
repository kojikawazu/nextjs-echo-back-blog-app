---
description: Go + Echo バックエンド API 設計・レイヤー分離
globs: "backend/**"
---

# API ルール（Go + Echo）

## アーキテクチャ

- **レイヤード（Handler → Service → Repository）**。各層はインターフェースを介して依存し、コンストラクタインジェクションで組み立てる。
- **ハンドラーは薄く保つ**。責務はリクエストのバインド・パラメータ取得・Cookie/JWT の解釈・Service 呼び出し・レスポンス DTO 返却のみ。ビジネスロジック・業務上の入力検証（必須・形式・範囲）をハンドラに書かない（パスパラメータの型変換などリクエスト形式のチェックはハンドラで行ってよい。[`error-handling.md`](error-handling.md)）。
- **Service** は入力検証とビジネスロジックを担い、Repository を呼び出す。HTTP（ステータス・`echo.Context`）を扱わない。
- **Repository** はインターフェースを定義し、実装を差し替え可能にする（Service はインターフェースにのみ依存する）。SQL の実行と結果のマッピングを担い、ドライバ固有のエラー（`pgx.ErrNoRows` 等）はリポジトリの sentinel（`ErrXxxNotFound`）に変換して外に漏らさない。

## ディレクトリ構成

Go アプリは `backend/`（module `backend`）に置く。

```
backend/
├── main.go                 # エントリポイント
├── routes/                 # ルーティング定義・DI の組み立て
├── middlewares/            # ミドルウェア・HTTPErrorHandler
├── handlers/<domain>/      # Handler 層（薄い。委譲のみ）
├── services/<domain>/      # Service 層（入力検証・ビジネスロジック。errors.go に sentinel error）
├── repositories/<domain>/  # Repository 層（interface + 実装 + モック）
├── models/                 # ドメインモデル（API 都合の JSON タグを持たない。JWT Claims の直列化用タグは可）
├── dto/                    # レスポンス DTO・エラーレスポンス
├── apperrors/              # エラーの分類（Kind）・コード・共通 sentinel
├── config/ logger/ supabase/ utils/  # 設定・ログ・DB 接続・共通処理
├── testsupport/            # IT / E2E 共有ヘルパー
└── e2e/                    # API-E2E
```

- ドメイン単位（`blogs` / `blog_users` / `blog_likes` / `blog_comments` / `auth`）でパッケージを分ける。新しいドメインも同じ形で各層に追加する。
- 使われないパッケージ・中身のない実装を置かない。

### 規模が大きくなったら検討するもの（現時点では採用しない）

現在の規模（4 テーブル・CRUD 中心）では差分とリスクに見合わないため採用しない（#125）。ドメインや読み取りの最適化要件が増えた時点で再検討する。

- **Command / Query の分離**（書き込みは usecase、読み込みは読み取り専用に最適化した query-service）
- **`internal/` によるパッケージの隠蔽**、`cmd/api/main.go` へのエントリポイント移動

## レスポンス DTO（DB をそのまま返さない）

- ドメインモデル・DB 構造体をそのまま `c.JSON()` しない。**レスポンス DTO 構造体**（`dto/`）にマッピングし、公開してよいフィールドだけを厳選して返す（内部フィールド・機密の漏洩防止、API 契約と内部スキーマの疎結合化）。
- JSON タグはレスポンス DTO 側で管理する（ドメイン構造体に API 都合の JSON タグを持たせない）。

## エラー

- Service・Handler は `apperrors.Error` の sentinel（ドメイン固有は `services/<domain>/errors.go`、共通は `apperrors`）を返し、呼び出し側は `errors.Is` で判定する。**エラーメッセージの文字列比較で分岐しない**。
- 想定外のエラーは `fmt.Errorf("...: %w", err)` で包んで返す（原因を辿れるようにする）。
- ハンドラはエラーを `return` するだけにし、`c.JSON` でエラーを書かない。ステータスと統一エラーレスポンス（`dto.ErrorResponse`）への変換は `middlewares.HTTPErrorHandler` に集約する。

## 共通方針

- RESTful 設計（リソース指向エンドポイント）
- レスポンス形式: JSON（`c.JSON()`）
- エラーは必ずハンドリングする（`_` で無視しない）。関数は最後の戻り値で error を返す。
- 依存注入は `routes.SetupRoutes` でコンストラクタインジェクション。
- HTTP ステータス: 400（不正リクエスト）、401（未認証）、403（権限不足）、404（不在）、409（競合）、429（レートリミット超過）、500（サーバーエラー）
- センシティブ情報をログ・レスポンスに含めない。
