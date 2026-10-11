package services_blog_users

import "backend/apperrors"

// ユーザーのサービスが返すエラー。呼び出し側は errors.Is で判定する。
var (
	ErrEmailAndPasswordRequired = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Email and password are required")
	ErrInvalidEmailFormat       = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid email format")
	ErrIDRequired               = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Id is required")
	ErrNameRequired             = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Name is required")
	ErrEmailRequired            = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Email is required")
	ErrPasswordRequired         = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Password is required")
	ErrNewPasswordRequired      = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "New password is required")
	// ErrPasswordTooLong は bcrypt が扱える長さ（72 バイト）を超えるパスワードのエラー。
	ErrPasswordTooLong        = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Password is too long")
	ErrInvalidCurrentPassword = apperrors.New(apperrors.KindInvalid, apperrors.CodeInvalidCurrentPassword, "Invalid current password")
	// ErrInvalidCredentials はログイン時にユーザー不在・パスワード不一致のいずれかであることを示す。
	// どちらであるかを呼び出し側に区別させない（メールアドレスの登録有無を推測させないため）。
	ErrInvalidCredentials = apperrors.New(apperrors.KindUnauthorized, apperrors.CodeInvalidCredentials, "Invalid credentials")
	ErrUserNotFound       = apperrors.New(apperrors.KindNotFound, apperrors.CodeUserNotFound, "User not found")
)
