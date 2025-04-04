package model

import (
	"fmt"
	"io"
	"time"
)

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Avatar    *Image
}

// Order represents an order in the system
type Order struct {
	ID            string
	UserID        string
	Number        string
	Created       time.Time
	TotalAmount   float64
	Currency      string
	PaymentStatus PaymentStatusEnum
	Lines         []*OrderLine
}

// OrderLine represents a line item in an order
type OrderLine struct {
	ID          string
	OrderID     string
	VariantID   string
	VariantName string
	ProductSlug string
	Price       float64
	Currency    string
	Quantity    int
}

// ProductVariant represents a variant of a product
type ProductVariant struct {
	ID      string
	Name    string
	Product *Product
	Pricing *Pricing
}

// Product represents a product in the system
type Product struct {
	ID          string
	Name        string
	Description string
	Slug        string
	Thumbnail   *Image
	Category    *Category
	Variants    []*ProductVariant
}

// Category represents a product category
type Category struct {
	ID   string
	Name string
}

// Image represents an image with URL and alt text
type Image struct {
	URL string
	Alt string
}

// Pricing represents the pricing information
type Pricing struct {
	Price *Price
}

// Price represents a price with gross amount
type Price struct {
	Gross *Money
}

// Money represents an amount of money with currency
type Money struct {
	Amount   string
	Currency string
}

// PaymentStatusEnum represents the possible payment statuses
type PaymentStatusEnum string

const (
	PaymentStatusPending PaymentStatusEnum = "PENDING"
	PaymentStatusPaid    PaymentStatusEnum = "PAID"
	PaymentStatusFailed  PaymentStatusEnum = "FAILED"
)

// OrderConnection represents a connection of orders
type OrderConnection struct {
	Edges []*OrderEdge
}

// OrderEdge represents an edge in the order connection
type OrderEdge struct {
	Node *Order
}

// DateTime represents a date and time
type DateTime time.Time

// MarshalGQL implements the graphql.Marshaler interface
func (t DateTime) MarshalGQL(w io.Writer) {
	w.Write([]byte(time.Time(t).Format(time.RFC3339)))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *DateTime) UnmarshalGQL(v interface{}) error {
	if str, ok := v.(string); ok {
		parsed, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return err
		}
		*t = DateTime(parsed)
		return nil
	}
	return fmt.Errorf("DateTime must be a string")
}
