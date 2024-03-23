package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models"

	"github.com/gin-gonic/gin"
)

func LinkEmailValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var linkEmailRequest models.LinkEmailRequest

		if payloadValidationError := context.ShouldBindJSON(&linkEmailRequest); payloadValidationError != nil {
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

		context.Set("request", linkEmailRequest)
		context.Next()
	}
}

func LinkPhoneValidator() gin.HandlerFunc {
	return func(context *gin.Context) {
		var linkPhoneRequest models.LinkPhoneRequest

		if payloadValidationError := context.ShouldBindJSON(&linkPhoneRequest); payloadValidationError != nil {
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

		context.Set("request", linkPhoneRequest)
		context.Next()
	}
}
