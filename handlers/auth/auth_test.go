package handlers_auth

import (
	"backend/models"
	services_auth "backend/services/auth"
	services_users "backend/services/users"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// .envファイルの読み込み
	err := godotenv.Load("../../.env.test")
	if err != nil {
		panic("Error loading ../../.env.test file")
	}

	// テストを実行
	code := m.Run()

	// テスト終了後に終了コードで終了
	os.Exit(code)
}

func TestHandler_Login_Success(t *testing.T) {
	// Echo インスタンスの作成
	e := echo.New()

	// リクエストボディの準備
	reqBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// HTTP リクエストとレスポンスの作成
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// モックの作成と期待値の設定
	mockAuthService := new(services_auth.MockAuthService)
	mockUserService := new(services_users.MockUserService)

	mockAuthService.On("Login", "test@example.com", "password123").Return(nil)

	mockUserService.On("FetchUserByEmailAndPassword", "test@example.com", "password123").Return(&models.UserData{
		ID:    "user123",
		Email: "test@example.com",
		Name:  "Test User",
	}, nil)

	// AuthHandler の作成
	handler := NewAuthHandler(mockUserService, mockAuthService)

	// テスト対象のハンドラー関数の呼び出し
	handler.Login(c)

	assert.Equal(t, http.StatusOK, rec.Code)
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	assert.Equal(t, "Login successful", response["message"])

	// クッキーの確認
	cookies := rec.Result().Cookies()
	assert.NotEmpty(t, cookies)

	// 特定のクッキーを確認
	var tokenCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			tokenCookie = c
			break
		}
	}
	assert.NotNil(t, tokenCookie)
	assert.NotEmpty(t, tokenCookie.Value)

	// 期待値のアサーション
	mockAuthService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
