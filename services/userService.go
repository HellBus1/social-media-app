package services

import (
	"context"
	"log"
	"social-media-app/models"
	"social-media-app/models/post"
	"time"
	"fmt"
    "net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
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
	err = tx.QueryRow(ctx, `
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

func GetUserByUserId(DB *pgxpool.Pool, userID int) (*models.Users, error) {
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

	var user models.Users
	var nullableImageUrl *string
	var nullableEmail *string
	var nullablePhone *string
	err = tx.QueryRow(ctx, `
			SELECT id, name, password, email, phone, image_url, credential_type, created_at, updated_at
			FROM users
			WHERE id = $1
	`, userID).
			Scan(&user.ID, &user.Name, &user.Password, &nullableEmail, &nullablePhone, &nullableImageUrl, &user.CredentialType, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
			log.Println("Failed to query user by ID:", err)
			return nil, err
	}

	if nullableImageUrl != nil {
		user.ImageURL = *nullableImageUrl
	}

	if nullableEmail != nil {
		user.Email = *nullableEmail
	}

	if nullablePhone != nil {
		user.Phone = *nullablePhone
	}

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &user, nil
}

func LinkEmail(DB *pgxpool.Pool, userReq models.LinkEmailRequest, userID int, cc *gin.Context) (*models.LinkEmailResponse, error) {
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

	// Check if email is already filled
	var emailFilled bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND email IS NOT NULL)", userID).Scan(&emailFilled)
	if err != nil {
		fmt.Println("Error checking if email is filled:", err)
		return nil, err
	}
	if emailFilled {
		// Email is already filled, return with an appropriate message
		cc.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Email is already filled for this user"})

		return nil, nil
	}

	// Proceed with the update if email is not filled
	result, err := tx.Exec(ctx, "UPDATE users SET email=$1, updated_at=$2 WHERE id=$3 AND email IS NULL", userReq.Email, time.Now(), userID)
	if err != nil {
		fmt.Println("Failed to update link email:", err)
		return nil, err
	}
	rowsAffected := result.RowsAffected()
	log.Println("Rows affected: ", rowsAffected)

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.LinkEmailResponse{
		Email: userReq.Email,
	}, nil
}

func LinkPhone(DB *pgxpool.Pool, userReq models.LinkPhoneRequest, userID int, cc *gin.Context) (*models.LinkPhoneResponse, error) {
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

	// Check if phone is already filled
	var phoneFilled bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND phone IS NOT NULL)", userID).Scan(&phoneFilled)
	if err != nil {
		fmt.Println("Error checking if phone is filled:", err)
		return nil, err
	}
	if phoneFilled {
		// Phone is already filled, return with an appropriate message
		cc.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "message": "Phone is already filled for this user"})

		return nil, nil
	}

	// Proceed with the update if phone is not filled
	result, err := tx.Exec(ctx, "UPDATE users SET phone=$1, updated_at=$2 WHERE id=$3 AND phone IS NULL", userReq.Phone, time.Now(), userID)
	if err != nil {
		fmt.Println("Failed to update phone:", err)
		return nil, err
	}
	rowsAffected := result.RowsAffected()
	log.Println("Rows affected: ", rowsAffected)

	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return nil, err
	}

	return &models.LinkPhoneResponse{
		Phone: userReq.Phone,
	}, nil
}