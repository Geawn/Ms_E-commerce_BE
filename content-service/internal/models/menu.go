package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"not null"`
	Slug string `gorm:"uniqueIndex;not null"`
}

type Collection struct {
	gorm.Model
	Name string `gorm:"not null"`
	Slug string `gorm:"uniqueIndex;not null"`
}

type MenuItem struct {
	gorm.Model
	MenuID       uint
	Name         string `gorm:"not null"`
	Level        int    `gorm:"not null;default:1"`
	URL          string
	CategoryID   *uint
	Category     *Category `gorm:"foreignKey:CategoryID"`
	CollectionID *uint
	Collection   *Collection `gorm:"foreignKey:CollectionID"`
	PageID       *uint
	Page         *Page `gorm:"foreignKey:PageID"`
	ParentID     *uint
	Children     []MenuItem `gorm:"foreignKey:ParentID"`
}

type Menu struct {
	gorm.Model
	Name    string     `gorm:"not null"`
	Slug    string     `gorm:"uniqueIndex;not null"`
	Channel string     `gorm:"not null"`
	Items   []MenuItem `gorm:"foreignKey:MenuID"`
}