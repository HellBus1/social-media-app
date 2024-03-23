package services

import (
	"context"
	"log"
	"social-media-app/models/post"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCommentByPostId(DB *pgxpool.Pool, postId int) (*[]post.PostCommentResponse, error) {
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

	var comments []post.PostCommentResponse
	rows, err := DB.Query(ctx, `
			SELECT c.comment, u.id, u.name, u.image_url, COUNT(f.friend_id) AS friend_count, c.created_at 
			FROM comments c
			INNER JOIN users u ON c.commentator_id = u.id
			LEFT JOIN friendship f ON u.id = f.user_id
			WHERE c.post_id = $1
			GROUP BY c.comment, u.id, u.name, u.image_url, c.created_at
			ORDER BY c.created_at
	`, postId)
	if err != nil {
		log.Println("Failed to query comments:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment string
		var creator post.PostCreatorResponse
		err := rows.Scan(&comment, &creator.UserId, &creator.Name, &creator.ImageUrl, &creator.FriendCount, &creator.CreatedAt)
		if err != nil {
			log.Println("Failed to scan row:", err)
			continue
		}

		comments = append(comments, post.PostCommentResponse{
			Comment: comment,
			Creator: creator,
		})
	}

	if err := rows.Err(); err != nil {
		log.Println("Error while iterating rows:", err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &comments, nil
}