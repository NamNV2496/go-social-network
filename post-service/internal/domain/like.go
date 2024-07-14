package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameLike = goqu.T("like")
)

type Like struct {
	Id        uint64    `db:"id" goqu:"omitnil"`
	PostId    uint64    `db:"post_id" goqu:"omitnil"`
	UserId    string    `db:"user_id" goqu:"omitnil"`
	Like      bool      `db:"like" goqu:"omitnil"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil"`
}
