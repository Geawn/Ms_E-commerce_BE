package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuItem struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name"`
	Level      int            `json:"level"`
	CategoryID *string        `json:"categoryId"`
	Category   *Category      `json:"category,omitempty"`
	CollectionID *string      `json:"collectionId"`
	Collection *Collection    `json:"collection,omitempty"`
	PageID     *string        `json:"pageId"`
	Page       *Page          `json:"page,omitempty"`
	URL        string         `json:"url"`
	ParentID   *string        `json:"parentId"`
	Parent     *MenuItem      `json:"parent,omitempty"`
	Children   []MenuItem     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Menu struct {
	ID        string      `json:"id" gorm:"primaryKey"`
	Slug      string      `json:"slug" gorm:"uniqueIndex"`
	Channel   string      `json:"channel"`
	Items     []MenuItem  `json:"items" gorm:"foreignKey:MenuID"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Category struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type Collection struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Slug string `json:"slug"`
	Name string `json:"name"`
} 