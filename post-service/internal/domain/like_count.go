package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameLikeCount = goqu.T("like_count")
)

type LikeCount struct {
	Id        uint64    `db:"id" goqu:"omitnil"`
	PostId    uint64    `db:"post_id" goqu:"omitnil"`
	TotalLike int       `db:"total_like" goqu:"omitnil"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil"`
}
