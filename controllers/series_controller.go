package controllers

import (
	"media-content-api/config"
	"media-content-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateSeriesInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Seasons     int    `json:"seasons"`
	Episodes    int    `json:"episodes"`
	Genre       string `json:"genre"`
	Creator     string `json:"creator"`
	PosterURL   string `json:"poster_url"`
}

type UpdateSeriesInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Seasons     int    `json:"seasons"`
	Episodes    int    `json:"episodes"`
	Genre       string `json:"genre"`
	Creator     string `json:"creator"`
	PosterURL   string `json:"poster_url"`
}

func GetSeries(c *gin.Context) {
	var series []models.Series
	config.DB.Find(&series)

	result := make([]map[string]interface{}, 0)
	for _, s := range series {
		var avgRating float64
		config.DB.Model(&models.Review{}).
			Where("series_id = ?", s.ID).
			Select("COALESCE(AVG(rating), 0)").
			Row().Scan(&avgRating)

		result = append(result, map[string]interface{}{
			"id":           s.ID,
			"title":        s.Title,
			"description":  s.Description,
			"release_year": s.ReleaseYear,
			"seasons":      s.Seasons,
			"episodes":     s.Episodes,
			"rating":       avgRating,
			"genre":        s.Genre,
			"creator":      s.Creator,
			"poster_url":   s.PosterURL,
			"created_by":   s.CreatedBy,
			"created_at":   s.CreatedAt,
			"updated_at":   s.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetSeriesByID(c *gin.Context) {
	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}

	var avgRating float64
	config.DB.Model(&models.Review{}).
		Where("series_id = ?", series.ID).
		Select("COALESCE(AVG(rating), 0)").
		Row().Scan(&avgRating)

	c.JSON(http.StatusOK, gin.H{"data": map[string]interface{}{
		"id":           series.ID,
		"title":        series.Title,
		"description":  series.Description,
		"release_year": series.ReleaseYear,
		"seasons":      series.Seasons,
		"episodes":     series.Episodes,
		"rating":       avgRating,
		"genre":        series.Genre,
		"creator":      series.Creator,
		"poster_url":   series.PosterURL,
		"created_by":   series.CreatedBy,
		"created_at":   series.CreatedAt,
		"updated_at":   series.UpdatedAt,
	}})
}

func CreateSeries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input CreateSeriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	series := models.Series{
		Title:       input.Title,
		Description: input.Description,
		ReleaseYear: input.ReleaseYear,
		Seasons:     input.Seasons,
		Episodes:    input.Episodes,
		Genre:       input.Genre,
		Creator:     input.Creator,
		PosterURL:   input.PosterURL,
		CreatedBy:   userID.(uint),
	}

	config.DB.Create(&series)
	c.JSON(http.StatusCreated, gin.H{"data": series})
}

func UpdateSeries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	role, _ := c.Get("role")

	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}

	if series.CreatedBy != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own series"})
		return
	}

	var input UpdateSeriesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := models.Series{}
	if input.Title != "" {
		updates.Title = input.Title
	}
	if input.Description != "" {
		updates.Description = input.Description
	}
	if input.ReleaseYear != 0 {
		updates.ReleaseYear = input.ReleaseYear
	}
	if input.Seasons != 0 {
		updates.Seasons = input.Seasons
	}
	if input.Episodes != 0 {
		updates.Episodes = input.Episodes
	}
	if input.Genre != "" {
		updates.Genre = input.Genre
	}
	if input.Creator != "" {
		updates.Creator = input.Creator
	}
	if input.PosterURL != "" {
		updates.PosterURL = input.PosterURL
	}

	config.DB.Model(&series).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"data": series})
}

func DeleteSeries(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	role, _ := c.Get("role")

	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}

	if series.CreatedBy != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own series"})
		return
	}

	config.DB.Delete(&series)
	c.JSON(http.StatusOK, gin.H{"message": "Series deleted successfully"})
}
