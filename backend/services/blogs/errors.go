package services_blogs

import "backend/apperrors"

// ブログのサービスが返すエラー。呼び出し側は errors.Is で判定する。
var (
	ErrInvalidUserID      = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid userId")
	ErrInvalidID          = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid id")
	ErrInvalidTitle       = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid title")
	ErrInvalidGithubURL   = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid githubUrl")
	ErrInvalidCategory    = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid category")
	ErrInvalidDescription = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid description")
	ErrInvalidTags        = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid tags")
	ErrInvalidCount       = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "Invalid count")
	ErrBlogNotFound       = apperrors.New(apperrors.KindNotFound, apperrors.CodeBlogNotFound, "Blog not found")
)
