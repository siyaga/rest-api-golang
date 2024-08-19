package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/siyaga/go_rest_api/config"
	"golang.org/x/crypto/bcrypt" // Import bcrypt package
	"gorm.io/gorm"
)

// User struct
type User struct {
 gorm.Model
 ID       uuid.UUID `gorm:"type:uuid;"`
 Username string    `json:"username"`
 Email    string    `json:"email"`
 Password string    `json:"password"`
}
// Users struct
type Users struct {
 Users []User `json:"users"`
}
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
 // UUID version 4
 user.ID = uuid.New()
 hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword) // Store the hashed password
 return
}
// Custom JSON Marshaling for Datetime Fields
func (u *User) MarshalJSON() ([]byte, error) {
	type UserAlias User
	return json.Marshal(&struct {
		UserAlias
		CreatedAt string `json:"CreatedAt"`
		UpdatedAt string `json:"UpdatedAt"`
		DeletedAt *string `json:"DeletedAt"`
	}{
		UserAlias: UserAlias(*u),
		CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: func() *string {
			if !u.DeletedAt.Time.IsZero() { // Access the Time field
				formatted := u.DeletedAt.Time.Format("2006-01-02 15:04:05")
				return &formatted // Return a pointer to the formatted string
			}
			return nil // Return nil for null
		}(),
	})
}

// LoginRequest struct for login API
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse struct for login API response
type LoginResponse struct {
	Id uuid.UUID `json:"id"`
	Token string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// VerifyPassword function to compare plain password with hashed password
func VerifyPassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
// JWT Claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT function to generate JWT token
func GenerateJWT(username string) (string, error) {
	// Create the Claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			Issuer:    "go_rest_api",                        // Issuer of the token
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	secretKey := []byte(config.Config("JWT_SECRET")) // Replace with your actual secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT function to validate JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secretKey := []byte(config.Config("JWT_SECRET")) // Replace with your actual secret key
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	// Validate the token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}