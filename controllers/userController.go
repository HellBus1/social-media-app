package controllers

import (
	"net/http"
	"social-media-app/helpers"
	"social-media-app/models"
	"fmt"
	"database/sql"

	"github.com/gin-gonic/gin"
	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRegister(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}
	fmt.Println(DB)//

	user, ok := ctx.MustGet("request").(models.UsersForAuth)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Request from context"})
		return
	}

	if !(user.CredentialType == "phone" || user.CredentialType == "email") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "credentialType should be phone or email",
		})
		return
	}

	var query string
	if (user.CredentialType == "phone") {
		query = "SELECT EXISTS (SELECT 1 FROM users WHERE phone = $1)"
	} else {
		query = "SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)"
	}

	row := DB.QueryRow(ctx, query, user.CredentialValue)
	var exists bool
	err := row.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if phone or email exists"})
		return
	}
	if err == sql.ErrNoRows {
		exists = false
	}
	if exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": "phone or email already exists"})
		return
	}

	// Before creating the user, perform necessary operations
	models.BeforeCreateUser(&user)

	// Save the user to the database
	var query_register string
	if (user.CredentialType == "phone") {
		query_register = "INSERT INTO users (name, password, phone, credential_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	} else {
		query_register = "INSERT INTO users (name, password, email, credential_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"
	}
	_, err = DB.Exec(ctx, query_register, user.Name, user.Password, user.CredentialValue, user.CredentialType, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		fmt.Println("109")//
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	var email string
	var phone string
	if user.Email.Valid {
		email = user.Email.String
	}
	if user.Phone.Valid {
		phone = user.Phone.String
	}

	// Generate JWT token
	var username string
	switch user.CredentialType {
	case "email":
		username = email
	case "phone":
		username = phone
	}
	token, err := helpers.GenerateToken(user.ID, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
			"message": "Failed to generate token",
		})
		return
	}
	
	// Construct response data
	var responseData gin.H
	if (user.CredentialType == "phone") {
		responseData = gin.H{
			"message": "User registered successfully",
			"data": gin.H{
				"phone": user.CredentialValue,
				"name": user.Name,
				"accessToken": token,
			},
		}
	} else {
		responseData = gin.H{
			"message": "User registered successfully",
			"data": gin.H{
				"email": user.CredentialValue,
				"name": user.Name,
				"accessToken": token,
			},
		}
	}

	ctx.JSON(http.StatusCreated, responseData)
}

func UserLogin(ctx *gin.Context) {
	// var user models.Users
	var user models.UsersForAuth

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

	var query string
	if (Request.CredentialType == "phone") {
		query = "SELECT id, name, password, email, phone, image_url, credential_type, created_at, updated_at FROM users WHERE phone = $1"
	} else {
		query = "SELECT id, name, password, email, phone, image_url, credential_type, created_at, updated_at FROM users WHERE email = $1"
	}
	row := DB.QueryRow(ctx, query, Request.CredentialValue)
	err := row.Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.Phone, &user.ImageURL, &user.CredentialType, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		fmt.Println("272")//
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not Found", "message": "User not found"})
		return
	}

	// Compare password
	comparePass := helpers.ComparePassword([]byte(user.Password), []byte(Request.Password))
	if !comparePass {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Invalid password"})
		return
	}

	var email string
	var phone string
	if user.Email.Valid {
		email = user.Email.String
	}
	if user.Phone.Valid {
		phone = user.Phone.String
	}

	var username string
	switch user.CredentialType {
	case "email":
		username = email
	case "phone":
		username = phone
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
			"email": email,
			"phone": phone,
			"name": user.Name,
			"accessToken": token,
		},
	}
	// If the email or phone value is null, change it to an empty string
	if !user.Email.Valid {
		responseData["data"].(gin.H)["email"] = ""
	}
	if !user.Phone.Valid {
		responseData["data"].(gin.H)["phone"] = ""
	}

	ctx.JSON(http.StatusOK, responseData)
}
