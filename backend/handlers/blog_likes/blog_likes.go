package handlers_blog_likes

import (
	services_blogs_likes "backend/services/blog_likes"
	utils_cookie "backend/utils/cookie"
)

// BlogLikeHandler はブログいいね系エンドポイントのハンドラ。
type BlogLikeHandler struct {
	BlogLikeService services_blogs_likes.BlogLikeService
	CookieUtils     utils_cookie.CookieUtils
}

// NewBlogLikeHandler は BlogLikeHandler を生成する。
//
// 引数:
//   - blogLikeService: いいねデータの取得/作成/削除に用いるいいねサービス
//   - cookieUtils: 訪問者トークンの取得/生成に用いるクッキーユーティリティ
//
// 戻り値:
//   - *BlogLikeHandler: 生成したいいねハンドラ
func NewBlogLikeHandler(blogLikeService services_blogs_likes.BlogLikeService, cookieUtils utils_cookie.CookieUtils) *BlogLikeHandler {
	return &BlogLikeHandler{
		BlogLikeService: blogLikeService,
		CookieUtils:     cookieUtils,
	}
}
