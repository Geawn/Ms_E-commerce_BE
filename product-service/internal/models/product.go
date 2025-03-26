package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Slug        string  `json:"slug" gorm:"uniqueIndex;not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"not null"`
	CategoryID  uint    `json:"category_id" gorm:"not null"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
	Variants    []ProductVariant `json:"variants" gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	gorm.Model
	Name        string    `json:"name" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	Products    []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

type ProductVariant struct {
	gorm.Model
	ProductID   uint      `json:"product_id" gorm:"not null"`
	Product     Product   `json:"product" gorm:"foreignKey:ProductID"`
	Name        string    `json:"name" gorm:"not null"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"not null"`
	Attributes  []VariantAttribute `json:"attributes" gorm:"foreignKey:VariantID"`
}

type VariantAttribute struct {
	gorm.Model
	VariantID   uint   `json:"variant_id" gorm:"not null"`
	Variant     ProductVariant `json:"variant" gorm:"foreignKey:VariantID"`
	Name        string `json:"name" gorm:"not null"`
	Value       string `json:"value" gorm:"not null"`
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
