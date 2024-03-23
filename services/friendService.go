package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"social-media-app/models/friend"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateFriendService(DB *pgxpool.Pool, Request friend.FriendRequest, userId int) error {
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

	// Check if the friendship exists
	var friendshipExists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM friendship WHERE user_id = $1 AND friend_id = $2 OR user_id = $2 AND friend_id = $1)", userId, Request.UserID).Scan(&friendshipExists)
	if err != nil {
		log.Println("Error executing select statement:", err)
		return err
	}

	if friendshipExists {
		log.Println("Friendship has already exist")
		return errors.New("friendship has already exist")
	}

	// Insert the first friendship (user A adding user B)
	_, err = tx.Exec(ctx, "INSERT INTO friendship (user_id, friend_id, created_at, updated_at) VALUES ($1, $2, $3, $4)", userId, Request.UserID, time.Now(), time.Now())
	if err != nil {
		log.Println("Error executing insert statement:", err)
		return err
	}

	// Insert the second friendship (user B adding user A)
	_, err = tx.Exec(ctx, "INSERT INTO friendship (user_id, friend_id, created_at, updated_at) VALUES ($1, $2, $3, $4)", Request.UserID, userId, time.Now(), time.Now())
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

func GetListOfFriendService(DB *pgxpool.Pool, userId int, pageSize int, offset int) ([]map[string]interface{}, error) {
	ctx := context.Background()

	// Query to fetch list of friends along with friend count for each user
	query := `
        SELECT u.id AS user_id, u.name, u.image_url, COUNT(DISTINCT f.id) AS friend_count, u.created_at
        FROM users u
        JOIN friendship f ON u.id = f.friend_id
        WHERE f.user_id = $1
        GROUP BY u.id
        ORDER BY u.created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := DB.Query(ctx, query, userId, pageSize, offset)
	if err != nil {
		log.Println("Error executing SQL query:", err)
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	// Create slice to store result data
	var listFriendData []map[string]interface{}

	// Iterate over rows
	for rows.Next() {
		var userData = make(map[string]interface{})
		var userId int
		var name, imageUrl string
		var friendCount int
		var createdAt time.Time

		err := rows.Scan(&userId, &name, &imageUrl, &friendCount, &createdAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		userData["user_id"] = userId
		userData["name"] = name
		userData["image_url"] = imageUrl
		userData["friend_count"] = friendCount
		userData["created_at"] = createdAt

		listFriendData = append(listFriendData, userData)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over rows:", err)
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return listFriendData, nil
}

func RemoveFriendService(DB *pgxpool.Pool, userId int, Request friend.FriendRequest) error {
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

	// Check if the friendship exists
	var friendshipExists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM friendship WHERE user_id = $1 AND friend_id = $2 OR user_id = $2 AND friend_id = $1)", userId, Request.UserID).Scan(&friendshipExists)
	if err != nil {
		log.Println("Error executing select statement:", err)
		return err
	}

	if !friendshipExists {
		log.Println("Friendship does not exist")
		return errors.New("friendship does not exist")
	}

	// Remove the first friendship (user A removing user B)
	_, err = tx.Exec(ctx, "DELETE FROM friendship WHERE user_id = $1 AND friend_id = $2", userId, Request.UserID)
	if err != nil {
		log.Println("Error executing delete statement:", err)
		return err
	}

	// Remove the second friendship (user B removing user A)
	_, err = tx.Exec(ctx, "DELETE FROM friendship WHERE user_id = $1 AND friend_id = $2", Request.UserID, userId)
	if err != nil {
		log.Println("Error executing delete statement:", err)
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}
