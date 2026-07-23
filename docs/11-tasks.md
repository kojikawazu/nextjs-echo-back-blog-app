# タスク

本ドキュメントは、ブログバックエンドAPIプロジェクトにおける完了済みタスクと今後の改善タスクをまとめたものである。完了済みタスクは実際のコードベースから逆算して記録し、改善タスクは既知の制限事項・技術的負債に基づいて定義している。

## 目次

- [完了済みタスク](#完了済みタスク)
  - [バックエンドAPI実装](#バックエンドapi実装)
  - [認証システム](#認証システム)
  - [アーキテクチャ](#アーキテクチャ)
  - [インフラストラクチャ](#インフラストラクチャ)
  - [テスト](#テスト)
- [改善タスク（未着手）](#改善タスク未着手)
  - [セキュリティ（優先度: 高）](#セキュリティ優先度-高)
  - [データベース・モデル（優先度: 中）](#データベースモデル優先度-中)
  - [ビルド・開発環境（優先度: 中）](#ビルド開発環境優先度-中)
  - [コード品質（優先度: 低）](#コード品質優先度-低)
  - [インフラストラクチャ（優先度: 低）](#インフラストラクチャ優先度-低)
- [タスク優先度サマリー](#タスク優先度サマリー)

---

## 完了済みタスク

### バックエンドAPI実装

- [x] **ブログCRUDエンドポイント実装**
  - 全件取得（`GET /api/blogs`）
  - ユーザー別取得（`GET /api/blogs/user/:userId`）
  - 個別取得（`GET /api/blogs/detail/:id`）
  - カテゴリ一覧取得（`GET /api/blogs/categories`）
  - タグ一覧取得（`GET /api/blogs/tags`）
  - 人気記事取得（`GET /api/blogs/popular/:count`）
  - 新規作成（`POST /api/blogs/create`）
  - 更新（`PUT /api/blogs/update/:id`）
  - 削除（`DELETE /api/blogs/delete/:id`）

- [x] **ブログユーザーエンドポイント実装**
  - ユーザー情報取得（`GET /api/users/detail`）
  - ユーザー情報更新（`PUT /api/users/update`）

- [x] **いいね機能エンドポイント実装**
  - 訪問者ID別いいね一覧取得（`GET /api/blog-likes`）
  - 訪問者ID生成（`GET /api/blog-likes/generate-visit-id`）
  - いいね状態確認（`GET /api/blog-likes/is-liked/:blogId`）
  - いいね追加（`POST /api/blog-likes/create/:blogId`）
  - いいね削除（`DELETE /api/blog-likes/delete/:blogId`）

- [x] **コメント機能エンドポイント実装**
  - ブログ別コメント取得（`GET /api/comments/blog/:blogId`）
  - コメント新規作成（`POST /api/comments/create`）

- [x] **ヘルスチェックエンドポイント実装**
  - サーバー稼働確認（`GET /`）

### 認証システム

- [x] **JWT認証システム実装**
  - ログイン（`POST /api/users/login`）: メール・パスワードによる認証、JWTトークン発行
  - 認証確認（`GET /api/users/auth-check`）: クッキー内JWTトークンの検証
  - ログアウト（`POST /api/users/logout`）: クッキー削除によるトークン無効化
  - JWTクレーム: ユーザーID、メール、ユーザー名を含む
  - 有効期限: 1時間
  - 署名アルゴリズム: HS256

- [x] **訪問者ID（Visit ID）トラッキング実装**
  - UUID v4による一意の訪問者ID生成
  - JWTトークン形式で`visit-id-token`クッキーに格納
  - 有効期限: 1時間
  - 既存Visit IDの重複生成防止チェック

- [x] **クッキー管理ユーティリティ実装**
  - 認証クッキーの追加・削除
  - 環境別クッキー設定（本番: Secure=true/SameSite=None、開発: Secure=false/SameSite=Lax）
  - トークン検証ユーティリティ

### アーキテクチャ

- [x] **クリーンアーキテクチャ（3層構造）実装**
  - Handler層: HTTPリクエスト/レスポンス処理
  - Service層: ビジネスロジック、バリデーション
  - Repository層: データベースアクセス（Supabase/PostgreSQL）

- [x] **インターフェースベースの依存性注入（DI）**
  - 全ドメイン（auth, blogs, blog_users, blog_likes, blog_comments）でインターフェース定義
  - モック実装（`*_mock.go`）の提供

- [x] **ミドルウェア設定**
  - Echoロガーミドルウェア
  - リカバリーミドルウェア
  - CORSミドルウェア（環境変数による許可オリジン設定）

- [x] **ロガー実装**
  - ログレベル別出力（INFO, ERROR, WARN, DEBUG, TEST）
  - テストモード時のログ抑制機能

- [x] **Supabaseクライアント実装**
  - pgxpoolによるコネクションプール管理（最大10接続、アイドルタイム30秒）
  - Simple Protocol優先設定（Prepared Statement競合防止）
  - グレースフルシャットダウン時のプールクローズ
  - 起動時テストクエリによる接続確認

### インフラストラクチャ

- [x] **Cloud Run移行（AWS App Runner → Google Cloud Run）**
  - コスト削減: 月額 $12 → $0
  - アプリケーションコード変更なし
  - 移行レポート作成（`docs/cloud-run-migration-report.md`）

- [x] **Terraform IaCセットアップ**
  - GCPプロバイダー設定（`terraform/main.tf`）
  - Cloud Runサービス定義（`terraform/cloud-run.tf`）
  - Artifact Registryリポジトリ定義（`terraform/artifact-registry.tf`）
  - Secret Manager設定（`terraform/secrets.tf`）
  - IAMロール設定（`terraform/cloud-run-iam.tf`）
  - 変数定義（`terraform/variables.tf`）

- [x] **CI/CD（GitHub Actions）**
  - mainブランチへのpush時自動デプロイ
  - Dockerイメージビルド & Artifact Registryへのプッシュ
  - Cloud Runへのデプロイ
  - 古いコンテナイメージの自動クリーンアップ（最新5件を保持）

- [x] **Docker マルチステージビルド**
  - ビルドステージ: `golang:1.19`
  - 実行ステージ: `gcr.io/distroless/base`（軽量・セキュア）

- [x] **Supabaseプロジェクト移行**
  - pg_dump/psqlによるデータ移行
  - 4テーブル（blog_users, blogs, blog_likes, blog_comments）の移行
  - 移行手順書作成（`docs/supabase-migration-guide.md`）

### テスト

- [x] **単体テスト・結合テスト実装**
  - Handler層テスト: blogs, auth, blog_users, blog_likes, blog_comments
  - Service層テスト: blogs, auth, blog_comments, blog_likes
  - Repository層テスト: blogs, blog_users, blog_comments, blog_likes
  - テストフレームワーク: `testify v1.9.0`
  - モックベースのテスト（インターフェースDIを活用）

- [x] **IT（Repository層）の testcontainers 化・CI 導入**
  - 実 Supabase 直結から、使い捨て PostgreSQL コンテナ（`testcontainers-go` v0.35.0 + postgres モジュール）へ移行
  - 共有ヘルパー `testsupport`（`//go:build integration`）で TestMain を統一。`SUPABASE_URL` 指定時は実 DB へフォールバック
  - 決定的シード（`testsupport/testdata/schema.sql` / `seed.sql`）・接続断の異常系追加・脆いアサーション（SQLSTATE 全文一致）を緩和
  - `supabase/client.go` の sslmode を `DB_SSLMODE`（既定 `require`）で環境変数化
  - CI（`ci.yml`）に `integration-test` ジョブを追加（`go test -tags=integration ./repositories/...`）

- [x] **E2E（API-E2E）の実装・CI 導入**
  - `httptest` で本番同等の Echo アプリ（middlewares + routes）を起動し、実 HTTP フローで Handler→Service→Repository→DB を検証
  - DB は `testsupport`（testcontainers）を再利用（`Setup()` を切り出し、build tag を `integration || e2e` に拡張）
  - `e2e/`（`//go:build e2e`）に集約。範囲は正常系＋準正常系（認証フロー・ブログCRUD・いいね/コメント・ヘルスチェック）
  - CI（`ci.yml`）に `e2e-test` ジョブを追加（`JWT_SECRET_KEY` を env 注入・`go test -tags=e2e ./e2e/...`）

## 改善タスク（未着手）

### セキュリティ（優先度: 高）

- [ ] **パスワードのハッシュ化（bcrypt）**
  - 現状: `blog_users.password`が平文保存
  - 対応: `golang.org/x/crypto/bcrypt`を使用したハッシュ化
  - 影響範囲: ユーザー登録、ログイン認証、既存パスワードのマイグレーション
  - 備考: `golang.org/x/crypto`は既にgo.modに含まれている

- [ ] **レート制限の導入**
  - 現状: 全エンドポイントにレート制限なし
  - 対応: Echo の`middleware.RateLimiter`またはカスタムミドルウェアの導入
  - 優先対象: ログインエンドポイント（ブルートフォース攻撃防止）

- [ ] **CSRF保護の追加**
  - 現状: クッキーベース認証にもかかわらずCSRFトークン検証なし
  - 対応: Echo の`middleware.CSRF`の導入、またはカスタムCSRFトークン実装
  - 備考: SameSite=None設定のため、CSRF攻撃のリスクが存在する

- [ ] **入力サニタイズ（XSS防止）**
  - 現状: ユーザー入力（コメント本文、ゲスト名等）のサニタイズ処理なし
  - 対応: HTMLエスケープ処理の追加、またはサニタイズライブラリの導入
  - 優先対象: コメント作成エンドポイント

### データベース・モデル（優先度: 中）

- [ ] **タグテーブルの正規化**
  - 現状: タグがカンマ区切り文字列として`blogs.tags`に格納
  - 対応: `tags`テーブルと`blog_tags`中間テーブルの作成（多対多リレーション）
  - 効果: タグ別検索の効率化、タグ統計の容易化、データ整合性の向上

- [ ] **likes / comment_cnt の型変更（int8 → int32）**
  - 現状: `models.BlogData`の`Likes`と`CommentCnt`が`int8`型（最大値127）
  - 対応: Go側の型を`int32`に変更、DBカラムを`integer`型に変更
  - 理由: 人気記事でのオーバーフロー防止

- [ ] **APIページネーション実装**
  - 現状: ブログ一覧・コメント一覧が全件取得
  - 対応: `limit`/`offset`パラメータまたはカーソルベースのページネーション追加
  - 対象エンドポイント: `GET /api/blogs`, `GET /api/blogs/user/:userId`, `GET /api/comments/blog/:blogId`

### ビルド・開発環境（優先度: 中）

- [x] **Go バージョンの統一**
  - 対応: `go.mod`を`go 1.22`、`Dockerfile`のビルドステージを`golang:1.22`に統一
  - 備考: testcontainers-go（Go 1.21+ 要求）導入に合わせて実施

### コード品質（優先度: 低）

- [ ] **リクエストバリデーションミドルウェアの追加**
  - 現状: 各Service層で個別にバリデーション実装
  - 対応: Echo のバリデータ（`echo.Validator`）を用いた統一的なリクエストバリデーション
  - 効果: バリデーションロジックの重複削減、エラーレスポンス形式の統一

- [ ] **構造化ログの導入（JSON形式）**
  - 現状: 標準`log`パッケージベースのテキスト形式ログ
  - 対応: `zerolog`または`zap`による構造化ログへの移行
  - 効果: Cloud Loggingとの統合改善、ログ検索・分析の効率化

- [ ] **エラーハンドリングの統一**
  - 現状: エラーメッセージの文字列比較による分岐（`err.Error() == "..."`)
  - 対応: カスタムエラー型の導入、またはエラーコード体系の整備
  - 効果: エラー処理の堅牢化、APIエラーレスポンスの標準化

- [ ] **認証ミドルウェアの導入**
  - 現状: 認証が必要なエンドポイントの保護が不明確
  - 対応: JWT検証ミドルウェアを作成し、認証が必要なルートグループに適用
  - 効果: 認証ロジックの一元管理、エンドポイント保護の明確化

### インフラストラクチャ（優先度: 低）

- [ ] **ヘルスチェックの強化**
  - 現状: ルートパス（`GET /`）で文字列を返すだけのシンプルなヘルスチェック
  - 対応: データベース接続状態を含むヘルスチェックエンドポイントの追加
  - 効果: Cloud Runのヘルスチェック設定との連携改善

- [ ] **Graceful Shutdownのタイムアウト設定**
  - 現状: シグナル受信時に即時シャットダウン
  - 対応: `context.WithTimeout`を用いたタイムアウト付きシャットダウン
  - 効果: 処理中リクエストの安全な完了

## タスク優先度サマリー

| 優先度 | カテゴリ | タスク数 | 主要タスク |
|--------|---------|---------|-----------|
| 高 | セキュリティ | 4 | パスワードハッシュ化、CSRF保護、レート制限、入力サニタイズ |
| 中 | データベース | 3 | タグ正規化、型変更、ページネーション（Goバージョン統一は完了） |
| 低 | コード品質・インフラ | 5 | バリデーションMW、構造化ログ、エラーハンドリング、認証MW、ヘルスチェック強化 |
