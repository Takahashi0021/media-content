package models

import (
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	MovieID   *uint          `gorm:"index" json:"movie_id,omitempty"`
	SeriesID  *uint          `gorm:"index" json:"series_id,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// Убираем связи, чтобы не было проблем с Preload
}

func (Favorite) TableName() string {
	return "favorites"
}
