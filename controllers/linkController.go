package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models"
	"social-media-app/services"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LinkEmail(ginContext *gin.Context) {
	handleDBError := func(err error, message string) bool {
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": message})
			return true
		}
		return false
	}

	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		handleDBError(nil, "Failed to get DB from context")
		return
	}

	Request, ok := ginContext.MustGet("request").(models.LinkEmailRequest)
	if !ok {
		handleDBError(nil, "Failed to get Request from context")
		return
	}

	// Must unique email(don't allow duplicate email) | 400
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 and credential_type = 'email')"
	row := DB.QueryRow(ginContext, query, Request.Email)
	var exists bool
	if err := row.Scan(&exists); handleDBError(err, "Error checking if email exists") {
		return
	}
	if exists {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Must unique email(don't allow duplicate email) | 409
	query = "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	row = DB.QueryRow(ginContext, query, Request.Email)
	err = row.Scan(&exists)
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
	if handleDBError(err, fmt.Sprintf("Failed to link email: %s", err)) {
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Successfully linked email", "data": email})
}

func LinkPhone(ginContext *gin.Context) {
	handleDBError := func(err error, message string) bool {
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": message})
			return true
		}
		return false
	}

	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		handleDBError(nil, "Failed to get DB from context")
		return
	}

	Request, ok := ginContext.MustGet("request").(models.LinkPhoneRequest)
	if !ok {
		handleDBError(nil, "Failed to get Request from context")
		return
	}

	// Must unique phone(don't allow duplicate phone) | 400
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1 and credential_type = 'phone')"
	row := DB.QueryRow(ginContext, query, Request.Phone)
	var exists bool
	if err := row.Scan(&exists); handleDBError(err, "Error checking if phone exists") {
		return
	}
	if exists {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Phone already exists"})
		return
	}

	// Must unique phone(don't allow duplicate phone) | 409
	query = "SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)"
	row = DB.QueryRow(ginContext, query, Request.Phone)
	err = row.Scan(&exists)
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
	if handleDBError(err, fmt.Sprintf("Failed to link phone number: %s", err)) {
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"message": "Successfully linked phone number", "data": phone})
}
