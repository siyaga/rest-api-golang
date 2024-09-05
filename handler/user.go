package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/siyaga/go_rest_api/database"
	"github.com/siyaga/go_rest_api/model"
	"github.com/siyaga/go_rest_api/response"
)

// CreateUser function
func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)
	// Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	
	if err != nil {
		return response.ResponseError(c, 500, "Something's wrong with your input", err)
	}

	// Validate if username already exists
	var existingUserByUsername model.User
	db.Where("username = ? AND deleted_at IS NULL", user.Username).First(&existingUserByUsername)
	if existingUserByUsername.ID != uuid.Nil {
		return response.ResponseError(c, 400, "Username already exists", nil)
	}

	// Validate if email already exists
	var existingUserByEmail model.User
	db.Where("email = ? AND deleted_at IS NULL", user.Email).First(&existingUserByEmail)
	if existingUserByEmail.ID != uuid.Nil {
		return response.ResponseError(c, 400, "Email already exists", nil)
	}

	err = db.Create(&user).Error
	if err != nil {
		return response.ResponseError(c, 500, "Could not create user", err)
	}
	// Return the created user
	return response.ResponseSuccessOneData(c, "User has created", user)
}


// Get All Users from db
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	var users []model.User
	var total int64

	// Get query parameters
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	searchQuery := c.Query("search")

	// Default values if page and limit are not provided
	page := 1
	limit := 10

	// Parse page and limit if provided
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return response.ResponseError(c, 400, "Invalid page parameter", err)
		}
	}
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return response.ResponseError(c, 400, "Invalid limit parameter", err)
		}
	}


	// Calculate offset
	offset := (page - 1) * limit

	// Find users with pagination and search, excluding the password
	query := db.Model(&users).Count(&total).Offset(offset).Limit(limit)
	if searchQuery != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}
	query.Find(&users)

	// Return users with pagination information
	return response.ResponseSuccessManyData(c, "Users Found", users, page, limit, int(total))
}

// // Get All Users from db
// func GetAllUsers(c *fiber.Ctx) error {
// 	db := database.DB.Db
// 	var users []model.User
// 	var total int64

// 	// Get query parameters
// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	limit, _ := strconv.Atoi(c.Query("limit", "10"))

// 	// Calculate offset
// 	offset := (page - 1) * limit

// 	// Find users with pagination
// 	db.Model(&users).Count(&total).Offset(offset).Limit(limit).Find(&users)

// 	// Return users with pagination information
// 	return response.ResponseSuccessManyData(c, "Users Found", users, page, limit, int(total))
// }

 // GetSingleUser from db
func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db
	// get id params
	id := c.Params("id")
	var user model.User
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return response.ResponseError(c, 404, "User not found", nil)
	}
	return response.ResponseSuccessOneData(c, "User Found", user)
}

//  // update a user in db
// func UpdateUserUsername(c *fiber.Ctx) error {
// 	type updateUser struct {
// 		Username string `json:"username"`
// 	}
// 	db := database.DB.Db
// 	var user model.User
// 	// get id params
// 	id := c.Params("id")
// 	// find single user in the database by id
// 	db.Find(&user, "id = ?", id)
// 	if user.ID == uuid.Nil {
// 		return response.ResponseError(c, 404, "User not found", nil)
// 	}
// 	var updateUserData updateUser
// 	err := c.BodyParser(&updateUserData)
// 	if err != nil {
// 		return response.ResponseError(c, 500, "Something's wrong with your input", err)
// 	}
// 	user.Username = updateUserData.Username
// 	// Save the Changes
// 	db.Save(&user)
// 	// Return the updated user
// 	return response.ResponseSuccessOneData(c, "users Found", user)
// }

// update a user in db
func UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	db := database.DB.Db
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return response.ResponseError(c, 404, "User not found", nil)
	}
	
	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return response.ResponseError(c, 500, "Something's wrong with your input", err)
	}

	// --- Validation Logic ---

	// Check if username is being updated and if it's already taken by another user
	if updateUserData.Username != "" && updateUserData.Username != user.Username {
		var existingUserByUsername model.User
		db.Where("username = ? AND deleted_at IS NULL", updateUserData.Username).First(&existingUserByUsername)
		if existingUserByUsername.ID != uuid.Nil {
			return response.ResponseError(c, 400, "Username already exists", nil)
		}
		user.Username = updateUserData.Username
	}

	// Check if email is being updated and if it's already taken by another user
	if updateUserData.Email != "" && updateUserData.Email != user.Email {
		var existingUserByEmail model.User
		db.Where("email = ? AND deleted_at IS NULL", updateUserData.Email).First(&existingUserByEmail)
		if existingUserByEmail.ID != uuid.Nil {
			return response.ResponseError(c, 400, "Email already exists", nil)
		}
		user.Email = updateUserData.Email
	}

	if updateUserData.Password != "" {
		user.Password = updateUserData.Password
	}

	// Save the Changes
	db.Save(&user)

	
	return response.ResponseSuccessOneData(c, "User Found", user)

}


 // delete user in db by ID
func DeleteUserByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var user model.User
	// get id params
	id := c.Params("id")
	// find single user in the database by id
	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return response.ResponseError(c, 404, "User not found", nil)
	}
	err := db.Delete(&user, "id = ?", id).Error
	if err != nil {
		return response.ResponseError(c, 404, "Failed to delete user", nil)
	}
	return response.ResponseSuccessOneData(c, "User deleted", nil)
}