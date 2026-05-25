package services

import (
	"errors"
	"media-content-api/config"
	"media-content-api/models"
)

type FavoriteService struct{}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{}
}

func (s *FavoriteService) AddFavorite(userID uint, movieID, seriesID *uint) error {
	if movieID == nil && seriesID == nil {
		return errors.New("either movie_id or series_id is required")
	}

	favorite := models.Favorite{
		UserID:   userID,
		MovieID:  movieID,
		SeriesID: seriesID,
	}

	return config.DB.Create(&favorite).Error
}

func (s *FavoriteService) RemoveFavorite(userID uint, favoriteID uint) error {
	result := config.DB.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		return errors.New("favorite not found")
	}
	return result.Error
}

// ИСПРАВЛЕНО: убираем Preload
func (s *FavoriteService) GetUserFavorites(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	// Убираем .Preload("User") - это вызывало ошибку
	err := config.DB.Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

func (s *FavoriteService) GetMovieFavoritesCount(movieID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Favorite{}).Where("movie_id = ?", movieID).Count(&count).Error
	return count, err
}

func (s *FavoriteService) GetSeriesFavoritesCount(seriesID uint) (int64, error) {
	var count int64
	err := config.DB.Model(&models.Favorite{}).Where("series_id = ?", seriesID).Count(&count).Error
	return count, err
}

func (s *FavoriteService) IsFavorite(userID uint, movieID, seriesID *uint) (bool, error) {
	var count int64
	query := config.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)

	if movieID != nil {
		query = query.Where("movie_id = ?", movieID)
	}
	if seriesID != nil {
		query = query.Where("series_id = ?", seriesID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}
