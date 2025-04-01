package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string  `gorm:"uniqueIndex;not null" json:"email"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Password  string  `json:"-"`
	Role      string  `gorm:"default:'user'" json:"role"`
	Profile   Profile `gorm:"foreignKey:UserID" json:"profile"`
}

type Profile struct {
	gorm.Model
	UserID      uint      `json:"userId"`
	PhoneNumber string    `json:"phoneNumber"`
	Addresses   []Address `gorm:"foreignKey:ProfileID" json:"addresses"`
	Avatar      *Avatar   `gorm:"foreignKey:ProfileID" json:"avatar"`
}

type Address struct {
	gorm.Model
	ProfileID  uint   `json:"profileId"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
	IsDefault  bool   `json:"isDefault"`
}

type Avatar struct {
	gorm.Model
	ProfileID uint   `json:"profileId"`
	URL       string `json:"url"`
	Alt       string `json:"alt"`
}
