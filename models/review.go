package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	MovieID   *uint          `gorm:"index" json:"movie_id,omitempty"`
	SeriesID  *uint          `gorm:"index" json:"series_id,omitempty"`
	Rating    int            `gorm:"not null" json:"rating"`
	Comment   string         `json:"comment"`
	Likes     int            `gorm:"default:0" json:"likes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Связи
	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Movie  Movie  `gorm:"foreignKey:MovieID" json:"movie,omitempty"`
	Series Series `gorm:"foreignKey:SeriesID" json:"series,omitempty"`
}

func (Review) TableName() string {
	return "reviews"
}
