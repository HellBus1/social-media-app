package models

import (
	"github.com/go-playground/validator/v10"
	"social-media-app/helpers"
	"time"
	"database/sql"
)

type Users struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required,min=5,max=50" validate:"required,min=5,max=50"`
	Password string `json:"password" binding:"required,min=5,max=15" validate:"required,min=5,max=15"`

	Email     string    `json:"email,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	CredentialValue string `json:"credentialValue" binding:"required" validate:"required"` //TODO: not yet validation phone and email value

	ImageURL  string    `json:"image_url"`
	CredentialType string `json:"credentialType" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersForAuth struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required,min=5,max=50" validate:"required,min=5,max=50"`
	Password string `json:"password" binding:"required,min=5,max=15" validate:"required,min=5,max=15"`

	Email     sql.NullString    `json:"email"`
	Phone     sql.NullString    `json:"phone"`
	CredentialValue string `json:"credentialValue" binding:"required" validate:"required"` //TODO: not yet validation phone and email value

	ImageURL  string    `json:"image_url"`
	CredentialType string `json:"credentialType" binding:"required" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	CredentialType string `json:"credentialType" binding:"required" validate:"required"`
	CredentialValue string `json:"credentialValue" binding:"required" validate:"required"` //TODO: not yet validation phone and email value
	Password string `json:"password" binding:"required,min=5,max=15" validate:"required,min=5,max=15"`
}

// HashPassword hashes the password before creating the user
func (u *UsersForAuth) HashPassword() error {
	// Hash the password using a hashing function like bcrypt
	hashedPassword, err := helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

// BeforeCreateUser is a function to be called before creating a new user
func BeforeCreateUser(user *UsersForAuth) {
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
