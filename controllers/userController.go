package controllers

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserLogin(ctx *gin.Context) {
	var user models.Users

	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		fmt.Println("20")//
		fmt.Println(ok)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	Request, ok := ctx.MustGet("request").(models.UserRequest)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Request from context"})
		return
	}

	if !(Request.CredentialType == "phone" || Request.CredentialType == "email") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "credentialType should be phone or email",
		})
		return
	}

	query := "SELECT id, name, password, email, phone, image_url, credential_type, created_at, updated_at FROM users WHERE phone = $1 OR email = $1"
	row := DB.QueryRow(ctx, query, Request.CredentialValue)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Phone, &user.ImageURL, &user.CredentialType, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "User not found"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Server error occurred"})
		return
	}

	// Compare password
	comparePass := helpers.ComparePassword([]byte(user.Password), []byte(Request.Password))
	if !comparePass {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid password"})
		return
	}

	var username string
	switch user.CredentialType {
	case "email":
		username = user.Email
	case "phone":
		username = user.Phone
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(user.ID, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "message": "Failed to generate token"})
		return
	}

	// Construct response data
	responseData := gin.H{
		"message": "User logged successfully",
		"data": gin.H{
			"email":       user.Email,
			"phone":       user.Phone,
			"name":        user.Name,
			"accessToken": token,
		},
	}

	ctx.JSON(http.StatusOK, responseData)
}
