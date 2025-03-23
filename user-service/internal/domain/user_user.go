package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameUserUser = goqu.T("user_user")
)

type UserUser struct {
	Id        uint64    `db:"id" goqu:"omitnil" esMapping:"key:id,type:-"`
	UserId    string    `db:"user_id" goqu:"omitnil" esMapping:"key:user_id,type:string"`
	Follower  string    `db:"follower" goqu:"omitnil" esMapping:"key:follower,type:string"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil" esMapping:"key:created_at,type:date,format:yyyy-MM-dd HH:mm:ss"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil" esMapping:"key:updated_at,type:date,format:yyyy-MM-dd HH:mm:ss"`
}
