package controllers

import (
	"net/http"
	"strconv"

	"media-content-api/services"

	"github.com/gin-gonic/gin"
)

var reviewService = services.NewReviewService()

type CreateReviewInput struct {
	MovieID  *uint  `json:"movie_id"`
	SeriesID *uint  `json:"series_id"`
	Rating   int    `json:"rating" binding:"required,min=1,max=10"`
	Comment  string `json:"comment"`
}

type UpdateReviewInput struct {
	Rating  int    `json:"rating" binding:"min=1,max=10"`
	Comment string `json:"comment"`
}

func CreateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := reviewService.CreateReview(userID.(uint), input.MovieID, input.SeriesID, input.Rating, input.Comment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": review})
}

func GetReviews(c *gin.Context) {
	var movieID, seriesID *uint

	if mid := c.Query("movie_id"); mid != "" {
		id, err := strconv.ParseUint(mid, 10, 32)
		if err == nil {
			movieID = new(uint)
			*movieID = uint(id)
		}
	}

	if sid := c.Query("series_id"); sid != "" {
		id, err := strconv.ParseUint(sid, 10, 32)
		if err == nil {
			seriesID = new(uint)
			*seriesID = uint(id)
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	reviews, total, err := reviewService.GetReviews(movieID, seriesID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reviews,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func UpdateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var input UpdateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := reviewService.UpdateReview(uint(reviewID), userID.(uint), input.Rating, input.Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

func DeleteReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	role, _ := c.Get("role")
	isAdmin := role == "admin"

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := reviewService.DeleteReview(uint(reviewID), userID.(uint), isAdmin); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

func LikeReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	if err := reviewService.LikeReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review liked successfully"})
}

func GetAverageRating(c *gin.Context) {
	var movieID, seriesID *uint

	if mid := c.Query("movie_id"); mid != "" {
		id, err := strconv.ParseUint(mid, 10, 32)
		if err == nil {
			movieID = new(uint)
			*movieID = uint(id)
		}
	}

	if sid := c.Query("series_id"); sid != "" {
		id, err := strconv.ParseUint(sid, 10, 32)
		if err == nil {
			seriesID = new(uint)
			*seriesID = uint(id)
		}
	}

	avg, count, err := reviewService.GetAverageRating(movieID, seriesID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avg,
		"total_reviews":  count,
	})
}
