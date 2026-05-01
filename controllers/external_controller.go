package controllers

import (
	"fmt"
	"media-content-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var omdbService = services.NewOMDbService()

// GET /api/external/movie/:title
func GetExternalMovie(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie title is required"})
		return
	}

	movie, err := omdbService.GetMovieByTitle(title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

// GET /api/external/search?q=query&page=1
func SearchExternalMovies(c *gin.Context) {
	query := c.Query("q")
	page := c.DefaultQuery("page", "1")

	var pageNum int
	if _, err := fmt.Sscanf(page, "%d", &pageNum); err != nil {
		pageNum = 1
	}

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	results, err := omdbService.SearchMovies(query, pageNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "query": query, "page": pageNum})
}
