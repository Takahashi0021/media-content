package controllers

import (
	"media-content-api/config"
	"media-content-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMovies(c *gin.Context) {
	var movies []models.Movie
	config.DB.Find(&movies)
	c.JSON(http.StatusOK, gin.H{"data": movies})
}

func GetMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": movie})
}

func CreateMovie(c *gin.Context) {
	var input models.Movie
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func UpdateMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	var input models.Movie
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&movie).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": movie})
}

func DeleteMovie(c *gin.Context) {
	var movie models.Movie
	if err := config.DB.First(&movie, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	config.DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
