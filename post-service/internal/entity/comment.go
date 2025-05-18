package entity

type Comment struct {
	UserId        string   `json:"user_id,omitempty"`
	CommentId     uint64   `json:"comment_id,omitempty"`
	CommentText   string   `json:"comment_text,omitempty"`
	CommentLevel  uint32   `json:"comment_level,omitempty"`
	CommentParent uint64   `json:"comment_parent,omitempty"`
	Images        []string `json:"images,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Date          string   ` json:"date,omitempty"`
}
type CreateCommentRequest struct {
	PostId      uint64   `json:"post_id,omitempty"`
	Comment     *Comment `json:"comment,omitempty"`
	Application []string `json:"application,omitempty"`
}

type CreateCommentResponse struct {
	CommentId uint64 `json:"comment_id,omitempty"`
}

type GetCommentRequest struct {
	PostId     uint64 `json:"post_id,omitempty"`
	PageNumber uint32 ` json:"page_number,omitempty"`
	PageSize   uint32 `json:"page_size,omitempty"`
}

type GetCommentResponse struct {
	Comment []*Comment `json:"comment,omitempty"`
}
