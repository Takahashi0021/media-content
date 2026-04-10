package main

import (
	"log"
	"media-content-api/config"
	"media-content-api/models"
	"media-content-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config.ConnectDatabase()

	err = config.DB.AutoMigrate(&models.Movie{}, &models.Series{}, &models.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	log.Println("Database migration completed - tables created")

	router := routes.SetupRouter()

	log.Println("Server starting on port 8080...")
	router.Run(":8080")
}
