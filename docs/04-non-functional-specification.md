# 非機能仕様書

本ドキュメントは、ブログWebアプリケーションバックエンド（Go + Echo フレームワーク）の非機能要件を、実装コードから逆引きして整理したものである。

## 目次

- [1. パフォーマンス](#1-パフォーマンス)
  - [1.1 データベース接続プール（pgxpool）](#11-データベース接続プールpgxpool)
  - [1.2 Cloud Run オートスケーリング](#12-cloud-run-オートスケーリング)
  - [1.3 HTTPレスポンス](#13-httpレスポンス)
- [2. スケーラビリティ](#2-スケーラビリティ)
  - [2.1 ステートレス設計](#21-ステートレス設計)
  - [2.2 Cloud Run 水平スケーリング](#22-cloud-run-水平スケーリング)
- [3. 可用性](#3-可用性)
  - [3.1 Cloud Run マネージド可用性](#31-cloud-run-マネージド可用性)
  - [3.2 グレースフルシャットダウン](#32-グレースフルシャットダウン)
- [4. 信頼性](#4-信頼性)
  - [4.1 Echo Recover ミドルウェア](#41-echo-recover-ミドルウェア)
  - [4.2 エラーハンドリングパターン](#42-エラーハンドリングパターン)
  - [4.3 構造化ログ](#43-構造化ログ)
  - [4.4 Echo Logger ミドルウェア](#44-echo-logger-ミドルウェア)
  - [4.5 初期化時のフェイルファスト](#45-初期化時のフェイルファスト)
- [5. 保守性](#5-保守性)
  - [5.1 クリーンアーキテクチャ層](#51-クリーンアーキテクチャ層)
  - [5.2 インターフェースベースの依存性注入（DI）](#52-インターフェースベースの依存性注入di)
  - [5.3 テスタビリティ](#53-テスタビリティ)
  - [5.4 設定の外部化](#54-設定の外部化)
  - [5.5 Dockerコンテナ構成](#55-dockerコンテナ構成)
  - [5.6 Infrastructure as Code](#56-infrastructure-as-code)

---

## 1. パフォーマンス

### 1.1 データベース接続プール（pgxpool）

Supabase（PostgreSQL）への接続には `jackc/pgx/v4/pgxpool` を使用し、以下の設定でコネクションプールを構成している。

| 設定項目 | 値 | 説明 |
|---|---|---|
| MaxConns | 10 | 同時接続の最大数 |
| MaxConnIdleTime | 30秒 | アイドル状態のコネクションが破棄されるまでの時間 |
| PreferSimpleProtocol | true | Prepared Statementの競合を防ぐためにSimple Protocolを優先 |
| SSLMode | require | 接続URLに `?sslmode=require` を付与 |

- 起動時に `Pool.Ping()` で接続確認を行い、さらに `SELECT 1` のテストクエリを実行して疎通を検証する。
- コネクションプールはグローバル変数 `supabase.Pool` として保持され、全リポジトリから共有される。

### 1.2 Cloud Run オートスケーリング

- Cloud Run上にデプロイされるため、リクエスト量に応じた自動スケーリングが適用される。
- コンテナリソース制限: CPU 1000m（1コア）、メモリ 512Mi。
- トラフィックの100%を最新リビジョンに割り当てる設定（`latest_revision = true`）。

### 1.3 HTTPレスポンス

- JSON形式でのレスポンス返却を標準とする。
- いいねデータ取得時には `Cache-Control: no-store` ヘッダーを設定し、キャッシュを無効化する。

## 2. スケーラビリティ

### 2.1 ステートレス設計

- 認証はJWTトークンベースで実装されており、サーバー側にセッション情報を保持しない。
- JWTトークンはHTTP-only Cookieに格納され、各リクエストでクライアントから送信される。
- 訪問者追跡（visit-id）もJWTトークン + Cookieベースであり、サーバー側状態を持たない。
- これにより、Cloud Runの水平スケーリング（インスタンス追加）が問題なく機能する。

### 2.2 Cloud Run 水平スケーリング

- Cloud Runのマネージドオートスケーリングにより、リクエスト数の増加に応じてコンテナインスタンスが自動追加される。
- ステートレス設計のため、任意のインスタンスがどのリクエストにも応答可能。

## 3. 可用性

### 3.1 Cloud Run マネージド可用性

- Google Cloud Runのマネージドプラットフォーム上で稼働するため、基盤レベルの高可用性はGCPが保証する。
- ヘルスチェックエンドポイント `GET /` が用意されており、`"Service is running"` を返す。

### 3.2 グレースフルシャットダウン

アプリケーションは `SIGINT` および `SIGTERM` シグナルを検知し、以下の順序で安全にシャットダウンを行う。

1. シグナル受信のログ出力: `"Shutting down server..."`
2. Echoサーバーの停止: `e.Close()`
3. Supabaseコネクションプールのクローズ: `supabase.ClosePool()`

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
```

- サーバー停止時に `http.ErrServerClosed` が返却された場合はFatalとしない（正常終了と判定）。

## 4. 信頼性

### 4.1 Echo Recover ミドルウェア

- `middleware.Recover()` を使用しており、パニックが発生した場合でもプロセスがクラッシュせず、500エラーレスポンスを返す。
- リクエストごとにリカバリが適用されるため、1つのリクエストの障害が他のリクエストに影響しない。

### 4.2 エラーハンドリングパターン

各レイヤーで以下のエラーハンドリングパターンを統一的に採用している。

- **Handler層**: エラーメッセージの文字列マッチングにより、適切なHTTPステータスコード（400, 401, 404, 500）を返す。
- **Service層**: バリデーションエラーを `errors.New()` で生成し、上位層に返却する。
- **Repository層**: データベースエラーをログ出力した上で、上位層にそのまま返却する。

エラーレスポンスはすべて `{"error": "メッセージ"}` のJSON形式で統一されている。

### 4.3 構造化ログ

5段階のログレベルを定義し、目的に応じて使い分けている。

| ログレベル | 出力先 | プレフィックス | 用途 |
|---|---|---|---|
| Info | stdout | `INFO:` | 正常処理の記録 |
| Error | stderr | `ERROR:` | エラーの記録 |
| Warn | stdout | `WARN:` | 警告の記録 |
| Debug | stdout | `DEBUG:` | デバッグ情報 |
| Test | stdout | `TEST:` | テスト用ログ |

- 各ログには日付、時刻、ファイル名（短縮形）が自動付与される（`log.Ldate|log.Ltime|log.Lshortfile`）。
- `TEST_MODE=true` の環境変数が設定されている場合、Info/Error/Warn/Debug のログ出力が `io.Discard` に抑制される。
- リクエストコンテキスト付きログユーティリティ（`utils/log`）により、HTTPメソッド、リクエストパス、User-Agentを含むログを出力可能。

### 4.4 Echo Logger ミドルウェア

- `middleware.Logger()` により、すべてのHTTPリクエスト/レスポンスのアクセスログが自動記録される。

### 4.5 初期化時のフェイルファスト

- Supabase接続の初期化に失敗した場合、`log.Fatalf` でプロセスを即座に終了する。
- `JWT_SECRET_KEY` 環境変数が未設定の場合、`config` パッケージの `init()` で `log.Fatal` が呼ばれプロセスが終了する。

## 5. 保守性

### 5.1 クリーンアーキテクチャ層

以下の4層構造で責務を明確に分離している。

| レイヤー | ディレクトリ | 責務 |
|---|---|---|
| Handler | `handlers/` | HTTPリクエスト/レスポンスの処理、バリデーション結果に応じたステータスコード返却 |
| Service | `services/` | ビジネスロジック、入力バリデーション |
| Repository | `repositories/` | データベースアクセス（SQL実行、データスキャン） |
| Model | `models/` | データ構造の定義 |

各ドメインごとにサブディレクトリが分けられている（`blogs/`, `blog_users/`, `blog_likes/`, `blog_comments/`, `auth/`）。

### 5.2 インターフェースベースの依存性注入（DI）

すべてのService層およびRepository層はインターフェースとして定義され、コンストラクタ関数でインスタンスを生成する。

- `BlogRepository` インターフェース → `BlogRepositoryImpl` 構造体
- `BlogService` インターフェース → `BlogServiceImpl` 構造体（`BlogRepository` を注入）
- `CookieUtils` インターフェース → `CookieUtilsImpl` 構造体
- `AuthService` インターフェース → `AuthServiceImpl` 構造体

依存関係の組み立ては `routes/routes.go` の `SetupRoutes()` で行われる。

### 5.3 テスタビリティ

- 各インターフェースに対応するモックファイル（`*_mock.go`）が用意されている。
- `testify` パッケージ（`github.com/stretchr/testify`）を使用したユニットテストが各レイヤーに存在する。
- Cookie操作のモック（`cookie_mock.go`）により、認証が必要なハンドラのテストが可能。
- `TEST_MODE` 環境変数によるログ抑制でテスト出力をクリーンに保つ。

### 5.4 設定の外部化

- 環境依存の設定はすべて環境変数で管理する（`PORT`, `JWT_SECRET_KEY`, `ENV`, `ALLOWED_ORIGINS`, `SUPABASE_URL`）。
- `.env` ファイルのロードには `joho/godotenv` を使用し、ファイルが存在しない場合はスキップする（Cloud Run環境では環境変数がSecret Managerから注入される）。

### 5.5 Dockerコンテナ構成

- マルチステージビルドを採用し、ビルドステージ（`golang:1.19`）と実行ステージ（`gcr.io/distroless/base`）を分離。
- distrolessベースイメージにより、最小限のランタイム環境でセキュアかつ軽量なコンテナを実現。

### 5.6 Infrastructure as Code

- Terraformによりインフラ構成をコード管理（`terraform/` ディレクトリ）。
- GCPプロバイダ `~> 5.0`、Terraform `>= 1.6` を要件とする。
- Artifact Registry、Cloud Run、Secret Manager、IAMのリソースをすべてTerraformで定義。
