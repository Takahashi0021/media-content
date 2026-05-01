package main

import (
	"log"
	"os"

	"media-content-api/config"
	"media-content-api/models"
	"media-content-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Error loading .env file")
		}
	}

	config.ConnectDatabase()

	if err := config.DB.AutoMigrate(&models.User{}, &models.Movie{}, &models.Series{}); err != nil {
		log.Printf("Warning: Auto migration error: %v", err)
	}

	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
