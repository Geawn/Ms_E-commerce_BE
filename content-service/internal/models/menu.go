package models

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Slug      string         `gorm:"uniqueIndex" json:"slug"`
	Channel   string         `json:"channel"`
	Items     []MenuItem     `json:"items" gorm:"foreignKey:MenuID"`
}

type MenuItem struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	MenuID       uint           `json:"menu_id"`
	Name         string         `json:"name"`
	Level        int            `json:"level"`
	URL          string         `json:"url"`
	CategoryID   *uint          `json:"category_id"`
	Category     *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	CollectionID *uint          `json:"collection_id"`
	Collection   *Collection    `json:"collection,omitempty" gorm:"foreignKey:CollectionID"`
	PageID       *uint          `json:"page_id"`
	Page         *Page          `json:"page,omitempty" gorm:"foreignKey:PageID"`
	ParentID     *uint          `json:"parent_id"`
	Children     []MenuItem     `json:"children" gorm:"foreignKey:ParentID"`
}

type Category struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Description string         `json:"description"`
}

type Collection struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Slug        string         `gorm:"uniqueIndex" json:"slug"`
	Description string         `json:"description"`
}
