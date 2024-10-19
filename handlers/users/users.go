package handlers_users

import (
	utils "backend/utils/log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// リクエストボディで指定されたemailとpasswordでユーザーを取得する。
// 有効なemailフォーマットかをチェックし、データベースに該当ユーザーがいない場合、404エラーを返す。
func (h *UserHandler) GetUserByEmailAndPassword(c echo.Context) error {
	utils.LogInfo(c, "Fetching user by email and password...")

	// JSONのリクエストボディからemailとpasswordを取得
	type RequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		utils.LogError(c, "Failed to bind request body: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserByEmailAndPassword(reqBody.Email, reqBody.Password)
	if err != nil {
		switch err.Error() {
		case "email and password are required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Email and password are required",
			})
		case "invalid email format":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
		case "user not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			utils.LogError(c, "Error fetching user: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch user",
			})
		}
	}

	utils.LogInfo(c, "Fetched user successfully")
	return c.JSON(http.StatusOK, user)
}

// ユーザーIDでユーザーデータを取得する
func (h *UserHandler) FetchUser(c echo.Context) error {
	utils.LogInfo(c, "Fetching user...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "token")
	if err != nil {
		utils.LogError(c, "Error getting cookie: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting cookie",
		})
	}

	// JWTトークンを解析してユーザーIDを取得
	userId, err := h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting userId from token: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting userId from token",
		})
	}

	// サービス層からユーザーデータを取得
	user, err := h.UserService.FetchUserById(userId)
	if err != nil {
		switch err.Error() {
		case "id is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Id is required",
			})
		case "user not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			utils.LogError(c, "Error fetching user: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch user",
			})
		}
	}

	// パスワードだけ抜く
	user.Password = ""

	utils.LogInfo(c, "Fetched user successfully")
	return c.JSON(http.StatusOK, user)
}

// ユーザーデータを更新する
func (h *UserHandler) UpdateUser(c echo.Context) error {
	utils.LogInfo(c, "Updating user...")

	// クッキーからJWTトークンを取得
	cookieValue, err := h.CookieUtils.GetAuthCookieValue(c, "token")
	if err != nil {
		utils.LogError(c, "Error getting cookie: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting cookie",
		})
	}

	// JWTトークンを解析してユーザーIDを取得
	userId, err := h.CookieUtils.GetUserIdFromToken(c, cookieValue)
	if err != nil {
		utils.LogError(c, "Error getting userId from token: "+err.Error())
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Error getting userId from token",
		})
	}

	// JSONのリクエストボディからname, email, passwordを取得
	type RequestBody struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		NewPassword string `json:"newPassword"`
	}

	// リクエストボディをバインド
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		utils.LogError(c, "Failed to bind request body: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// サービス層からユーザーデータを更新
	user, err := h.UserService.UpdateUser(userId, reqBody.Name, reqBody.Email, reqBody.Password, reqBody.NewPassword)
	if err != nil {
		switch err.Error() {
		case "id is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "id is required",
			})
		case "name is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Name is required",
			})
		case "email is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Email is required",
			})
		case "password is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Password is required",
			})
		case "new password is required":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "New password is required",
			})
		case "invalid email format":
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid email format",
			})
		case "user not found":
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		default:
			utils.LogError(c, "Error updating user: "+err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update user",
			})
		}
	}

	// トークンの有効期限を1時間に設定
	expirationTime := h.CookieUtils.GetAuthCookieExpirationTime()
	// トークンの再作成
	tokenString, err := h.CookieUtils.CreateToken(user)
	if err != nil {
		utils.LogError(c, "Error creating token: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create token",
		})
	}
	// トークンをクッキーに更新
	h.CookieUtils.UpdateAuthCookie(c, tokenString, expirationTime)

	utils.LogInfo(c, "Updated user successfully")
	return c.JSON(http.StatusOK, user)
}
