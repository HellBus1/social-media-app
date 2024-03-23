package services

import (
	"context"
	"log"
	"social-media-app/models"
	"social-media-app/models/post"
	"time"
	"fmt"

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

func LinkEmail(DB *pgxpool.Pool, userReq models.LinkEmailRequest, userID int) (*models.LinkEmailResponse, error) {
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

	result, err := tx.Exec(ctx, "UPDATE users SET email=$1, updated_at=$2 WHERE id=$3", userReq.Email, time.Now(), userID)
	if err != nil {
		fmt.Println("93")//
		fmt.Println(err)
		log.Println("Failed to update link email")
		return nil, err
	}
	rowsAffected := result.RowsAffected()
	if err != nil {
		fmt.Println("100")
		fmt.Println(err)
		log.Println("Failed to get the number of affected rows")
		return nil, err
	}
	log.Println("rowsAffected: ", rowsAffected)//

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.LinkEmailResponse{
		Email: userReq.Email,
	}, nil
}

func LinkPhone(DB *pgxpool.Pool, userReq models.LinkPhoneRequest, userID int) (*models.LinkPhoneResponse, error) {
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

	result, err := tx.Exec(ctx, "UPDATE users SET phone=$1, updated_at=$2 WHERE id=$3", userReq.Phone, time.Now(), userID)
	if err != nil {
		log.Println("Failed to update link phone")
		return nil, err
	}
	rowsAffected := result.RowsAffected()
	if err != nil {
		log.Println("Failed to get the number of affected rows")
		return nil, err
	}
	log.Println("rowsAffected: ", rowsAffected)//

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.LinkPhoneResponse{
		Phone: userReq.Phone,
	}, nil
}
