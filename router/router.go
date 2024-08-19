package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siyaga/go_rest_api/handler"
	// "github.com/siyaga/go_rest_api/middleware"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
 // grouping
 api := app.Group("/api")
//  api.Use(middleware.JWTMiddleware())
 v1 := api.Group("/user")
 // routes
 v1.Get("/", handler.GetAllUsers)
 v1.Get("/:id", handler.GetSingleUser)
 v1.Post("/", handler.CreateUser)
 v1.Put("/:id", handler.UpdateUser)
 v1.Delete("/:id", handler.DeleteUserByID)

// Add new authority route
api.Post("/authority", handler.Login) // Assuming you want to use the Login handler for authority

}