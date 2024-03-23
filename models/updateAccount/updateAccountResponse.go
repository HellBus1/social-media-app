package updateaccount

import "time"

type UpdateAccountResponse struct {
	ID        int       `json:"id"`
	ImageURL  string    `json:"imageUrl"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
}
