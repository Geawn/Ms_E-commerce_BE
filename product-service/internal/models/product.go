package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name           string    `gorm:"not null"`
	Slug           string    `gorm:"uniqueIndex;not null"`
	Description    string    `gorm:"type:text"`
	SeoTitle       string    `gorm:"type:text"`
	SeoDescription string    `gorm:"type:text"`
	CategoryID     uint      `gorm:"not null"`
	Category       *Category `gorm:"foreignKey:CategoryID"`
	Variants       []*ProductVariant
	Attributes     []*ProductAttribute
	Collections    []*Collection `gorm:"many2many:product_collections;"`
	Reviews        []*Review
	Rating         float64
	Thumbnail      *Image         `gorm:"foreignKey:ProductID"`
	Pricing        ProductPricing `gorm:"embedded"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Category struct {
	gorm.Model
	Name           string     `gorm:"not null"`
	Slug           string     `gorm:"uniqueIndex;not null"`
	Description    string     `gorm:"type:text"`
	SeoTitle       string     `gorm:"type:text"`
	SeoDescription string     `gorm:"type:text"`
	Products       []*Product `json:"products" gorm:"foreignKey:CategoryID"`
}

type Collection struct {
	gorm.Model
	Name           string     `gorm:"not null"`
	Slug           string     `gorm:"uniqueIndex;not null"`
	Description    string     `gorm:"type:text"`
	SeoTitle       string     `gorm:"type:text"`
	SeoDescription string     `gorm:"type:text"`
	Products       []*Product `gorm:"many2many:product_collections;"`
}

type ProductVariant struct {
	gorm.Model
	ProductID         uint           `gorm:"not null"`
	Product           *Product       `gorm:"foreignKey:ProductID"`
	Name              string         `gorm:"not null"`
	QuantityAvailable int            `gorm:"not null"`
	Pricing           ProductPricing `gorm:"embedded"`
}

type VariantAttribute struct {
	gorm.Model
	VariantID uint            `gorm:"not null"`
	Variant   *ProductVariant `gorm:"foreignKey:VariantID"`
	Name      string          `gorm:"not null"`
	Value     string          `gorm:"not null"`
}

type ProductAttribute struct {
	gorm.Model
	ProductID uint     `gorm:"not null"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	Name      string   `gorm:"not null"`
	Values    string   `gorm:"type:text"`
}

// GetValues returns the Values as a slice of strings
func (pa *ProductAttribute) GetValues() ([]string, error) {
	var values []string
	if err := json.Unmarshal([]byte(pa.Values), &values); err != nil {
		return nil, err
	}
	return values, nil
}

// SetValues sets the Values from a slice of strings
func (pa *ProductAttribute) SetValues(values []string) error {
	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return err
	}
	pa.Values = string(jsonBytes)
	return nil
}

type Review struct {
	gorm.Model
	ProductID uint     `gorm:"not null"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	UserID    uint     `gorm:"not null"`
	User      *User    `gorm:"foreignKey:UserID"`
	Rating    float64  `gorm:"not null"`
	Comment   string   `gorm:"type:text"`
}

type User struct {
	gorm.Model
	Name    string `gorm:"not null"`
	Email   string `gorm:"uniqueIndex;not null"`
	Reviews []*Review
}

type Image struct {
	gorm.Model
	URL       string   `json:"url" gorm:"not null"`
	Alt       string   `json:"alt"`
	Size      int      `json:"size"`
	Format    string   `json:"format"`
	ProductID uint     `gorm:"not null"`
	Product   *Product `gorm:"foreignKey:ProductID"`
}

type Price struct {
	Amount   float64 `gorm:"not null"`
	Currency string  `gorm:"not null"`
}

type PriceRange struct {
	Start Price `gorm:"embedded"`
	Stop  Price `gorm:"embedded"`
}

type ProductPricing struct {
	PriceRange PriceRange `gorm:"embedded"`
}

// ProductCollection represents the many-to-many relationship between products and collections
type ProductCollection struct {
	ProductID    uint        `gorm:"primaryKey"`
	CollectionID uint        `gorm:"primaryKey"`
	Product      *Product    `gorm:"foreignKey:ProductID"`
	Collection   *Collection `gorm:"foreignKey:CollectionID"`
}

// ProductConnection represents a paginated list of products
type ProductConnection struct {
	Edges      []*ProductEdge
	PageInfo   *PageInfo
	TotalCount int
}

// ProductEdge represents a single product in a connection
type ProductEdge struct {
	Node   *Product
	Cursor string
}

// PageInfo represents pagination information
type PageInfo struct {
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     string
	EndCursor       string
}
