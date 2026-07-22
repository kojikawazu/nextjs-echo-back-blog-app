package routes

import (
	utils_cookie "backend/utils/cookie"

	handlers_auth "backend/handlers/auth"
	handlers_comments "backend/handlers/blog_comments"
	handlers_blogs_likes "backend/handlers/blog_likes"
	handlers_blog_users "backend/handlers/blog_users"
	handlers_blogs "backend/handlers/blogs"

	repositories_comments "backend/repositories/blog_comments"
	repositories_blogs_likes "backend/repositories/blog_likes"
	repositories_blog_users "backend/repositories/blog_users"
	repositories_blogs "backend/repositories/blogs"

	services_auth "backend/services/auth"
	services_comments "backend/services/blog_comments"
	services_blogs_likes "backend/services/blog_likes"
	services_blog_users "backend/services/blog_users"
	services_blogs "backend/services/blogs"

	"net/http"

	"github.com/labstack/echo/v4"
)

// SetupRoutes は依存を初期化し、API エンドポイントのルーティングを設定する。
//
// 引数:
//   - e: ルーティングを設定する対象の Echo インスタンス
func SetupRoutes(e *echo.Echo) {
	// ヘルスチェックエンドポイントの追加
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is running")
	})

	// RepositoryとServiceとHandlerの初期化
	cookieUtils := utils_cookie.NewCookieUtils()

	userRepository := repositories_blog_users.NewBlogUsersRepository()
	blogRepository := repositories_blogs.NewBlogRepository()
	BlogLikeRepository := repositories_blogs_likes.NewBlogLikeRepository()
	commentRepository := repositories_comments.NewCommentRepository()

	authService := services_auth.NewAuthService()
	userService := services_blog_users.NewUserService(userRepository)
	blogService := services_blogs.NewBlogService(blogRepository)
	blogLikeService := services_blogs_likes.NewBlogLikeService(BlogLikeRepository)
	commentService := services_comments.NewCommentService(commentRepository)

	authHandler := handlers_auth.NewAuthHandler(userService, authService)
	blogUsersHandler := handlers_blog_users.NewBlogUsersHandler(userService, cookieUtils)
	blogHandler := handlers_blogs.NewBlogHandler(blogService, cookieUtils)
	blogLikeHandler := handlers_blogs_likes.NewBlogLikeHandler(blogLikeService, cookieUtils)
	commentHandler := handlers_comments.NewCommentHandler(commentService)

	// APIエンドポイントの設定
	api := e.Group("/api")
	{
		// ユーザー関連のエンドポイント
		blogUsers := api.Group("/users")
		{
			blogUsers.POST("/login", authHandler.Login)
			blogUsers.GET("/auth-check", authHandler.CheckAuth)
			blogUsers.POST("/logout", authHandler.Logout)

			blogUsers.GET("/detail", blogUsersHandler.FetchBlogUsers)
			blogUsers.PUT("/update", blogUsersHandler.UpdateBlogUsers)
		}
		// ブログ関連のエンドポイント
		blogs := api.Group("/blogs")
		{
			blogs.GET("", blogHandler.FetchBlogs)
			blogs.GET("/user/:userId", blogHandler.FetchBlogsByUserId)
			blogs.GET("/detail/:id", blogHandler.FetchBlogById)
			blogs.GET("/categories", blogHandler.FetchBlogCategories)
			blogs.GET("/tags", blogHandler.FetchBlogTags)
			blogs.GET("/popular/:count", blogHandler.FetchBlogPopular)
			blogs.POST("/create", blogHandler.CreateBlog)
			blogs.PUT("/update/:id", blogHandler.UpdateBlog)
			blogs.DELETE("/delete/:id", blogHandler.DeleteBlog)
		}
		// ブログいいね関連のエンドポイント
		blogLikes := api.Group("/blog-likes")
		{
			blogLikes.GET("", blogLikeHandler.FetchBlogLikesByVisitId)
			blogLikes.GET("/generate-visit-id", blogLikeHandler.GenerateVisitorId)
			blogLikes.GET("/is-liked/:blogId", blogLikeHandler.IsBlogLiked)
			blogLikes.POST("/create/:blogId", blogLikeHandler.CreateBlogLike)
			blogLikes.DELETE("/delete/:blogId", blogLikeHandler.DeleteBlogLike)
		}
		// コメント関連のエンドポイント
		comments := api.Group("/comments")
		{
			comments.GET("/blog/:blogId", commentHandler.FetchCommentsByBlogId)
			comments.POST("/create", commentHandler.CreateComment)
		}
	}
}
