package models

import (
	"time"

	"gorm.io/gorm"
)

type Page struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	Title          string         `json:"title"`
	Slug           string         `gorm:"uniqueIndex" json:"slug"`
	Content        string         `json:"content"`
	SeoTitle       string         `json:"seo_title"`
	SeoDescription string         `json:"seo_description"`
	IsPublished    bool           `json:"is_published" gorm:"default:false"`
	PublishedAt    *time.Time     `json:"published_at"`
}
