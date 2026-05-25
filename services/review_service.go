package services

import (
	"errors"
	"media-content-api/config"
	"media-content-api/models"

	"gorm.io/gorm"
)

type ReviewService struct{}

func NewReviewService() *ReviewService {
	return &ReviewService{}
}

func (s *ReviewService) CreateReview(userID uint, movieID, seriesID *uint, rating int, comment string) (*models.Review, error) {
	if rating < 1 || rating > 10 {
		return nil, errors.New("rating must be between 1 and 10")
	}

	if movieID == nil && seriesID == nil {
		return nil, errors.New("either movie_id or series_id is required")
	}

	review := models.Review{
		UserID:   userID,
		MovieID:  movieID,
		SeriesID: seriesID,
		Rating:   rating,
		Comment:  comment,
	}

	result := config.DB.Create(&review)
	if result.Error != nil {
		return nil, result.Error
	}

	config.DB.Preload("User").First(&review, review.ID)
	return &review, nil
}

func (s *ReviewService) GetReviews(movieID, seriesID *uint, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := config.DB.Model(&models.Review{}).Preload("User")

	if movieID != nil {
		query = query.Where("movie_id = ?", movieID)
	}
	if seriesID != nil {
		query = query.Where("series_id = ?", seriesID)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error

	return reviews, total, err
}

func (s *ReviewService) UpdateReview(reviewID, userID uint, rating int, comment string) error {
	var review models.Review
	if err := config.DB.Where("id = ? AND user_id = ?", reviewID, userID).First(&review).Error; err != nil {
		return errors.New("review not found")
	}

	if rating < 1 || rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

	return config.DB.Model(&review).Updates(map[string]interface{}{
		"rating":  rating,
		"comment": comment,
	}).Error
}

func (s *ReviewService) DeleteReview(reviewID, userID uint, isAdmin bool) error {
	query := config.DB.Where("id = ?", reviewID)
	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	result := query.Delete(&models.Review{})
	if result.RowsAffected == 0 {
		return errors.New("review not found")
	}
	return result.Error
}

func (s *ReviewService) LikeReview(reviewID uint) error {
	return config.DB.Model(&models.Review{}).Where("id = ?", reviewID).Update("likes", gorm.Expr("likes + ?", 1)).Error
}

func (s *ReviewService) GetAverageRating(movieID, seriesID *uint) (float64, int64, error) {
	var avg float64
	var count int64

	query := config.DB.Model(&models.Review{})

	if movieID != nil {
		query = query.Where("movie_id = ?", movieID)
	}
	if seriesID != nil {
		query = query.Where("series_id = ?", seriesID)
	}

	query.Count(&count)
	query.Select("COALESCE(AVG(rating), 0)").Row().Scan(&avg)

	return avg, count, nil
}
