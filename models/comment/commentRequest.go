package comment

import (
	"github.com/go-playground/validator/v10"
)

type CommentRequest struct {
	PostID  int    `json:"postId" binding:"required" validate:"required,gt=0"`
	Comment string `json:"comment" binding:"required,min=2,max=500" validate:"required,min=2,max=500"`
}

var Validate = validator.New()
