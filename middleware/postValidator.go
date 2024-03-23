package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/post"

	"github.com/gin-gonic/gin"
)

func PostValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var postRequest post.PostRequest

		if payloadValidationError := context.ShouldBindJSON(&postRequest); payloadValidationError != nil {
			var errors []string

			if payloadValidationError.Error() == "EOF" {
				errors = append(errors, "Request body is missing")
			} else {
				errors = helpers.GeneralValidator(payloadValidationError)
			}

			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   errors,
				"message": "Failed to validate",
			})
			return
		}

		context.Set("request", postRequest)
		context.Next()
	}
}
