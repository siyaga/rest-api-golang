package model

import (
	"encoding/json"

	"github.com/google/uuid"
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