package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/database"
	"github.com/Geawn/Ms_E-commerce_BE/account-service/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	// Load .env file if it hasn't been loaded
	if os.Getenv("JWT_SECRET") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	jwtKey = []byte(secret)
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(user models.AccountUser, deviceInfo, ipAddress string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	// Save token to database
	userToken := models.UserToken{
		UserID:     user.ID,
		Token:      tokenString,
		TokenType:  "access",
		ExpiresAt:  expirationTime,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	}

	if err := database.DB.Create(&userToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	// First check if token exists in database
	var userToken models.UserToken
	if err := database.DB.Where("token = ? AND expires_at > ?", tokenStr, time.Now()).First(&userToken).Error; err != nil {
		return nil, errors.New("token not found or expired")
	}

	// Update last used time
	now := time.Now()
	database.DB.Model(&userToken).Update("last_used_at", now)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func InvalidateToken(tokenStr string) error {
	return database.DB.Where("token = ?", tokenStr).Delete(&models.UserToken{}).Error
}

func InvalidateAllUserTokens(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.UserToken{}).Error
}

func CleanupExpiredTokens() error {
	return database.DB.Where("expires_at < ?", time.Now()).Delete(&models.UserToken{}).Error
}
