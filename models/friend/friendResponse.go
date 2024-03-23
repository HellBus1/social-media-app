package friend

import "time"

type FriendResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	FriendID  int       `json:"friendId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
