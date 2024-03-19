package services

import (
	"context"
	"log"
	"social-media-app/models/post"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreatePost inserts a new post with associated tags into the database
func CreatePost(DB *pgxpool.Pool, postReq post.PostRequest, userID int) (*post.PostResponse, error) {
	ctx := context.Background()

	tx, err := DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return nil, err
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

	// Insert post
	var postID int
	err = tx.QueryRow(ctx, "INSERT INTO posts (user_id, post_in_html, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id",
		userID, postReq.FeedInHtml, time.Now(), time.Now()).Scan(&postID)
	if err != nil {
		log.Println("Failed to insert into posts")
		return nil, err
	}

	// Insert tags
	for _, tagName := range postReq.Tags {
		_, err = tx.Exec(ctx, "INSERT INTO tags (post_id, name, created_at, updated_at) VALUES ($1, $2, $3, $4)",
			postID, tagName, time.Now(), time.Now())
		if err != nil {
			log.Println("Failed to insert into tags")
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &post.PostResponse{
		PostId: strconv.Itoa(postID), 
		FeedInHtml: postReq.FeedInHtml, 
		Tags: postReq.Tags,
	}, nil
}
