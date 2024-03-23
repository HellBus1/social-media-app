package post

type PostRequest struct {
	FeedInHtml string   `json:"feedInHtml" binding:"required,min=2,max=500" validate:"required,min=2,max=500"`
	Tags       []string `json:"tags" binding:"required" validate:"required"`
}
