package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models/friend"
	"social-media-app/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateFriend(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// Bind request body to FriendRequest struct
	request, _ := ctx.Get("request")
	friendRequest, ok := request.(friend.FriendRequest)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Use middleware to achive this, but now use hard code
	// userData := ctx.MustGet("userData").(jwt5.MapClaims)
	// userId := int(userData["id"].(float64))

	userId := 1

	err := services.CreateFriendService(DB, friendRequest, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to create friend %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully add friend",
		// "data":    friendData,
	})
}

func GetListOfFriend(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// userData := ctx.MustGet("userData").(jwt5.MapClaims)
	// userId := int(userData["id"].(float64))

	userId := 1

	pageNum, err := strconv.Atoi(ctx.Param("pageNum"))
	if err != nil || pageNum <= 0 {
		pageNum = 1
	}
	pageSize := 10
	offset := (pageNum - 1) * pageSize

	listFriendData, err := services.GetListOfFriendService(DB, userId, pageSize, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch list of friends", "message": err.Error()})
		return
	}

	var totalCount int
	err = DB.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	formattedData := formatListOfFriends(listFriendData)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully fetched list of friends",
		"data":    formattedData,
		"meta": gin.H{
			"limit":  pageSize,
			"offset": offset,
			"total":  totalCount,
		},
	})
}

func formatListOfFriends(friendData []map[string]interface{}) []map[string]interface{} {
	formattedData := make([]map[string]interface{}, len(friendData))

	for i, friend := range friendData {
		formattedData[i] = map[string]interface{}{
			"userId":      friend["user_id"],
			"name":        friend["name"],
			"imageUrl":    friend["image_url"],
			"friendCount": friend["friend_count"],
			"createdAt":   friend["created_at"].(time.Time).Format(time.RFC3339),
		}
	}

	return formattedData
}

func RemoveFriend(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// Bind request body to FriendRequest struct
	request, _ := ctx.Get("request")
	friendRequest, ok := request.(friend.FriendRequest)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// userData := ctx.MustGet("userData").(jwt5.MapClaims)
	// userId := int(userData["id"].(float64))
	userId := 1 // Replace this with the actual logic to get the user ID

	err := services.RemoveFriendService(DB, userId, friendRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to remove friend %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully delete friend",
	})
}
