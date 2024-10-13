package routes

import (
	handlers_auth "backend/handlers/auth"
	handlers_blogs "backend/handlers/blogs"

	repositories_blogs "backend/repositories/blogs"
	repositories_users "backend/repositories/users"

	services_auth "backend/services/auth"
	services_blogs "backend/services/blogs"
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
	userRepository := repositories_users.NewUserRepository()
	blogRepository := repositories_blogs.NewBlogRepository()

	authService := services_auth.NewAuthService()
	userService := services_users.NewUserService(userRepository)
	blogService := services_blogs.NewBlogService(blogRepository)

	authHandler := handlers_auth.NewAuthHandler(userService, authService)
	BlogHandler := handlers_blogs.NewBlogHandler(blogService)

	// APIエンドポイントの設定
	api := e.Group("/api")
	{
		// ユーザー関連のエンドポイント
		users := api.Group("/users")
		{
			users.POST("/login", authHandler.Login)
			users.GET("/auth-check", authHandler.CheckAuth)
			users.POST("/logout", authHandler.Logout)
		}
		// ブログ関連のエンドポイント
		blogs := api.Group("/blogs")
		{
			blogs.GET("", BlogHandler.FetchBlogs)
			blogs.GET("/user", BlogHandler.FetchBlogByUserId)
		}
	}
}
