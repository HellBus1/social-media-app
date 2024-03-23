package controllers

import (
	"fmt"
	"net/http"
	"social-media-app/models/comment"
	"social-media-app/services"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

func CreateComment(ctx *gin.Context) {
	DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}

	// Bind request body to CommentRequest struct
	request, _ := ctx.Get("request")
    commentRequest, ok := request.(comment.CommentRequest)
    if !ok {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

	// Use middleware to achive this, but now use hard code
	// commentatorData := ctx.MustGet("commentatorData").(jwt5.MapClaims)
	// commentatorId := int(commentatorData["id"].(float64))

	commentatorId := 1

	err := services.CreateCommentService(DB, commentRequest, commentatorId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to create comment %s", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully add post",
		// "data":    commentData,
	})
}
