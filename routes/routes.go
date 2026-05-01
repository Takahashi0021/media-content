package routes

import (
	"media-content-api/controllers"
	"media-content-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Public routes (no authentication required)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// External API routes (public - no auth needed for external API calls)
	external := router.Group("/api/external")
	{
		external.GET("/movie/:title", controllers.GetExternalMovie)
		external.GET("/search", controllers.SearchExternalMovies)
	}

	// Protected routes (authentication required)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/user", controllers.GetCurrentUser)
		protected.GET("/users", middleware.AdminMiddleware(), controllers.GetAllUsers)

		// Movie routes (protected)
		protected.GET("/movies", controllers.GetMovies)
		protected.GET("/movies/:id", controllers.GetMovie)
		protected.POST("/movies", controllers.CreateMovie)
		protected.PUT("/movies/:id", controllers.UpdateMovie)
		protected.DELETE("/movies/:id", middleware.AdminMiddleware(), controllers.DeleteMovie)

		// Series routes (protected)
		protected.GET("/series", controllers.GetSeries)
		protected.GET("/series/:id", controllers.GetSeriesByID)
		protected.POST("/series", controllers.CreateSeries)
		protected.PUT("/series/:id", controllers.UpdateSeries)
		protected.DELETE("/series/:id", middleware.AdminMiddleware(), controllers.DeleteSeries)
	}

	return router
}
