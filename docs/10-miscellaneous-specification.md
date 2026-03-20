# その他仕様書

本ドキュメントは、ブログバックエンドAPIプロジェクト（Go + Echo）に関する用語集、外部参照、設計判断の根拠、既知の制限事項、移行履歴、フロントエンド連携に関する補足仕様をまとめたものである。

## 用語集

| 用語 | 説明 |
|------|------|
| **Echo** | Go言語用の軽量・高パフォーマンスなWebフレームワーク。本プロジェクトではv4.12.0を使用。ルーティング、ミドルウェア、リクエスト/レスポンスのバインディング機能を提供する。 |
| **JWT (JSON Web Token)** | 認証トークンの標準規格。本プロジェクトでは`golang-jwt/jwt v3.2.2`を使用し、HS256アルゴリズムで署名する。有効期限は1時間。HTTP-onlyクッキーに格納して管理する。 |
| **Supabase** | PostgreSQLベースのBaaS（Backend as a Service）。本プロジェクトではデータベースとして利用し、`pgx/v4`ドライバーを通じてコネクションプール（最大10接続）で接続する。 |
| **Cloud Run** | Google Cloudのサーバーレスコンテナ実行プラットフォーム。リクエスト駆動で起動し、未使用時はゼロにスケールダウンするため、低トラフィック環境でのコスト最適化に適している。 |
| **Artifact Registry** | Google CloudのDockerイメージ管理サービス。CI/CDパイプラインでビルドしたコンテナイメージの保存先として使用する。 |
| **Secret Manager** | Google CloudのシークレットＧ管理サービス。`SUPABASE_URL`、`JWT_SECRET_KEY`、`ALLOWED_ORIGINS`等の機密情報を安全に管理し、Cloud Runの環境変数として注入する。 |
| **Visit ID（訪問者ID）** | 未認証のサイト訪問者を一意に識別するためのUUID。JWTトークンとして`visit-id-token`クッキーに格納される。いいね機能における匿名ユーザーの重複防止に使用する。 |
| **pgx** | Go言語用のPostgreSQLドライバー。本プロジェクトではv4を使用し、`pgxpool`によるコネクションプーリングを行う。Simple Protocolを優先設定し、Prepared Statementの競合を防止している。 |
| **Terraform** | HashiCorp製のIaC（Infrastructure as Code）ツール。GCPリソース（Cloud Run、Artifact Registry、Secret Manager、IAM）をコードで定義・管理する。 |
| **Distroless** | Google提供の最小限コンテナイメージ。シェルやパッケージマネージャを含まず、アプリケーションバイナリの実行に必要な最小限のランタイムのみを含む。セキュリティとイメージサイズの最適化に寄与する。 |
| **Workload Identity Federation** | GCPの認証方式。GitHub ActionsからGCPリソースへのアクセスにサービスアカウントキーを使用する。 |
| **CORS** | Cross-Origin Resource Sharing。フロントエンド（Next.js）からのクロスオリジンリクエストを許可するために設定。`ALLOWED_ORIGINS`環境変数で許可オリジンをカンマ区切りで指定する。 |
| **Clean Architecture** | 本プロジェクトで採用しているアーキテクチャパターン。Handler（プレゼンテーション層）、Service（ビジネスロジック層）、Repository（データアクセス層）の3層に分離する。 |

## 外部参照・URL

| 項目 | URL / 参照先 |
|------|-------------|
| フロントエンドリポジトリ | https://github.com/kojikawazu/nextjs-echo-blog-front-app-renew |
| Echo フレームワーク | https://echo.labstack.com/ |
| Supabase | https://supabase.com/ |
| golang-jwt | https://github.com/golang-jwt/jwt |
| pgx ドライバー | https://github.com/jackc/pgx |
| Google Cloud Run | https://cloud.google.com/run |
| Terraform GCP Provider | https://registry.terraform.io/providers/hashicorp/google/latest |
| Distroless イメージ | https://github.com/GoogleContainerTools/distroless |

## 設計判断と根拠

### JWT をHTTP-onlyクッキーに格納

認証トークンをAuthorizationヘッダーではなくHTTP-onlyクッキーに格納する方式を採用している。

- **根拠**: XSS攻撃によるトークン窃取のリスクを低減するため。HTTP-onlyフラグにより、JavaScriptからクッキー値にアクセスできなくなる。
- **実装詳細**: 本番環境では`Secure=true`、`SameSite=None`を設定し、HTTPS通信でのみ送信される。開発環境では`Secure=false`、`SameSite=Lax`を設定する。
- **クッキー名**: 認証トークンは`token`、訪問者IDトークンは`visit-id-token`。

### 匿名訪問者のUUIDトラッキング

いいね機能において、ログインしていない匿名訪問者を`visit-id`で識別する方式を採用している。

- **根拠**: ブログの閲覧者がログインせずにいいねを付けられるUXを実現するため。IPアドレスベースの識別はプライバシー上の懸念があるため、クッキーベースのUUIDを採用した。
- **制限**: クッキーの有効期限は1時間であり、期限切れ後は新しいVisit IDが発行される。

### ゲストコメント（認証不要）

コメント投稿に認証を要求しない設計を採用している。

- **根拠**: ブログ閲覧者が気軽にコメントできるようにするため。ゲスト名（`guest_user`）とコメント本文のみで投稿可能。
- **リスク**: スパムコメントへの対策が不足している（後述の制限事項参照）。

### ブログコンテンツのGitHub URL参照

ブログ本文をデータベースに格納せず、GitHubのMarkdownファイルへのURLを保持する設計を採用している。

- **根拠**: コンテンツ管理をGitHubに委譲し、バージョン管理やMarkdownエディタの利点を活用するため。フロントエンドがGitHub APIを通じてMarkdownを取得・レンダリングする。
- **フィールド**: `blogs`テーブルの`github_url`カラムに格納。

### タグのカンマ区切り文字列

タグを正規化テーブルではなく、`blogs`テーブルの`tags`カラムにカンマ区切りの文字列として格納している。

- **根拠**: 実装の簡素化。タグの種類・数が限定的であり、複雑なクエリ要件がないため。
- **処理**: サービス層（`FetchBlogTags`）でカンマ区切りを分割し、重複除去・ソートを行ってからレスポンスとして返す。

### いいね数・コメント数の非正規化

`blogs`テーブルに`likes`（いいね数）と`comment_cnt`（コメント数）を直接保持している。

- **根拠**: 一覧表示時のJOIN/COUNT集計クエリを回避し、読み取りパフォーマンスを向上させるため。
- **制限**: `int8`型で定義されているため、値の上限が127に制限される（後述の制限事項参照）。

### インターフェースベースの依存性注入（DI）

全レイヤー間の依存関係をGoのインターフェースで定義し、実装の差し替えを可能にしている。

- **根拠**: テスト容易性の確保。各レイヤーのモック実装（`*_mock.go`）を用いた単体テストが可能。
- **構成**: `Repository`インターフェース → `Service`が依存注入 → `Handler`が依存注入、という階層。

## 既知の制限事項と技術的負債

### セキュリティ

| 項目 | 詳細 | 影響度 |
|------|------|--------|
| **パスワードの平文保存** | `blog_users`テーブルの`password`カラムにパスワードが平文で格納されている。ハッシュ化（bcrypt等）が未実装。 | 高 |
| **CSRF保護なし** | Echo のCSRFミドルウェアが未設定。クッキーベース認証のため、CSRFトークンによる保護が必要。 | 高 |
| **レート制限なし** | APIエンドポイントにレート制限が設定されていない。ブルートフォース攻撃やDoS攻撃に対する防御が不足。 | 中 |
| **入力サニタイズ不足** | XSS攻撃防止のための入力サニタイズ処理が未実装。特にコメントのゲスト名・本文が対象。 | 中 |

### データ型の制約

| 項目 | 詳細 | 影響度 |
|------|------|--------|
| **likes / comment_cnt が int8** | `models.BlogData`の`Likes`と`CommentCnt`が`int8`型（最大値127）。人気記事のいいね数やコメント数がオーバーフローする可能性がある。 | 中 |
| **タグの正規化不足** | タグがカンマ区切り文字列で格納されており、タグ別検索やタグの統計集計が非効率。タグテーブルの正規化が望ましい。 | 低 |

### ビルド・バージョン

| 項目 | 詳細 | 影響度 |
|------|------|--------|
| **Go バージョンの不一致** | `go.mod`では`go 1.20`を指定しているが、`Dockerfile`のビルドステージでは`golang:1.19`を使用している。ビルド環境と開発環境のバージョン不整合がある。 | 低 |

### 機能面

| 項目 | 詳細 | 影響度 |
|------|------|--------|
| **ページネーション未実装** | ブログ一覧やコメント一覧にページネーションがなく、全件取得される。データ量増加時にパフォーマンス劣化のリスクがある。 | 中 |
| **リクエストバリデーションミドルウェア未実装** | バリデーションは各サービス層で個別に実装されており、統一的なバリデーションミドルウェアがない。 | 低 |
| **構造化ログ未対応** | ログは標準`log`パッケージベースのテキスト形式。JSON形式の構造化ログ（例: `zerolog`, `zap`）への移行が望ましい。 | 低 |

## 移行履歴

### AWS App Runner → Google Cloud Run（2026年3月）

ホスティング基盤をAWS App RunnerからGoogle Cloud Runに移行した。

- **目的**: ランニングコストの削減
- **コスト効果**: 月額 $12 → $0（年間約 $144 の削減）
- **変更範囲**:
  - Terraform: AWSプロバイダーからGCPプロバイダーに変更
  - GitHub Actions: デプロイワークフローをCloud Run向けに変更
  - コンテナイメージ保存先: Artifact Registryに変更
  - DNS: CloudFlare DNSでエンドポイントURL変更に対応
- **アプリケーションへの影響**: バックエンドAPIの仕様変更なし。Go/Echoのコード変更なし。
- **詳細**: `docs/cloud-run-migration-report.md` 参照

### Supabase プロジェクト移行

旧Supabaseプロジェクトから新Supabaseプロジェクトへデータを移行した。

- **対象テーブル**: `blog_users`, `blogs`, `blog_likes`, `blog_comments`
- **移行方式**: `pg_dump` / `psql` によるエクスポート・インポート
- **接続情報変更**: 移行元（`aws-0-ap-northeast-1`, ポート6543）→ 移行先（`aws-1-ap-northeast-1`, ポート5432）
- **詳細**: `docs/supabase-migration-guide.md` 参照

## フロントエンド連携

### フロントエンドアプリケーション

- **リポジトリ**: https://github.com/kojikawazu/nextjs-echo-blog-front-app-renew
- **フレームワーク**: Next.js
- **通信方式**: REST API（JSON形式）

### CORS設定

フロントエンドからのリクエストを許可するため、以下のCORS設定が適用されている。

- **許可オリジン**: `ALLOWED_ORIGINS`環境変数にカンマ区切りで指定
- **許可メソッド**: GET, POST, PUT, DELETE
- **許可ヘッダー**: Origin, Content-Type, Authorization, Access-Control-Allow-Credentials
- **クレデンシャル**: `AllowCredentials=true`（クッキー送受信のため必須）

### クッキー連携

フロントエンドは以下のクッキーを使用してバックエンドと連携する。

| クッキー名 | 用途 | 設定元 |
|-----------|------|--------|
| `token` | JWT認証トークン | ログイン時にバックエンドが設定 |
| `visit-id-token` | 訪問者ID（いいね機能用） | `/api/blog-likes/generate-visit-id`呼び出し時にバックエンドが設定 |

### 環境変数一覧

| 環境変数名 | 用途 | 管理方法 |
|-----------|------|---------|
| `PORT` | APIサーバーポート（デフォルト: 8080） | Cloud Run自動設定 / ローカル.env |
| `ENV` | 環境識別（`production`でHTTPS-onlyクッキー有効化） | Secret Manager |
| `JWT_SECRET_KEY` | JWTトークン署名キー | Secret Manager |
| `SUPABASE_URL` | Supabase PostgreSQL接続URL | Secret Manager |
| `ALLOWED_ORIGINS` | CORS許可オリジン（カンマ区切り） | Secret Manager |
