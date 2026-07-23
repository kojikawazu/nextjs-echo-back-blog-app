//go:build e2e

package e2e

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 正常系: ログイン → 認証確認 → ログアウト → 認証確認(未認証) の一連フロー。
func TestE2E_AuthFlow(t *testing.T) {
	client := newClient(t)

	// ログイン → 200 + token Cookie
	loginResp := doJSON(t, client, http.MethodPost, "/api/users/login", map[string]string{
		"email":    os.Getenv("TEST_USER_EMAIL"),
		"password": os.Getenv("TEST_USER_PASSWD"),
	})
	defer loginResp.Body.Close()
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// 認証確認 → 200 + ユーザー情報（Cookie が引き継がれる）
	checkResp := doJSON(t, client, http.MethodGet, "/api/users/auth-check", nil)
	assert.Equal(t, http.StatusOK, checkResp.StatusCode)
	var checkBody map[string]string
	decodeBody(t, checkResp, &checkBody)
	assert.Equal(t, os.Getenv("TEST_USER_EMAIL"), checkBody["email"])
	assert.Equal(t, os.Getenv("TEST_USER_ID"), checkBody["user_id"])

	// ログアウト → 200
	logoutResp := doJSON(t, client, http.MethodPost, "/api/users/logout", nil)
	defer logoutResp.Body.Close()
	assert.Equal(t, http.StatusOK, logoutResp.StatusCode)

	// ログアウト後の認証確認 → 401
	afterResp := doJSON(t, client, http.MethodGet, "/api/users/auth-check", nil)
	defer afterResp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, afterResp.StatusCode)
}

// 準正常系: 誤ったパスワードでのログインは 401。
func TestE2E_Login_WrongPassword(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodPost, "/api/users/login", map[string]string{
		"email":    os.Getenv("TEST_USER_EMAIL"),
		"password": "wrong-password",
	})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// 準正常系: email/password 欠落は 400。
func TestE2E_Login_MissingFields(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodPost, "/api/users/login", map[string]string{})
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// 準正常系: Cookie 無しでの認証確認は 401。
func TestE2E_CheckAuth_NoCookie(t *testing.T) {
	client := newClient(t)

	resp := doJSON(t, client, http.MethodGet, "/api/users/auth-check", nil)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
