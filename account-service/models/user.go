package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Password  string         `json:"password,omitempty" gorm:"not null"`
	Role      string         `json:"role" gorm:"not null;default:user"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	Tokens    []UserToken    `json:"-" gorm:"foreignKey:UserID;references:UserID"`
}

// TableName specifies the table name for AccountUser model
func (AccountUser) TableName() string {
	return "account_users"
}

type UserToken struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     string         `json:"user_id" gorm:"not null;index"`
	Token      string         `json:"token" gorm:"not null;uniqueIndex"`
	TokenType  string         `json:"token_type" gorm:"not null;default:'bearer'"`
	ExpiresAt  time.Time      `json:"expires_at" gorm:"not null"`
	LastUsedAt *time.Time     `json:"last_used_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	DeviceInfo string         `json:"device_info"`
	IPAddress  string         `json:"ip_address"`
	User       AccountUser    `json:"-" gorm:"foreignKey:UserID;references:UserID"`
}

// TableName specifies the table name for UserToken model
func (UserToken) TableName() string {
	return "user_tokens"
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AuthUserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

type AuthResponse struct {
	Token string           `json:"token"`
	User  AuthUserResponse `json:"user"`
}
