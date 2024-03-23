package services

import (
	"context"
	"log"
	"social-media-app/models/tag"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTagsByPostId(DB *pgxpool.Pool, postId int) (*[]string, error)  {
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

	var tags []string
	rows, err := DB.Query(ctx, `SELECT t.name FROM tags t WHERE t.post_id = $1`, postId)
	if err != nil {
					log.Println("Failed to query comments:", err)
					return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
					var tag tag.TagResponse
					err := rows.Scan(&tag.Name)
					if err != nil {
									log.Println("Failed to scan row:", err)
									continue
					}
	
					tags = append(tags, tag.Name)
	}
	
	if err := rows.Err(); err != nil {
					log.Println("Error while iterating rows:", err)
					return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit transaction:", err)
			return nil, err
	}

	return &tags, nil
}