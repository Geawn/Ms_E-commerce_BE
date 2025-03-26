package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `gorm:"not null"`
	Slug        string   `gorm:"uniqueIndex;not null"`
	Description string   `gorm:"type:text"`
	Price       float64  `gorm:"not null"`
	Stock       int      `gorm:"not null"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Variants    []ProductVariant
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	Slug        string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Products    []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

type ProductVariant struct {
	gorm.Model
	ProductID  uint               `gorm:"not null"`
	Product    Product            `gorm:"foreignKey:ProductID"`
	Name       string             `gorm:"not null"`
	Price      float64            `gorm:"not null"`
	Stock      int                `gorm:"not null"`
	Attributes []VariantAttribute `json:"attributes" gorm:"foreignKey:VariantID"`
}

type VariantAttribute struct {
	gorm.Model
	VariantID uint           `gorm:"not null"`
	Variant   ProductVariant `gorm:"foreignKey:VariantID"`
	Name      string         `gorm:"not null"`
	Value     string         `gorm:"not null"`
}

type Price struct {
	Amount   float64
	Currency string
}

type PriceRange struct {
	Start Price
	Stop  Price
}

type ProductPricing struct {
	PriceRange PriceRange
}
