package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models/updateAccount"
	"social-media-app/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateAccountController(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// Bind request body to CommentRequest struct
	request, _ := ctx.Get("request")
	updateAccountRequest, ok := request.(updateaccount.UpdateAccountRequest)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// userData := ctx.MustGet("userData").(jwt5.MapClaims)
	// userId := int(userData["id"].(float64))
	userId := 1

	err := services.UpdateAccountService(DB, updateAccountRequest, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to create friend %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully update user profile",
		// "data":    updateAccountData,
	})
}
