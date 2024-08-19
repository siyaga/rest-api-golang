package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siyaga/go_rest_api/model"
	"github.com/siyaga/go_rest_api/response"
)

// JWTMiddleware function
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.ResponseError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
		}

		// Split the token from the Authorization header
		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Validate the token
		claims, err := model.ValidateJWT(tokenString)
		if err != nil {
			return response.ResponseError(c, fiber.StatusUnauthorized, "Invalid token", err)
		}

		// Set the user's username in the context
		c.Locals("username", claims.Username)

		// Continue to the next handler
		return c.Next()
	}
}
