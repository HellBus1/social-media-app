package post

type PostResponse struct {
	PostId		 string		`json:"postId"`
	FeedInHtml string   `json:"feedInHtml"`
	Tags       []string `json:"tags"`
}
