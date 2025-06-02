package routes

import (
	"github.com/gin-gonic/gin"
	"sinkedin/handlers"
	"sinkedin/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// User routes
	userRoutes := r.Group("/api/users")
	{
		userRoutes.POST("/register", handlers.RegisterUser)
		userRoutes.POST("/login", handlers.LoginUser)
		userRoutes.GET("/:username", handlers.GetUserProfile)
		userRoutes.PUT("/:username", middleware.AuthMiddleware(), handlers.UpdateUserProfile)
		userRoutes.DELETE("/:username", middleware.AuthMiddleware(), handlers.DeleteUser)
	}

	// Post routes
	postRoutes := r.Group("/api/posts", middleware.AuthMiddleware())
	{
		postRoutes.POST("/", handlers.CreatePost)
		postRoutes.GET("/", handlers.GetPosts)
		postRoutes.GET("/:id", handlers.GetPost)
		postRoutes.PUT("/:id", handlers.UpdatePost)
		postRoutes.DELETE("/:id", handlers.DeletePost)
	}

	// Comment routes
	commentRoutes := r.Group("/api/comments", middleware.AuthMiddleware())
	{
		commentRoutes.POST("/", handlers.CreateComment)
		commentRoutes.GET("/post/:postId", handlers.GetPostComments)
		commentRoutes.GET("/:id", handlers.GetComment)
		commentRoutes.PUT("/:id", handlers.UpdateComment)
		commentRoutes.DELETE("/:id", handlers.DeleteComment)
	}

	// Like routes
	likeRoutes := r.Group("/api/likes", middleware.AuthMiddleware())
	{
		likeRoutes.POST("/:type/:id", handlers.ToggleLike)
		likeRoutes.GET("/:type/:id", handlers.GetLikes)
	}

	// Follow routes
	followRoutes := r.Group("/api/follow", middleware.AuthMiddleware())
	{
		followRoutes.POST("/:username", handlers.ToggleFollow)
		followRoutes.GET("/followers/:username", handlers.GetFollowers)
		followRoutes.GET("/following/:username", handlers.GetFollowing)
	}

	// Hashtag routes
	hashtagRoutes := r.Group("/api/hashtags")
	{
		hashtagRoutes.GET("/trending", handlers.GetTrendingHashtags)
		hashtagRoutes.GET("/:name/posts", handlers.GetHashtagPosts)
	}
}
