package services

import (
	"context"
	"errors"
	// "log"
	"social-media-app/models/comment"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFriends = errors.New("commentator and post owner are not friends")
var ErrPostNotFound = errors.New("post not found")

func CreateCommentService(DB *pgxpool.Pool, Request comment.CommentRequest, commentatorId int) error {
	ctx := context.Background()

	tx, err := DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Retrieve the post owner's ID
	var postOwnerId int
	err = tx.QueryRow(ctx, "SELECT user_id FROM Posts WHERE id = $1", Request.PostID).Scan(&postOwnerId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrPostNotFound
		}
		return err
	}

	// Check if the commentator and post owner are friends
	var isFriend bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM Friendship WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1))", commentatorId, postOwnerId).Scan(&isFriend)
	if err != nil {
		return err
	}

	if !isFriend {
		return ErrNotFriends
	}

	// Insert comment into database
	_, err = tx.Exec(ctx, "INSERT INTO comments (commentator_id, post_id, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", commentatorId, Request.PostID, Request.Comment, time.Now(), time.Now())
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
