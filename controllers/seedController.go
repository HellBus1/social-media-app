package controllers

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID         int
	Name       string
	Password   string
	Email      string
	Phone      string
	ImageURL   string
	Credential string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Post struct {
	UserID     int
	PostInHTML string
	CreatedAt  time.Time
}

type Comment struct {
	PostID        int
	CommentatorID int
	Comment       string
	CreatedAt     time.Time
}

type Tag struct {
	PostID      int
	Name        string
	CreatedAt   time.Time
}

func CreateSeed(ginContext *gin.Context) {
	DB, ok := ginContext.MustGet("DB").(*pgxpool.Pool)
	if !ok {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get DB from context"})
		return
	}
	// Seed users
	users := []*User{
		{ID: 1, Name: "User1", Password: "password1", Email: "user1@example.com", Phone: "1234567890", ImageURL: "image1.jpg", Credential: "credential1"},
		{ID: 2, Name: "User2", Password: "password2", Email: "user2@example.com", Phone: "1234567891", ImageURL: "image2.jpg", Credential: "credential2"},
		{ID: 3, Name: "User3", Password: "password2", Email: "user3@example.com", Phone: "1234567892", ImageURL: "image2.jpg", Credential: "credential2"},
		{ID: 4, Name: "User4", Password: "password2", Email: "user4@example.com", Phone: "1234567893", ImageURL: "image2.jpg", Credential: "credential2"},
		{ID: 5, Name: "User5", Password: "password2", Email: "user5@example.com", Phone: "1234567894", ImageURL: "image2.jpg", Credential: "credential2"},
		{ID: 6, Name: "User6", Password: "password2", Email: "user6@example.com", Phone: "1234567895", ImageURL: "image2.jpg", Credential: "credential2"},
		{ID: 7, Name: "User7", Password: "password2", Email: "user7@example.com", Phone: "1234567896", ImageURL: "image2.jpg", Credential: "credential2"},
		// Add more users here
	}

	for _, u := range users {
			_, err := DB.Exec(context.Background(), `
					INSERT INTO users (id, name, password, email, phone, image_url, credential_type, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, u.ID, u.Name, u.Password, u.Email, u.Phone, u.ImageURL, u.Credential, time.Now(), time.Now())
			if err != nil {
					log.Fatalf("Error inserting user %s: %v\n", u.Name, err)
			}
	}

	// Seed friendship (assuming everyone is friend with everyone else except self)
	for _, user := range users {
			for _, friend := range users {
					if user.ID != friend.ID {
							_, err := DB.Exec(context.Background(), `
									INSERT INTO friendship (user_id, friend_id, created_at, updated_at)
									VALUES ($1, $2, $3, $4)
							`, user.ID, friend.ID, time.Now(), time.Now())
							if err != nil {
									log.Fatalf("Error inserting friendship between %s and %s: %v\n", user.Name, friend.Name, err)
							}
					}
			}
	}

	// Seed posts (for users with ID 1 and 2)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 30; i++ {
			userID := 1 // Assuming the posts are created by user with ID 1 and 2 alternately
			if i%2 == 1 {
					userID = 2
			}
			post := &Post{
					UserID:     userID,
					PostInHTML: "Sample post content",
					CreatedAt:  time.Now().Add(-time.Duration(rand.Intn(100)) * time.Hour * 24), // Random past date
			}
			_, err := DB.Exec(context.Background(), `
					INSERT INTO posts (user_id, post_in_html, created_at, updated_at)
					VALUES ($1, $2, $3, $4)
			`, post.UserID, post.PostInHTML, post.CreatedAt, post.CreatedAt)
			if err != nil {
					log.Fatalf("Error inserting post for user %d: %v\n", post.UserID, err)
			}

			// Seed tags for each post
			numTags := rand.Intn(3) + 1 // Random number of tags (1 to 3)
			for j := 0; j < numTags; j++ {
					tag := &Tag{
							PostID:    i + 1,
							Name:      "Tag" + strconv.Itoa(j+1),
							CreatedAt: time.Now().Add(-time.Duration(rand.Intn(100)) * time.Hour * 24), // Random past date
					}
					_, err := DB.Exec(context.Background(), `
							INSERT INTO tags (post_id, name, created_at, updated_at)
							VALUES ($1, $2, $3, $4)
					`, tag.PostID, tag.Name, tag.CreatedAt, tag.CreatedAt)
					if err != nil {
							log.Fatalf("Error inserting tag for post %d: %v\n", tag.PostID, err)
					}
			}
	}

	// Seed comments
	for i := 0; i < 30; i++ {
			postID := i%15 + 1 // Assuming there are 30 posts and 10 users, so each post may have 2-3 comments
			commentatorID := (i/15 + 2) % 10 + 1 // Assuming commentator ID starts from 3 to 10
			comment := &Comment{
					PostID:         postID,
					CommentatorID:  commentatorID,
					Comment:        "Sample comment content",
					CreatedAt:      time.Now().Add(-time.Duration(rand.Intn(100)) * time.Hour * 24), // Random past date
			}
			_, err := DB.Exec(context.Background(), `
					INSERT INTO comments (post_id, commentator_id, comment, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5)
			`, comment.PostID, comment.CommentatorID, comment.Comment, comment.CreatedAt, comment.CreatedAt)
			if err != nil {
					log.Fatalf("Error inserting comment for post %d by user %d: %v\n", comment.PostID, comment.CommentatorID, err)
			}
	}

	log.Println("Seed data generated successfully!")

	ginContext.JSON(http.StatusAccepted, gin.H{"message": "successfully add post"})
}
