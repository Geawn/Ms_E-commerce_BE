package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/database"
	"github.com/Geawn/Ms_E-commerce_BE/account-service/models"
	"github.com/Geawn/Ms_E-commerce_BE/account-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var userServiceClient *utils.UserServiceClient

func init() {
	var err error
	userServiceClient, err = utils.NewUserServiceClient()
	if err != nil {
		panic("Failed to create user service client: " + err.Error())
	}
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	var existingUser models.AccountUser
	result := database.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
		return
	}

	// Create user in user-service via gRPC
	userID, err := userServiceClient.CreateUser(c.Request.Context(), req.Email, req.FirstName, req.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user in user service"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create new account user
	user := models.AccountUser{
		UserID:    strconv.FormatUint(userID, 10),
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  string(hashedPassword),
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user first
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Commit user creation first
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit user creation"})
		return
	}

	// Start new transaction for token creation
	tx = database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start token transaction"})
		return
	}

	// Verify user was created
	var createdUser models.AccountUser
	if err := tx.Where("user_id = ?", user.UserID).First(&createdUser).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying user creation"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(createdUser, "web", c.ClientIP())
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	var userToken models.UserToken
	// Check if token already exists
	var existingToken models.UserToken
	if err := tx.Where("token = ?", token).First(&existingToken).Error; err == nil {
		// Token exists, use it
		userToken = existingToken
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing token"})
		return
	} else {
		// Token doesn't exist, create new one
		userToken = models.UserToken{
			UserID:     createdUser.UserID,
			Token:      token,
			TokenType:  "access",
			ExpiresAt:  time.Now().Add(24 * time.Hour),
			DeviceInfo: "web",
			IPAddress:  c.ClientIP(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := tx.Create(&userToken).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user token"})
			return
		}
	}

	// Commit token creation
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit token creation"})
		return
	}

	response := models.AuthResponse{
		Token: token,
		User: models.AuthUserResponse{
			ID:        createdUser.ID,
			Username:  createdUser.Username,
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			Role:      createdUser.Role,
		},
	}

	c.JSON(http.StatusCreated, response)
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.AccountUser
	result := database.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user, "access", "24h")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Clear password before sending response
	user.Password = ""
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User: models.AuthUserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
		},
	})
}

func ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*utils.Claims)

	// Get user from database
	var user models.AccountUser
	if err := database.DB.Where("user_id = ?", userClaims.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Update password
	if err := database.DB.Model(&user).Update("password", string(newHash)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func Logout(c *gin.Context) {
	// Since we're using JWT, we don't need to do anything server-side
	// The client should simply remove the token from their storage
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
