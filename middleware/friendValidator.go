package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/friend"

	"github.com/gin-gonic/gin"
)

func FriendValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var friendRequest friend.FriendRequest

		if payloadValidationError := ctx.ShouldBindJSON(&friendRequest); payloadValidationError != nil {
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

		ctx.Set("request", friendRequest)
		ctx.Next()
	}
}
