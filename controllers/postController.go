package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models/post"
	"social-media-app/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePost(ginContext *gin.Context) {
	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	Request, ok := ginContext.MustGet("request").(post.PostRequest)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Request from context"})
		return
	}

	// TODO: handle with auth middleware
	// userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	// userID := int(userData["id"].(float64))
	userID := 1

	post, err := services.CreatePost(DB, Request, userID)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to insert post %s", err)})
		return
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully add post", "data": post})
}

func GetPost(ginContext *gin.Context) {
	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// TODO: handle with auth middleware
	// userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	// userID := int(userData["id"].(float64))
	userID := 1
	posts, err := services.GetPostsByUserId(DB, userID)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get posts %s", err)})
		return
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully get posts", "data": posts})
}
