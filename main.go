package main

import (
	"log"
	"os"

	"media-content-api/config"
	"media-content-api/models"
	"media-content-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Movie{}, &models.Series{})

	router := routes.SetupRouter()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
