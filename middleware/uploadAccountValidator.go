package middleware

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models/updateAccount"

	"github.com/gin-gonic/gin"
)

func UpdateAccountValidator() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var updateAccountRequest updateaccount.UpdateAccountRequest
        if payloadValidationError := ctx.ShouldBindJSON(&updateAccountRequest); payloadValidationError != nil {
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

        err := updateaccount.Validate.Struct(updateAccountRequest)
        if err != nil {
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
            return
        }

        ctx.Set("request", updateAccountRequest)
        ctx.Next()
    }
}
