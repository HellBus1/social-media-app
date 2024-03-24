package services

import (
	"context"
	"log"
	"social-media-app/models/updateAccount"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateAccountService(DB *pgxpool.Pool, Request updateaccount.UpdateAccountRequest, userId int) error {
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

	_, err = tx.Exec(ctx, "UPDATE users SET image_url=$1, name=$2, updated_at=$3 WHERE id=$4", Request.ImageURL, Request.Name, time.Now(), userId)
	if err != nil {
		log.Println("Error executing insert statement:", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}
