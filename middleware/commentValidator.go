package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/comment"

	"github.com/gin-gonic/gin"
)

func CommentValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var commentRequest comment.CommentRequest

		if payloadValidationError := ctx.ShouldBindJSON(&commentRequest); payloadValidationError != nil {
			var errors []string

			if payloadValidationError.Error() == "EOF" {
				errors = append(errors, "Request body is missing")
			} else {
				errors = helpers.GeneralValidator(payloadValidationError)
			}

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   errors,
				"message": "Failed to validate",
			})
			return
		}

		ctx.Set("request", commentRequest)
		ctx.Next()
	}
}