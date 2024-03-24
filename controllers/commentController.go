package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"social-media-app/models/comment"
	"social-media-app/services"

	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

func CreateComment(ctx *gin.Context) {
    DB, ok := ctx.MustGet("DB").(*pgxpool.Pool)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
        return
    }

    var commentRequest comment.CommentRequest
    if err := ctx.ShouldBindJSON(&commentRequest); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    commentatorId := 5

    // Check if the post exists
    var postExists bool
    err := DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM Posts WHERE id = $1)", commentRequest.PostID).Scan(&postExists)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": "Failed to check if post exists",
        })
        return
    }

    if !postExists {
        ctx.JSON(http.StatusNotFound, gin.H{
            "message": fmt.Sprintf("Post with ID %d not found", commentRequest.PostID),
        })
        return
    }

    err = services.CreateCommentService(DB, commentRequest, commentatorId)
    if err != nil {
        if errors.Is(err, services.ErrNotFriends) {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "message": "Commentator and post owner are not friends",
            })
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "message": fmt.Sprintf("Failed to create comment: %s", err),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Successfully added comment",
    })
}