package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameUserUser = goqu.T("user_user")
)

type UserUser struct {
	Id        uint64    `db:"id" goqu:"omitnil"`
	UserId    string    `db:"user_id" goqu:"omitnil"`
	Follower  string    `db:"follower_id" goqu:"omitnil"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil"`
}
