package handlers_blogs

import (
	services_blogs "backend/services/blogs"
	utils_cookie "backend/utils/cookie"
)

// BlogHandler はブログ系エンドポイントのハンドラ。
type BlogHandler struct {
	BlogService services_blogs.BlogService
	CookieUtils utils_cookie.CookieUtils
}

// NewBlogHandler は BlogHandler を生成する。
//
// 引数:
//   - blogService: ブログデータの取得/作成/更新/削除に用いるブログサービス
//   - cookieUtils: 認証トークンの取得/解析に用いるクッキーユーティリティ
//
// 戻り値:
//   - *BlogHandler: 生成したブログハンドラ
func NewBlogHandler(blogService services_blogs.BlogService, cookieUtils utils_cookie.CookieUtils) *BlogHandler {
	return &BlogHandler{
		BlogService: blogService,
		CookieUtils: cookieUtils,
	}
}
