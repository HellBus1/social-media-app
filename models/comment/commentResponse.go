package comment

import "time"

type CommentResponse struct {
	ID            int       `json:"id"`
	PostID        int       `json:"postId"`
	CommentatorID int       `json:"commentatorId"`
	Comment       string    `json:"comment"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
