package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models"
	"social-media-app/services"
	"database/sql"

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

	// Must unique email(don't allow duplicate email)
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	row := DB.QueryRow(ginContext, query, Request.Email)
	var exists bool
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if email exists"})
		return
	}
	if err == sql.ErrNoRows {
		exists = false
	}
	if exists {
		ginContext.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

	email, err := services.LinkEmail(DB, Request, userID, ginContext)
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

	// Must unique phone(don't allow duplicate phone)
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)"
	row := DB.QueryRow(ginContext, query, Request.Phone)
	var exists bool
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if phone exists"})
		return
	}
	if err == sql.ErrNoRows {
		exists = false
	}
	if exists {
		ginContext.JSON(http.StatusConflict, gin.H{"error": "Phone already exists"})
		return
	}

	userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

	phone, err := services.LinkPhone(DB, Request, userID, ginContext)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to link phone number %s", err)})
		return
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully link their email to phone number", "data": phone})
}
