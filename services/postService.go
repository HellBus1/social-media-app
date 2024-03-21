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

func GetPostsByUserId(DB *pgxpool.Pool, userID int) (*[]post.PostPaginatedResponse, error) {
	log.Println()
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

	var posts []post.PostPaginatedResponse
	rows, err := DB.Query(ctx, `
					SELECT p.id, p.post_in_html, p.created_at,
								 u.id, u.name, u.image_url, COUNT(f.friend_id) AS friend_count, p.created_at, p.updated_at
					FROM posts p
					INNER JOIN users u ON p.user_id = u.id
					LEFT JOIN friendship f ON u.id = f.user_id
					WHERE p.user_id = $1
					GROUP BY p.id, u.id, u.name, u.image_url, u.created_at
					ORDER BY p.created_at
	`, userID)
	if err != nil {
					log.Println("Failed to query posts:", err)
					return nil, err
	}
	defer rows.Close()

	user, getUserError := GetUserByIdWithFriendCount(DB, userID);
	if (getUserError != nil) {
		return nil, getUserError
	}

	for rows.Next() {
			var postItem post.PostPaginatedResponse
			var tempPostId int
			err := rows.Scan(&tempPostId, &postItem.Post.FeedInHtml, &postItem.Post.CreatedAt,
											 &postItem.Creator.UserId, &postItem.Creator.Name, &postItem.Creator.ImageUrl,
											 &postItem.Creator.FriendCount, &postItem.Creator.CreatedAt, &postItem.Post.UpdatedAt)
			if err != nil {
					log.Println("Failed to scan row:", err)
					continue
			}

			comments, getCommentsError := GetCommentByPostId(DB, tempPostId)
			if getCommentsError != nil {
					log.Println("Failed to get comments for post:", getCommentsError)
					continue
			}

			tags, getTagsError := GetTagsByPostId(DB, tempPostId)
			if getTagsError != nil {
				log.Println("Failed to get tags for post:", getTagsError)
				continue
			}

			postItem.PostId = strconv.Itoa(tempPostId)
			postItem.Comments = append(postItem.Comments, *comments...)
			postItem.Creator = *user
			postItem.Post.Tags = *tags

			posts = append(posts, postItem)
	}

	if err := rows.Err(); err != nil {
					log.Println("Error while iterating rows:", err)
					return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
			log.Println("Failed to commit transaction:", err)
			return nil, err
	}

	return &posts, nil
}