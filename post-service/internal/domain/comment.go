package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameComment = goqu.T("comment")
)

type Comment struct {
	Id            uint64    `db:"id" goqu:"omitnil"`
	PostId        uint64    `db:"post_id" goqu:"omitnil"`
	UserId        string    `db:"user_id" goqu:"omitnil"`
	CommentText   string    `db:"comment_text" goqu:"omitnil"`
	CommentLevel  int       `db:"comment_level" goqu:"omitnil"`
	CommentParent uint64    `db:"comment_parent" goqu:"omitnil"`
	Images        string    `db:"images" goqu:"omitnil"`
	Tags          string    `db:"tags" goqu:"omitnil"`
	CreatedAt     time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt     time.Time `db:"updated_at" goqu:"omitnil"`
}
