package models

import (
	"github.com/go-playground/validator/v10"
	"social-media-app/helpers"
	"time"
)

type Users struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"min=5,max=15" validate:"min=5,max=15"`
	Password  string    `json:"password" binding:"min=5,max=15" validate:"min=5,max=15"`
	Email     string    `json:"email" binding:"min=5,max=50" validate:"min=5,max=50"`
	Phone     string    `json:"phone" binding:"min=5,max=50" validate:"min=5,max=50"`
	ImageURL  string    `json:"image_url"`
	CredentialType string `json:"credential_type" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	CredentialType string `json:"credentialType" binding:"required" validate:"required"`
	CredentialValue string `json:"credentialValue" binding:"required" validate:"required"` //TODO: not yet validation phone and email value
	Password string `json:"password" binding:"required,min=5,max=15" validate:"required,min=5,max=15"`
}

// HashPassword hashes the password before creating the user
func (u *Users) HashPassword() error {
	// Hash the password using a hashing function like bcrypt
	hashedPassword, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// BeforeCreateUser is a function to be called before creating a new user
func BeforeCreateUser(user *Users) {
	// Perform any pre-create logic here
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.HashPassword()
}

func ValidateUser(user *Users) error {
	validate := validator.New()
	return validate.Struct(user)
}

type LinkEmailRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

type LinkEmailResponse struct {
	Email string `json:"email"`
}

type LinkPhoneRequest struct {
	Phone string `json:"phone" binding:"required,min=7,max=13,e164" validate:"required,min=7,max=13,e164"`
}

type LinkPhoneResponse struct {
	Phone string `json:"phone"`
}
