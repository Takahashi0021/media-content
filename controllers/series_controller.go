package controllers

import (
	"media-content-api/config"
	"media-content-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSeries(c *gin.Context) {
	var series []models.Series
	config.DB.Find(&series)
	c.JSON(http.StatusOK, gin.H{"data": series})
}

func GetSeriesByID(c *gin.Context) {
	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": series})
}

func CreateSeries(c *gin.Context) {
	var input models.Series
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Create(&input)
	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func UpdateSeries(c *gin.Context) {
	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}

	var input models.Series
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&series).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": series})
}

func DeleteSeries(c *gin.Context) {
	var series models.Series
	if err := config.DB.First(&series, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Series not found"})
		return
	}

	config.DB.Delete(&series)
	c.JSON(http.StatusOK, gin.H{"data": true})
}
