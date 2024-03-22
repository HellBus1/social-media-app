package post

import "time"

type PostResponse struct {
	PostId     string     `json:"postId"`
	FeedInHtml string     `json:"feedInHtml"`
	Tags       []string   `json:"tags"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

type PostPaginatedResponse struct {
	PostId   string                `json:"postId"`
	Post     PostResponse          `json:"post"`
	Comments []PostCommentResponse `json:"comments"`
	Creator  PostCreatorResponse   `json:"creator"`
}

type PostCommentResponse struct {
	Comment 	string								`json:"comment"`
	Creator 	PostCreatorResponse		`json:"creator"`
}

type PostCreatorResponse struct {
	UserId 			string 		`json:"userId,omitempty"`
	Name 				string 		`json:"name"`
	ImageUrl 		string 		`json:"imageUrl"`
	FriendCount int 			`json:"friendCount"`
	CreatedAt 	time.Time `json:"createdAt"`
}
