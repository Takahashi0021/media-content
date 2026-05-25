package routes

import (
	"media-content-api/controllers"
	"media-content-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Глобальная обработка OPTIONS запросов
	router.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.AbortWithStatus(204)
	})

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// External API routes
	external := router.Group("/api/external")
	{
		external.GET("/movie/:title", controllers.GetExternalMovie)
		external.GET("/search", controllers.SearchExternalMovies)
	}

	// Protected routes
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/user", controllers.GetCurrentUser)
		protected.GET("/users", middleware.AdminMiddleware(), controllers.GetAllUsers)

		protected.GET("/movies", controllers.GetMovies)
		protected.GET("/movies/:id", controllers.GetMovie)
		protected.POST("/movies", controllers.CreateMovie)
		protected.PUT("/movies/:id", controllers.UpdateMovie)
		protected.DELETE("/movies/:id", middleware.AdminMiddleware(), controllers.DeleteMovie)

		protected.GET("/series", controllers.GetSeries)
		protected.GET("/series/:id", controllers.GetSeriesByID)
		protected.POST("/series", controllers.CreateSeries)
		protected.PUT("/series/:id", controllers.UpdateSeries)
		protected.DELETE("/series/:id", middleware.AdminMiddleware(), controllers.DeleteSeries)

		protected.POST("/favorites", controllers.AddFavorite)
		protected.DELETE("/favorites/:id", controllers.RemoveFavorite)
		protected.GET("/favorites", controllers.GetUserFavorites)
		protected.GET("/favorites/movies/:id/count", controllers.GetMovieFavoritesCount)
		protected.GET("/favorites/series/:id/count", controllers.GetSeriesFavoritesCount)

		protected.POST("/reviews", controllers.CreateReview)
		protected.GET("/reviews", controllers.GetReviews)
		protected.PUT("/reviews/:id", controllers.UpdateReview)
		protected.DELETE("/reviews/:id", controllers.DeleteReview)
		protected.POST("/reviews/:id/like", controllers.LikeReview)
		protected.GET("/reviews/rating/average", controllers.GetAverageRating)
	}

	return router
}
