package services

import (
	"context"
	"log"
	"social-media-app/models/comment"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCommentService(DB *pgxpool.Pool, Request comment.CommentRequest, commentatorId int) error {
	ctx := context.Background()

	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("Failed to rollback transaction:", rollbackErr)
			}
			log.Println("Transaction aborted:", err)
		}
	}()

	// var id int

	// Prepare statement within the transaction
	_, err = tx.Exec(ctx, "INSERT INTO comments (commentator_id, post_id, comment, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)", commentatorId, Request.PostID, Request.Comment, time.Now(), time.Now())
	if err != nil {
		log.Println("Error executing insert statement:", err)
		return err
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}
	return nil
}
