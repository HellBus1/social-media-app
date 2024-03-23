package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models"
	"social-media-app/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LinkEmail(ginContext *gin.Context) {
	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	Request, ok := ginContext.MustGet("request").(models.LinkEmailRequest)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Request from context"})
		return
	}

	// TODO: handle with auth middleware
	// userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	// userID := int(userData["id"].(float64))
	userID := 1

	email, err := services.LinkEmail(DB, Request, userID)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to link email %s", err)})
		return
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": " successfully link their phone number to email", "data": email})
}

func LinkPhone(ginContext *gin.Context) {
	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	Request, ok := ginContext.MustGet("request").(models.LinkPhoneRequest)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Request from context"})
		return
	}

	// TODO: handle with auth middleware
	// userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	// userID := int(userData["id"].(float64))
	userID := 1

	phone, err := services.LinkPhone(DB, Request, userID)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to link phone number %s", err)})
		return
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully link their email to phone number", "data": phone})
}
