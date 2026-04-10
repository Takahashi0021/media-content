package routes

import (
	"media-content-api/controllers"
	"media-content-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

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
	}

	return router
}
