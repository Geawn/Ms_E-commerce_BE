package models

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title         string `gorm:"not null"`
	Slug          string `gorm:"uniqueIndex;not null"`
	Content       string `gorm:"type:text"`
	SeoTitle      string
	SeoDescription string
} 