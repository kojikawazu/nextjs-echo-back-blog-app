package services_blog_comments

import "backend/apperrors"

// コメントのサービスが返すエラー。呼び出し側は errors.Is で判定する。
var (
	ErrInvalidBlogID    = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid blogId")
	ErrInvalidGuestUser = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid guestUser")
	ErrInvalidComment   = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid comment")
)
