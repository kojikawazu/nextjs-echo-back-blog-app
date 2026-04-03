# テスト設計: WebAPI 全ドメイン強化

## 対象

- 対象機能: ブログWebアプリケーションバックエンドAPI（全ドメイン）
- スタック: Go + Echo v4 / testify v1.9.0
- テスト戦略: Handler層・Service層・Repository層の3層すべてを強化

---

## 強化の方針

### 現状の課題

| 課題 | 現状 | 改善後 |
|------|------|--------|
| テスト構造 | バリデーションケースごとに別関数（重複セットアップ多数） | テーブルドリブン形式（`[]struct`）に統一 |
| レスポンス検証 | `assert.Contains(body, "string")` が多数 | `assert.JSONEq` / `json.Unmarshal` で構造検証 |
| 未実装ハンドラーテスト | blog_likes 3エンドポイント・auth 2エンドポイント未実装 | 全エンドポイント網羅 |
| 不正リクエストのテスト | 不正JSON・Content-Type無しケースなし | 異常系として追加 |
| 認証フロー統合 | Login→CheckAuth→Logoutの一連シナリオなし | pipeline_test として追加 |

### テーブルドリブン方針

```go
// 推奨パターン（既存テストのリファクタリング対象）
func TestHandler_CreateBlog(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    map[string]string
        setupMock      func(m *MockBlogService, c *MockCookieUtils)
        expectedStatus int
        expectedBody   string
    }{
        {
            name: "正常系_全フィールド正常値でブログ作成成功",
            ...
        },
        {
            name: "準正常系_titleが空文字の場合400を返す",
            ...
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) { ... })
    }
}
```

### レスポンス検証方針

```go
// 改善前（文字列部分一致）
assert.Contains(t, rec.Body.String(), "title2")

// 改善後（JSON構造検証）
assert.JSONEq(t, `{"id":"1","title":"title2"}`, rec.Body.String())
// または
var res models.BlogData
json.Unmarshal(rec.Body.Bytes(), &res)
assert.Equal(t, "title2", res.Title)
```

---

## ドメイン別テストケース設計

---

## 1. Auth（認証）ドメイン

### 対象ファイル

- `handlers/auth/auth_impl.go`
- `services/auth/auth_impl.go`

### 1-1. Handler層: Login

#### 現状

Login成功1件のみ実装済み（`handlers/auth/auth_test.go`）

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | 正しいメール・パスワードでログイン成功 | `email: test@example.com, password: pass123` | 200, `{"message":"Login successful"}`, `token` Cookie セット | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | メールが空文字 → バリデーションエラー | `email: ""` | 400, `{"error":"..."}` | Unit | High |
| S-2 | パスワードが空文字 → バリデーションエラー | `password: ""` | 400, `{"error":"..."}` | Unit | High |
| S-3 | 存在しないメール → 認証失敗 | `email: notexist@example.com` | 401, `{"error":"..."}` | Unit | High |
| S-4 | パスワード不一致 → 認証失敗 | 正しいメール + 誤パスワード | 401, `{"error":"..."}` | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | 不正JSONボディ → パースエラー | `{invalid json}` | 400 | Unit | High |
| A-2 | 空ボディ | ボディなし | 400 | Unit | Medium |
| A-3 | UserServiceがエラーを返す → 500 | Serviceがエラーを返す | 500 | Unit | High |

### 1-2. Handler層: CheckAuth

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | 有効なJWT Cookieで認証確認成功 | `token` Cookie 有効値 | 200, `{"message":"Authenticated"}` | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | Cookieなし → 未認証 | Cookie欠如 | 401 | Unit | High |
| S-2 | 期限切れトークン → 未認証 | 期限切れJWT | 401 | Unit | High |
| S-3 | 不正なトークン値 → 未認証 | `token: "invalid"` | 401 | Unit | High |

### 1-3. Handler層: Logout

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | ログアウト成功 → Cookie削除 | 有効なCookie | 200, `{"message":"Logout successful"}`, Cookie削除確認 | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | Cookieなしでもログアウト成功 | Cookie欠如 | 200 | Unit | Medium |

### 1-4. 認証フロー統合テスト（パイプライン）

#### 現状

**未実装**

#### テストケース一覧

| # | テストケース | シナリオ | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| P-1 | Login → CheckAuth → Logout の一連フロー | 各ステップを順に実行 | 各ステップで正しいレスポンス | Integration | Medium |

---

## 2. Blogs ドメイン

### 対象ファイル

- `handlers/blogs/blogs_impl.go`
- `services/blogs/blogs_impl.go`
- `repositories/blogs/blogs_impl.go`

### 2-1. Handler層: 既存テストのリファクタリング

#### リファクタリング対象

| ファイル | 現状の問題 | 改善内容 |
|---------|-----------|---------|
| `blogs_FetchBlogs_test.go` | 3関数、`assert.Contains` | テーブルドリブン化、`assert.JSONEq` |
| `blogs_CreateBlog_test.go` | 1関数（正常系のみ）、`assert.Contains` | 準正常系・異常系追加、`assert.JSONEq` |
| `blogs_FetchBlogById_test.go` | 確認が必要 | テーブルドリブン化 |

#### FetchBlogs テストケース（リファクタリング後）

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | ブログ一覧取得成功（複数件） | なし | 200, BlogData配列（全フィールド検証） | Unit | High |
| N-2 | ブログ一覧取得成功（0件） | なし | 200, `[]` | Unit | Medium |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | Serviceがエラーを返す → 500 | Service エラー | 500, `{"error":"Error fetching blogs"}` | Unit | High |

#### CreateBlog テストケース（追加分）

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | titleが空文字 → 400 | `title: ""` | 400, `{"error":"Invalid title"}` | Unit | High |
| S-2 | githubUrlが空文字 → 400 | `githubUrl: ""` | 400 | Unit | High |
| S-3 | categoryが空文字 → 400 | `category: ""` | 400 | Unit | High |
| S-4 | descriptionが空文字 → 400 | `description: ""` | 400 | Unit | High |
| S-5 | tagsが空文字 → 400 | `tags: ""` | 400 | Unit | High |
| S-6 | 認証Cookie取得失敗 → 401 | Cookie欠如 | 401 | Unit | High |
| S-7 | JWT解析失敗 → 401 | 不正JWT | 401 | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | 不正JSONボディ → 400 | `{invalid json}` | 400 | Unit | High |
| A-2 | Serviceがエラーを返す → 500 | Service エラー | 500 | Unit | High |

### 2-2. Service層: 既存テストのリファクタリング

#### リファクタリング対象

| ファイル | 現状の問題 | 改善内容 |
|---------|-----------|---------|
| `blogs_CreateBlog_test.go` | 8関数（各バリデーション別） | テーブルドリブン化で1関数に統合 |
| `blogs_FetchBlogs_test.go` | 2関数 | テーブルドリブン化 |

---

## 3. Blog Likes（いいね）ドメイン

### 対象ファイル

- `handlers/blog_likes/blog_likes_impl.go`
- `services/blog_likes/blog_likes_impl.go`

### 3-1. Handler層: GenerateVisitId

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | VisitID生成成功 → 新規Cookie発行 | Cookie未設定 | 200, `visit-id-token` Cookie セット | Unit | High |
| N-2 | すでにVisitIDが存在する場合 → 既存を返す | Cookie設定済み | 200 | Unit | Medium |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | トークン生成失敗 → 500 | ServiceがVisitIDトークン生成エラー | 500 | Unit | High |

### 3-2. Handler層: IsBlogLiked

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | いいね済み → trueを返す | 有効なblogId + visitId | 200, `{"liked":true}` | Unit | High |
| N-2 | いいね未済 → falseを返す | 有効なblogId + visitId | 200, `{"liked":false}` | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | blogIdが空 → 400 | `blogId: ""` | 400 | Unit | High |
| S-2 | visit-id-token Cookie欠如 → 500 | Cookie未設定 | 500 | Unit | High |
| S-3 | visitId取得失敗 → 500 | JWT解析エラー | 500 | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | Serviceがエラーを返す → 500 | Service エラー | 500 | Unit | High |

### 3-3. Handler層: CreateBlogLike

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | いいね作成成功 | 有効なblogId + visitId | 201, BlogLikesData | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | blogIdが空 → 400 | `blogId: ""` | 400 | Unit | High |
| S-2 | visit-id-token Cookie欠如 → 500 | Cookie未設定 | 500 | Unit | High |
| S-3 | visitId取得失敗 → 500 | JWT解析エラー | 500 | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | Serviceがエラーを返す → 500 | Service エラー | 500 | Unit | High |

### 3-4. Handler層: DeleteBlogLike

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | いいね削除成功 | 有効なblogId + visitId | 200, `{"message":"Blog like deleted"}` | Unit | High |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | blogIdが空 → 400 | `blogId: ""` | 400 | Unit | High |
| S-2 | visit-id-token Cookie欠如 → 500 | Cookie未設定 | 500 | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | Serviceがエラーを返す → 500 | Service エラー | 500 | Unit | High |

### 3-5. Handler層: FetchBlogLikesByVisitId（既存テストのリファクタリング）

#### リファクタリング対象

| 現状 | 改善内容 |
|------|---------|
| 4関数（各ケース別、重複セットアップ） | テーブルドリブン化 |
| `assert.Contains` | `assert.JSONEq` |

---

## 4. Blog Comments（コメント）ドメイン

### 対象ファイル

- `handlers/blog_comments/blog_comments_impl.go`
- `services/blog_comments/blog_comments_impl.go`

### 4-1. Handler層: FetchCommentsByBlogId

#### 現状

テストファイルは存在するが内容確認が必要

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | コメント一覧取得成功（複数件） | `blogId: "1"` | 200, BlogCommentsData配列（全フィールド検証） | Unit | High |
| N-2 | コメント一覧取得成功（0件） | `blogId: "1"` | 200, `[]` | Unit | Medium |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | blogIdが空文字 → 400 | `blogId: ""` | 400 | Unit | High |

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | Serviceがエラーを返す → 500 | Service エラー | 500 | Unit | High |

### 4-2. Handler層: CreateComment（既存テストのリファクタリング）

#### リファクタリング対象

| 現状 | 改善内容 |
|------|---------|
| 6関数（各ケース別、重複セットアップ多数） | テーブルドリブン化で1関数に統合 |
| `assert.Contains` | `assert.JSONEq` |

#### 追加するテストケース

##### 異常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| A-1 | 不正JSONボディ → 400 | `{invalid json}` | 400 | Unit | High |

---

## 5. Blog Users（ユーザー）ドメイン

### 対象ファイル

- `handlers/blog_users/blog_users_impl.go`
- `services/blog_users/blog_users_impl.go`
- `repositories/blog_users/blog_users_impl.go`

### 5-1. Handler層: 既存テストのリファクタリング

#### リファクタリング対象

| ファイル | 現状の問題 | 改善内容 |
|---------|-----------|---------|
| `blog_users_FetchUser_test.go` | 確認が必要 | テーブルドリブン化・JSONEq検証 |
| `blog_users_UpdateUser_test.go` | 確認が必要 | テーブルドリブン化・JSONEq検証 |

### 5-2. Repository層: UpdateBlogUsers

#### 現状

**未実装**

#### テストケース一覧

##### 正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| N-1 | ユーザー情報更新成功 | 有効なid, name, email, password | 更新後のBlogUsersData | Integration | Medium |

##### 準正常系

| # | テストケース | 入力 | 期待結果 | テスト種別 | 優先度 |
|---|---|---|---|---|---|
| S-1 | 存在しないID → エラー | `id: "nonexistent"` | error返却 | Integration | Medium |
| S-2 | 不正なUUID形式 → エラー | `id: "2"` | DB error (SQLSTATE 22P02) | Integration | Medium |

---

## 6. テスト実装ファイル一覧

### 新規作成ファイル

| ファイルパス | 対象 | 優先度 |
|-------------|------|--------|
| `handlers/auth/auth_CheckAuth_test.go` | CheckAuth ハンドラー | High |
| `handlers/auth/auth_Logout_test.go` | Logout ハンドラー | High |
| `handlers/auth/auth_pipeline_test.go` | Login→CheckAuth→Logout 統合 | Medium |
| `handlers/blog_likes/blog_likes_GenerateVisitId_test.go` | GenerateVisitId ハンドラー | High |
| `handlers/blog_likes/blog_likes_IsBlogLiked_test.go` | IsBlogLiked ハンドラー | High |
| `handlers/blog_likes/blog_likes_CreateBlogLike_test.go` | CreateBlogLike ハンドラー | High |
| `handlers/blog_likes/blog_likes_DeleteBlogLike_test.go` | DeleteBlogLike ハンドラー | High |
| `repositories/blog_users/blog_users_UpdateBlogUsers_test.go` | UpdateBlogUsers リポジトリ | Medium |
| `repositories/blog_comments/blog_comments_CreateComment_test.go` | CreateComment リポジトリ | Medium |

### リファクタリング対象ファイル

| ファイルパス | リファクタリング内容 | 優先度 |
|-------------|-------------------|--------|
| `handlers/blogs/test/blogs_FetchBlogs_test.go` | テーブルドリブン化 + JSONEq | High |
| `handlers/blogs/test/blogs_CreateBlog_test.go` | 準正常系・異常系追加 + JSONEq | High |
| `handlers/blog_likes/blog_likes_FetchBlogLikesByVisitId_test.go` | テーブルドリブン化 + JSONEq | Medium |
| `handlers/blog_comments/blog_comments_CreateComment_test.go` | テーブルドリブン化 + JSONEq | Medium |
| `services/blogs/test/blogs_CreateBlog_test.go` | テーブルドリブン化（8関数→1関数） | Medium |
| `handlers/auth/auth_test.go` | Login失敗ケース追加 | High |

---

## 7. テスト構成

### ユニットテスト（Handler/Service）

- モック対象: 外部I/O（DB接続、Cookie操作）のみ
- モック禁止: ビジネスロジック（Service実装の内部）
- 検証: ステータスコード + JSONレスポンス構造（`assert.JSONEq` 推奨）

### 結合テスト（Repository）

- 実際のSupabase（PostgreSQL）に接続
- 環境変数: `.env.test` から読み込み
- パイプラインテスト: Create → Read → Update → Delete の一連操作

### テストの実行

```bash
# 全テスト実行
go test ./... -v

# Handler層のみ
go test ./handlers/... -v

# Service層のみ
go test ./services/... -v

# 特定ドメイン
go test ./handlers/blog_likes/... -v
```

---

## 8. モック方針

- **モック許可**: 外部I/O（HTTP通信、DB接続、Cookie操作）のみ
- **モック禁止**: ビジネスロジック、ドメインオブジェクト、ユーティリティ関数
- **Cookieモック**: 認証が必要なエンドポイントは `SetMockBlogCookies` ヘルパーを利用

---

## 9. 優先度サマリー

| 優先度 | 作業内容 | 対象ファイル数 |
|--------|---------|--------------|
| **High** | blog_likes 未実装ハンドラーテスト（4件） | 4 |
| **High** | auth CheckAuth・Logout テスト | 2 |
| **High** | 既存テストのテーブルドリブン化 + JSONEq化 | 6+ |
| **Medium** | auth Login失敗・不正ボディ追加 | 1 |
| **Medium** | repository 補完（UpdateBlogUsers, CreateComment） | 2 |
| **Medium** | 認証フロー統合テスト（pipeline） | 1 |

---

## 実装チェックリスト

### Handler層（新規）

- [ ] handlers/auth/auth_CheckAuth_test.go
- [ ] handlers/auth/auth_Logout_test.go
- [ ] handlers/auth/auth_pipeline_test.go
- [ ] handlers/blog_likes/blog_likes_GenerateVisitId_test.go
- [ ] handlers/blog_likes/blog_likes_IsBlogLiked_test.go
- [ ] handlers/blog_likes/blog_likes_CreateBlogLike_test.go
- [ ] handlers/blog_likes/blog_likes_DeleteBlogLike_test.go

### Handler層（リファクタリング）

- [ ] handlers/auth/auth_test.go（Login失敗ケース・不正ボディ追加）
- [ ] handlers/blogs/test/blogs_FetchBlogs_test.go（テーブルドリブン化）
- [ ] handlers/blogs/test/blogs_CreateBlog_test.go（準正常系・異常系追加）
- [ ] handlers/blog_likes/blog_likes_FetchBlogLikesByVisitId_test.go（テーブルドリブン化）
- [ ] handlers/blog_comments/blog_comments_CreateComment_test.go（テーブルドリブン化）

### Service層（リファクタリング）

- [ ] services/blogs/test/blogs_CreateBlog_test.go（テーブルドリブン化）
- [ ] services/blogs/test/blogs_FetchBlogs_test.go（テーブルドリブン化）

### Repository層（新規）

- [ ] repositories/blog_users/blog_users_UpdateBlogUsers_test.go
- [ ] repositories/blog_comments/blog_comments_CreateComment_test.go
