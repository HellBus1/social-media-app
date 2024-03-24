package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models/post"
	"social-media-app/services"
	"strconv"

	"github.com/gin-gonic/gin"
	jwt5 "github.com/golang-jwt/jwt/v5"
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

	userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))

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

	limit, _ := strconv.Atoi(ginContext.Query("limit"))
	offset, _ := strconv.Atoi(ginContext.Query("offset"))
	search := ginContext.Query("search")
	searchTags := ginContext.QueryArray("searchTag")
	
	if limit <= 0 {
			limit = 5
	}
	
	if offset < 0 {
			offset = 0
	}
	
	calculatedOffset := offset * limit
	
	if calculatedOffset < 0 {
			calculatedOffset = 0
	}

	userData := ginContext.MustGet("userData").(jwt5.MapClaims)
	userID := int(userData["id"].(float64))
	
	posts, err := services.GetPostsByUserId(DB, userID, search, searchTags, limit, calculatedOffset)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get posts %s", err)})
		return
	}

	meta := post.Meta{
		Limit: limit,
		Offset: offset,
		Total: len(*posts),
	}

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully get posts", "data": posts, "meta": meta})
}
