package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	ReleaseYear int            `json:"release_year"`
	Duration    int            `json:"duration"`
	Rating      float64        `json:"rating"`
	Genre       string         `json:"genre"`
	Director    string         `json:"director"`
	PosterURL   string         `json:"poster_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Movie) TableName() string {
	return "movies"
}
