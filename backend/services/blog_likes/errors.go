package services_blog_likes

import "backend/apperrors"

// いいねのサービスが返すエラー。呼び出し側は errors.Is で判定する。
var (
	ErrVisitIDEmpty         = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "VisitId is empty")
	ErrBlogIDOrVisitIDEmpty = apperrors.New(apperrors.KindInvalid, apperrors.CodeValidation, "BlogId or VisitId is empty")
	ErrAlreadyLiked         = apperrors.New(apperrors.KindConflict, apperrors.CodeAlreadyLiked, "Blog is already liked")
)
