package services

import (
	"context"
	"log"
	"social-media-app/models/post"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserByIdWithFriendCount(DB *pgxpool.Pool, userID int) (*post.PostCreatorResponse, error) {
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

	var user post.PostCreatorResponse
	err = DB.QueryRow(ctx, `
			SELECT u.id, u.name, u.image_url, COUNT(f.friend_id) AS friend_count, u.created_at 
			FROM users u
			LEFT JOIN friendship f ON u.id = f.user_id
			WHERE u.id = $1
			GROUP BY u.id, u.name, u.image_url, u.created_at
	`, userID).
			Scan(&user.UserId, &user.Name, &user.ImageUrl, &user.FriendCount, &user.CreatedAt)
	if err != nil {
			log.Println("Failed to query user:", err)
			return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &user, nil
}