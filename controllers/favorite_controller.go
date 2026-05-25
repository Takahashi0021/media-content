package controllers

import (
	"net/http"
	"strconv"

	"media-content-api/config"
	"media-content-api/models"
	"media-content-api/services"

	"github.com/gin-gonic/gin"
)

var favoriteService = services.NewFavoriteService()

type AddFavoriteInput struct {
	MovieID  *uint `json:"movie_id"`
	SeriesID *uint `json:"series_id"`
}

func AddFavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input AddFavoriteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.MovieID == nil && input.SeriesID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either movie_id or series_id is required"})
		return
	}

	var existsFavorite models.Favorite
	query := config.DB.Where("user_id = ?", userID)
	if input.MovieID != nil {
		query = query.Where("movie_id = ?", input.MovieID)
	}
	if input.SeriesID != nil {
		query = query.Where("series_id = ?", input.SeriesID)
	}

	if query.First(&existsFavorite).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already in favorites"})
		return
	}

	if err := favoriteService.AddFavorite(userID.(uint), input.MovieID, input.SeriesID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Favorite added successfully"})
}

func RemoveFavorite(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	favoriteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid favorite ID"})
		return
	}

	if err := favoriteService.RemoveFavorite(userID.(uint), uint(favoriteID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

func GetUserFavorites(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	favorites, err := favoriteService.GetUserFavorites(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": favorites})
}

func GetMovieFavoritesCount(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	count, err := favoriteService.GetMovieFavoritesCount(uint(movieID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"movie_id": movieID, "favorites_count": count})
}

func GetSeriesFavoritesCount(c *gin.Context) {
	seriesID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid series ID"})
		return
	}

	count, err := favoriteService.GetSeriesFavoritesCount(uint(seriesID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"series_id": seriesID, "favorites_count": count})
}
