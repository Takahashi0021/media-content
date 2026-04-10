package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "postgres"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Attempting to create database...")

		dsnWithoutDB := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Asia/Almaty",
			host, user, password, port)

		db, err := gorm.Open(postgres.Open(dsnWithoutDB), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to PostgreSQL: " + err.Error())
		}

		createDBQuery := fmt.Sprintf("CREATE DATABASE %s;", dbname)
		if err := db.Exec(createDBQuery).Error; err != nil {
			panic("Failed to create database: " + err.Error())
		}

		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to database after creation: " + err.Error())
		}

		log.Printf("Database %s created successfully", dbname)
	}

	DB = database
	log.Println("Database connected successfully")
}
