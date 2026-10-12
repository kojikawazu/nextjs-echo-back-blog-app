package utils_cookie

import (
	"backend/config"
	"backend/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// signTestToken は任意のアルゴリズム・鍵・クレームでトークンを生成する（不正なトークンを作るためのヘルパー）。
func signTestToken(t *testing.T, method jwt.SigningMethod, key any, claims jwt.Claims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(method, claims).SignedString(key)
	require.NoError(t, err)
	return token
}

// authClaims は aud・exp を指定した認証トークンのクレームを返す。
func authClaims(userID string, audience []string, expiresAt *jwt.NumericDate) *models.Claims {
	return &models.Claims{
		UserID:           userID,
		Email:            "test@example.com",
		Username:         "Test User",
		RegisteredClaims: jwt.RegisteredClaims{Audience: audience, ExpiresAt: expiresAt},
	}
}

// visitorClaims は aud・exp を指定した訪問者トークンのクレームを返す。
func visitorClaims(visitID string, audience []string, expiresAt *jwt.NumericDate) *models.ClaimsVisitId {
	return &models.ClaimsVisitId{
		VisitId:          visitID,
		RegisteredClaims: jwt.RegisteredClaims{Audience: audience, ExpiresAt: expiresAt},
	}
}

// TestNewAuthToken は発行した認証トークンが aud=auth・指定の exp・ユーザー情報を持つことを検証する（正常系）。
func TestNewAuthToken(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour).Truncate(time.Second)

	token, err := NewAuthToken(&models.BlogUsersData{ID: "user-1", Email: "a@example.com", Name: "Alice"}, expiresAt)
	require.NoError(t, err)

	claims, err := ParseAuthToken(token)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, "a@example.com", claims.Email)
	assert.Equal(t, "Alice", claims.Username)
	assert.Equal(t, jwt.ClaimStrings{AudienceAuth}, claims.Audience)
	assert.Equal(t, expiresAt, claims.ExpiresAt.Time)
}

// TestNewVisitorToken は発行した訪問者トークンが aud=visitor・指定の exp・訪問者IDを持つことを検証する（正常系）。
func TestNewVisitorToken(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour).Truncate(time.Second)

	token, err := NewVisitorToken("visitor-1", expiresAt)
	require.NoError(t, err)

	claims, err := ParseVisitorToken(token)
	require.NoError(t, err)
	assert.Equal(t, "visitor-1", claims.VisitId)
	assert.Equal(t, jwt.ClaimStrings{AudienceVisitor}, claims.Audience)
	assert.Equal(t, expiresAt, claims.ExpiresAt.Time)
}

// TestParseAuthToken は認証トークンの検証で、不正・期限切れ・種類違いのトークンを拒否することを検証する。
func TestParseAuthToken(t *testing.T) {
	future := jwt.NewNumericDate(time.Now().Add(time.Hour))
	past := jwt.NewNumericDate(time.Now().Add(-time.Hour))
	auth := []string{AudienceAuth}

	visitorToken, err := NewVisitorToken("visitor-1", time.Now().Add(time.Hour))
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			// 正常系
			name:  "valid",
			token: signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, authClaims("user-1", auth, future)),
		},
		{
			// 準正常系: 形式不正
			name:    "malformed",
			token:   "not-a-jwt",
			wantErr: jwt.ErrTokenMalformed,
		},
		{
			// 準正常系: 期限切れ
			name:    "expired",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, authClaims("user-1", auth, past)),
			wantErr: jwt.ErrTokenExpired,
		},
		{
			// 準正常系: exp が無い（以前は CheckAuth だけが受け入れていた。#156）
			name:    "missing_exp",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, authClaims("user-1", auth, nil)),
			wantErr: jwt.ErrTokenRequiredClaimMissing,
		},
		{
			// 準正常系: 別の鍵で署名（改ざん）
			name:    "wrong_key",
			token:   signTestToken(t, jwt.SigningMethodHS256, []byte("attacker-key"), authClaims("user-1", auth, future)),
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			// 準正常系: alg=none（署名なし）
			name:    "alg_none",
			token:   signTestToken(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, authClaims("user-1", auth, future)),
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			// 準正常系: 同じ鍵でも HS256 以外のアルゴリズムは受け付けない
			name:    "alg_hs512",
			token:   signTestToken(t, jwt.SigningMethodHS512, config.JwtKey, authClaims("user-1", auth, future)),
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			// 準正常系: 誰でも発行できる訪問者トークンを認証トークンとして使えない（#156）
			name:    "visitor_token",
			token:   visitorToken,
			wantErr: jwt.ErrTokenInvalidAudience,
		},
		{
			// 準正常系: aud を持たない旧形式のトークン
			name:    "missing_aud",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, authClaims("user-1", nil, future)),
			wantErr: jwt.ErrTokenRequiredClaimMissing,
		},
		{
			// 準正常系: user_id が空
			name:    "empty_user_id",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, authClaims("", auth, future)),
			wantErr: errEmptyUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseAuthToken(tt.token)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, claims)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "user-1", claims.UserID)
		})
	}
}

// TestParseVisitorToken は訪問者トークンの検証で、不正・期限切れ・種類違いのトークンを拒否することを検証する。
func TestParseVisitorToken(t *testing.T) {
	future := jwt.NewNumericDate(time.Now().Add(time.Hour))
	past := jwt.NewNumericDate(time.Now().Add(-time.Hour))
	visitor := []string{AudienceVisitor}

	authToken, err := NewAuthToken(&models.BlogUsersData{ID: "user-1"}, time.Now().Add(time.Hour))
	require.NoError(t, err)

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			// 正常系
			name:  "valid",
			token: signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, visitorClaims("visitor-1", visitor, future)),
		},
		{
			// 準正常系: 期限切れ
			name:    "expired",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, visitorClaims("visitor-1", visitor, past)),
			wantErr: jwt.ErrTokenExpired,
		},
		{
			// 準正常系: exp が無い
			name:    "missing_exp",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, visitorClaims("visitor-1", visitor, nil)),
			wantErr: jwt.ErrTokenRequiredClaimMissing,
		},
		{
			// 準正常系: alg=none（署名なし）
			name:    "alg_none",
			token:   signTestToken(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, visitorClaims("visitor-1", visitor, future)),
			wantErr: jwt.ErrTokenSignatureInvalid,
		},
		{
			// 準正常系: 認証トークンを訪問者トークンとして使えない
			name:    "auth_token",
			token:   authToken,
			wantErr: jwt.ErrTokenInvalidAudience,
		},
		{
			// 準正常系: visit_id が空
			name:    "empty_visit_id",
			token:   signTestToken(t, jwt.SigningMethodHS256, config.JwtKey, visitorClaims("", visitor, future)),
			wantErr: errEmptyVisitID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseVisitorToken(tt.token)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, claims)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, "visitor-1", claims.VisitId)
		})
	}
}
