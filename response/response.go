package response

import (
	"github.com/gofiber/fiber/v2"
)

// ResponseError function
func ResponseError(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{"status": status, "message": message, "data": data})
}

// ResponseSuccessOneData function
func ResponseSuccessOneData(c *fiber.Ctx, message string, data interface{}) error {
	
	return c.Status(200).JSON(fiber.Map{"status": 200, "message": message, "data": data})
}

// ResponseSuccessManyData function
func ResponseSuccessManyData(c *fiber.Ctx, message string, data interface{}, page int, limit int, count int) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  200,
		"message": message,
		"data":    data,
		"page":    page,
		"limit":   limit,
		"count":   count,
	})
}