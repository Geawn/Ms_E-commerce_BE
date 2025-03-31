package models

import (
	"time"

	"gorm.io/gorm"
)

type Page struct {
	ID             string         `json:"id" gorm:"primaryKey"`
	Slug           string         `json:"slug" gorm:"uniqueIndex"`
	Title          string         `json:"title"`
	SeoTitle       string         `json:"seoTitle"`
	SeoDescription string         `json:"seoDescription"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
} 