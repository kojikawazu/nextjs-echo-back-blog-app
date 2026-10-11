// Package dto は API のリクエスト/レスポンス構造体を定義する。
// ドメインモデル・DB 構造体をそのまま返さず、公開してよいフィールドだけをここで定義する。
package dto

import (
	"backend/models"
	"time"
)

// UserResponse はユーザー情報のレスポンス。パスワードは含めない。
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse はユーザーモデルをレスポンス DTO に変換する。
//
// 引数:
//   - user: 変換元のユーザーモデル
//
// 戻り値:
//   - UserResponse: 公開してよいフィールドのみを持つレスポンス
func NewUserResponse(user *models.BlogUsersData) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
