package services_auth

import "backend/apperrors"

// 認証のサービスが返すエラー。呼び出し側は errors.Is で判定する。
var (
	ErrEmailAndPasswordRequired = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Email and password are required")
	ErrInvalidEmailFormat       = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid email format")
)
