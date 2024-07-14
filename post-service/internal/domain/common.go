package domain

import "github.com/doug-martin/goqu/v9"

var (
	TabNameTopComment = goqu.T("top_comments")

	TabColId            = "id"
	TabColUserId        = "user_id"
	TabColPostId        = "post_id"
	TabColCommentLevel  = "comment_level"
	TabColCommentParent = "comment_parent"
	TabColCreatedAt     = "created_at"
	TabColUpdatedAt     = "updated_at"
	TabColLike          = "like"
	TabColTotalLike     = "total_like"
)
