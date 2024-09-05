package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siyaga/go_rest_api/database"
	"github.com/siyaga/go_rest_api/model"
	"github.com/siyaga/go_rest_api/response"
)

// Login endpoint handler
func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ResponseError(c, 400, "Invalid request body", err)
	}

	var user model.User
	if err := database.DB.Db.Where("username = ? AND deleted_at IS NULL", req.Username).First(&user).Error; err != nil {
		return response.ResponseError(c, 401, "Invalid credentials", err)
	}

	if err := model.VerifyPassword(user.Password, req.Password); err != nil {
		return response.ResponseError(c, 401, "Invalid credentials", err)
	}

	token, err := model.GenerateJWT(user.Username)
	if err != nil {
		return response.ResponseError(c, 500, "Failed to generate token", err)
	}

	return response.ResponseSuccessOneData(c, "Login successful", model.LoginResponse{Id: user.ID,Token: token,Username: user.Username, Email: user.Email})
}

// Protected endpoint handler
func GetProtectedData(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		
		return response.ResponseError(c, 401, "Unauthorized", nil)
	}

	claims, err := model.ValidateJWT(tokenString)
	if err != nil {
		return response.ResponseError(c, 401, "Unauthorized", err)
	}

	// Access user information from claims
	username := claims.Username

	// ... your logic to fetch protected data ...

	return response.ResponseSuccessOneData(c, "Protected data", username)
}