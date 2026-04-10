package models

import (
	"time"

	"gorm.io/gorm"
)

type Series struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	ReleaseYear int            `json:"release_year"`
	Seasons     int            `json:"seasons"`
	Episodes    int            `json:"episodes"`
	Rating      float64        `json:"rating"`
	Genre       string         `json:"genre"`
	Creator     string         `json:"creator"`
	PosterURL   string         `json:"poster_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Series) TableName() string {
	return "series"
}
