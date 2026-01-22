package routes

import (
	"blog/controllers"
	"blog/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.New()

	// 初始化限流器：每个IP每分钟最多100次请求
	limiter := middleware.NewFixedWindowLimiter()
	r.Use(limiter.Limit()) // 全局应用限流中间件

	//使用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(gin.Recovery())

	//创建控制器实例
	authController := &controllers.AuthContorller{}
	postController := &controllers.PostController{}
	commentController := &controllers.CommentController{}
	uploadController := &controllers.UploadController{}

	//api 路由组
	api := r.Group("/api/v1")
	{
		//认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		//需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			//用户信息
			authenticated.GET("/profile", authController.GetProfile)

			//文章相关路由
			post := authenticated.Group("/posts")
			{
				post.POST("", postController.CreatePost)
				post.PUT("/:id", postController.UpdatePost)
				post.DELETE("/:id", postController.DeletePost)
			}

			//评论相关路由
			comments := authenticated.Group("/post/:post_id/comments")
			{
				comments.POST("", commentController.CreateCommentRequest)
			}
		}
		//公开路由
		public := api.Group("")
		{
			//文章公开路由
			public.GET("/post", postController.GetPosts)
			public.GET("/post/:id", postController.GetPost)
		}

		//评论公开路由（单独分组避免路由冲突）
		comments := api.Group("/comments")
		{
			comments.GET("/post/:post_id", commentController.GetComments)
		}

		file := api.Group("/file")
		{
			file.POST("/upload", uploadController.UploadFile)
			file.GET("/download/:filename", uploadController.DownloadFile)
			file.DELETE("/delete/:filename", uploadController.DeleteFile)
		}
	}
	//健康检查
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})
	return r
}
