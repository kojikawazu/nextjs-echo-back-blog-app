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
