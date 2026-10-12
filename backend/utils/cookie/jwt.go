package utils_cookie

import (
	"backend/config"
	"backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// トークンの種類を表す aud クレームの値。
// 認証トークンと訪問者トークンは同じ鍵で署名するため、署名の検証だけでは種類を区別できない。
// 誰でも発行できる訪問者トークンを認証トークンとして使えないよう、検証時に aud を照合する（#156）。
const (
	AudienceAuth    = "auth"
	AudienceVisitor = "visitor"
)

var (
	// errEmptyUserID はクレームの user_id が空であることを示す。
	errEmptyUserID = errors.New("token has empty user_id")
	// errEmptyVisitID はクレームの visit_id が空であることを示す。
	errEmptyVisitID = errors.New("token has empty visit_id")
)

// NewAuthToken は認証トークン（token Cookie 用）を発行する。
//
// 引数:
//   - user: トークンに埋め込むユーザー情報
//   - expiresAt: トークンの有効期限
//
// 戻り値:
//   - string: 署名済みのJWTトークン文字列
//   - error: 署名に失敗した場合のエラー
func NewAuthToken(user *models.BlogUsersData, expiresAt time.Time) (string, error) {
	return signToken(&models.Claims{
		UserID:           user.ID,
		Email:            user.Email,
		Username:         user.Name,
		RegisteredClaims: registeredClaims(AudienceAuth, expiresAt),
	})
}

// NewVisitorToken は訪問者トークン（visit-id-token Cookie 用）を発行する。
//
// 引数:
//   - visitID: トークンに埋め込む訪問者ID
//   - expiresAt: トークンの有効期限
//
// 戻り値:
//   - string: 署名済みのJWTトークン文字列
//   - error: 署名に失敗した場合のエラー
func NewVisitorToken(visitID string, expiresAt time.Time) (string, error) {
	return signToken(&models.ClaimsVisitId{
		VisitId:          visitID,
		RegisteredClaims: registeredClaims(AudienceVisitor, expiresAt),
	})
}

// ParseAuthToken は認証トークンを検証し、クレームを返す。
// 署名アルゴリズム（HS256）・aud（auth）・exp（必須・期限内）・user_id（空でない）を検証する。
//
// 引数:
//   - tokenString: 検証対象のJWTトークン文字列
//
// 戻り値:
//   - *models.Claims: 検証に成功したトークンのクレーム
//   - error: 形式不正・署名不正・期限切れ・種類違い・user_id 空の場合のエラー
func ParseAuthToken(tokenString string) (*models.Claims, error) {
	claims := &models.Claims{}
	if err := parseToken(tokenString, claims, AudienceAuth); err != nil {
		return nil, err
	}
	if claims.UserID == "" {
		return nil, errEmptyUserID
	}
	return claims, nil
}

// ParseVisitorToken は訪問者トークンを検証し、クレームを返す。
// 署名アルゴリズム（HS256）・aud（visitor）・exp（必須・期限内）・visit_id（空でない）を検証する。
//
// 引数:
//   - tokenString: 検証対象のJWTトークン文字列
//
// 戻り値:
//   - *models.ClaimsVisitId: 検証に成功したトークンのクレーム
//   - error: 形式不正・署名不正・期限切れ・種類違い・visit_id 空の場合のエラー
func ParseVisitorToken(tokenString string) (*models.ClaimsVisitId, error) {
	claims := &models.ClaimsVisitId{}
	if err := parseToken(tokenString, claims, AudienceVisitor); err != nil {
		return nil, err
	}
	if claims.VisitId == "" {
		return nil, errEmptyVisitID
	}
	return claims, nil
}

// registeredClaims は aud と exp を設定した登録済みクレームを返す。
func registeredClaims(audience string, expiresAt time.Time) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{audience},
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}
}

// signToken はクレームを HS256 で署名する。
func signToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.JwtKey)
}

// parseToken はトークンを検証して claims に読み込む。
// alg を HS256 に固定し（none や別アルゴリズムを拒否）、aud の一致と exp の存在を必須にする。
func parseToken(tokenString string, claims jwt.Claims, audience string) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return config.JwtKey, nil
	},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithAudience(audience),
		jwt.WithExpirationRequired(),
	)
	return err
}
