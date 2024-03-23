package updateaccount

import (
    "github.com/go-playground/validator/v10"
)

type UpdateAccountRequest struct {
    ImageURL string `json:"imageUrl" binding:"required" validate:"required,url"`
    Name     string `json:"name" binding:"required" validate:"required,min=5,max=50"`
}

var Validate = validator.New()