package controllers

import (
	"media-content-api/config"
	"media-content-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateMovieInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Duration    int    `json:"duration"`
	Genre       string `json:"genre"`
	Director    string `json:"director"`
	PosterURL   string `json:"poster_url"`
}

type UpdateMovieInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Duration    int    `json:"duration"`
	Genre       string `json:"genre"`
	Director    string `json:"director"`
	PosterURL   string `json:"poster_url"`
}

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	config.DB.Find(&movies)

	result := make([]map[string]interface{}, 0)
	for _, movie := range movies {
		var avgRating float64
		config.DB.Model(&models.Review{}).
			Where("movie_id = ?", movie.ID).
			Select("COALESCE(AVG(rating), 0)").
			Row().Scan(&avgRating)

		result = append(result, map[string]interface{}{
			"id":           movie.ID,
			"title":        movie.Title,
			"description":  movie.Description,
			"release_year": movie.ReleaseYear,
			"duration":     movie.Duration,
			"rating":       avgRating,
			"genre":        movie.Genre,
			"director":     movie.Director,
			"poster_url":   movie.PosterURL,
			"created_by":   movie.CreatedBy,
			"created_at":   movie.CreatedAt,
			"updated_at":   movie.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	var avgRating float64
	config.DB.Model(&models.Review{}).
		Where("movie_id = ?", movie.ID).
		Select("COALESCE(AVG(rating), 0)").
		Row().Scan(&avgRating)

	c.JSON(http.StatusOK, gin.H{"data": map[string]interface{}{
		"id":           movie.ID,
		"title":        movie.Title,
		"description":  movie.Description,
		"release_year": movie.ReleaseYear,
		"duration":     movie.Duration,
		"rating":       avgRating,
		"genre":        movie.Genre,
		"director":     movie.Director,
		"poster_url":   movie.PosterURL,
		"created_by":   movie.CreatedBy,
		"created_at":   movie.CreatedAt,
		"updated_at":   movie.UpdatedAt,
	}})
}

func CreateMovie(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input CreateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := models.Movie{
		Title:       input.Title,
		Description: input.Description,
		ReleaseYear: input.ReleaseYear,
		Duration:    input.Duration,
		Genre:       input.Genre,
		Director:    input.Director,
		PosterURL:   input.PosterURL,
		CreatedBy:   userID.(uint),
	}

	config.DB.Create(&movie)
	c.JSON(http.StatusCreated, gin.H{"data": movie})
}

func UpdateMovie(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	role, _ := c.Get("role")

	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	if movie.CreatedBy != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own movies"})
		return
	}

	var input UpdateMovieInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := models.Movie{}
	if input.Title != "" {
		updates.Title = input.Title
	}
	if input.Description != "" {
		updates.Description = input.Description
	}
	if input.ReleaseYear != 0 {
		updates.ReleaseYear = input.ReleaseYear
	}
	if input.Duration != 0 {
		updates.Duration = input.Duration
	}
	if input.Genre != "" {
		updates.Genre = input.Genre
	}
	if input.Director != "" {
		updates.Director = input.Director
	}
	if input.PosterURL != "" {
		updates.PosterURL = input.PosterURL
	}

	config.DB.Model(&movie).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"data": movie})
}

func DeleteMovie(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	role, _ := c.Get("role")

	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	if movie.CreatedBy != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own movies"})
		return
	}

	config.DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
