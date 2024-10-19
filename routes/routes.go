package routes

import (
	utils_cookie "backend/utils/cookie"

	handlers_auth "backend/handlers/auth"
	handlers_blogs "backend/handlers/blogs"
	handlers_blogs_likes "backend/handlers/blogs_likes"
	handlers_comments "backend/handlers/comments"
	handlers_users "backend/handlers/users"

	repositories_blogs "backend/repositories/blogs"
	repositories_blogs_likes "backend/repositories/blogs_likes"
	repositories_comments "backend/repositories/comments"
	repositories_users "backend/repositories/users"

	services_auth "backend/services/auth"
	services_blogs "backend/services/blogs"
	services_blogs_likes "backend/services/blogs_likes"
	services_comments "backend/services/comments"
	services_users "backend/services/users"

	"net/http"

	"github.com/labstack/echo/v4"
)

// ルーティングを設定する関数
func SetupRoutes(e *echo.Echo) {
	// ヘルスチェックエンドポイントの追加
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Service is running")
	})

	// RepositoryとServiceとHandlerの初期化
	cookieUtils := utils_cookie.NewCookieUtils()

	userRepository := repositories_users.NewUserRepository()
	blogRepository := repositories_blogs.NewBlogRepository()
	BlogLikeRepository := repositories_blogs_likes.NewBlogLikeRepository()
	commentRepository := repositories_comments.NewCommentRepository()

	authService := services_auth.NewAuthService()
	userService := services_users.NewUserService(userRepository)
	blogService := services_blogs.NewBlogService(blogRepository)
	blogLikeService := services_blogs_likes.NewBlogLikeService(BlogLikeRepository)
	commentService := services_comments.NewCommentService(commentRepository)

	authHandler := handlers_auth.NewAuthHandler(userService, authService)
	UserHandler := handlers_users.NewUserHandler(userService, cookieUtils)
	BlogHandler := handlers_blogs.NewBlogHandler(blogService, cookieUtils)
	BlogLikeHandler := handlers_blogs_likes.NewBlogLikeHandler(blogLikeService, cookieUtils)
	CommentHandler := handlers_comments.NewCommentHandler(commentService)

	// APIエンドポイントの設定
	api := e.Group("/api")
	{
		// ユーザー関連のエンドポイント
		users := api.Group("/users")
		{
			users.POST("/login", authHandler.Login)
			users.GET("/auth-check", authHandler.CheckAuth)
			users.POST("/logout", authHandler.Logout)

			users.GET("/detail", UserHandler.FetchUser)
			users.PUT("/update", UserHandler.UpdateUser)
		}
		// ブログ関連のエンドポイント
		blogs := api.Group("/blogs")
		{
			blogs.GET("", BlogHandler.FetchBlogs)
			blogs.GET("/user/:userId", BlogHandler.FetchBlogsByUserId)
			blogs.GET("/detail/:id", BlogHandler.FetchBlogById)
			blogs.GET("/categories", BlogHandler.FetchBlogCategories)
			blogs.POST("/create", BlogHandler.CreateBlog)
			blogs.PUT("/update/:id", BlogHandler.UpdateBlog)
			blogs.DELETE("/delete/:id", BlogHandler.DeleteBlog)
		}
		// ブログいいね関連のエンドポイント
		blogLikes := api.Group("/blog-likes")
		{
			blogLikes.GET("", BlogLikeHandler.FetchBlogLikesByVisitId)
			blogLikes.GET("/generate-visit-id", BlogLikeHandler.GenerateVisitorId)
			blogLikes.GET("/is-liked/:blogId", BlogLikeHandler.IsBlogLiked)
			blogLikes.POST("/create/:blogId", BlogLikeHandler.CreateBlogLike)
			blogLikes.DELETE("/delete/:blogId", BlogLikeHandler.DeleteBlogLike)
		}
		// コメント関連のエンドポイント
		comments := api.Group("/comments")
		{
			comments.GET("/blog/:blogId", CommentHandler.FetchCommentsByBlogId)
			comments.POST("/create", CommentHandler.CreateComment)
		}
	}
}
