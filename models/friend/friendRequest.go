package friend

import (
	"github.com/go-playground/validator/v10"
)

type FriendRequest struct {
	UserID int `json:"userId" binding:"required" validate:"required,gt=0"`
}

var Validate = validator.New()

