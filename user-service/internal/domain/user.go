package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameUser = goqu.T("user")
)

type User struct {
	Id        int64     `db:"id" goqu:"omitnil" esMapping:"key:id,type:-"`
	Email     string    `db:"email" goqu:"omitnil" esMapping:"key:email,type:string"`
	Name      string    `db:"name" goqu:"omitnil" esMapping:"key:name,type:string"`
	Picture   string    `db:"picture" goqu:"omitnil" esMapping:"key:picture,type:-"`
	UserId    string    `db:"user_id" goqu:"omitnil" esMapping:"key:user_id,type:string"`
	Password  string    `db:"password" goqu:"omitnil" esMapping:"key:password,type:-"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil" esMapping:"key:created_at,type:date,format:yyyy-MM-dd HH:mm:ss"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil" esMapping:"key:updated_at,type:date,format:yyyy-MM-dd HH:mm:ss"`
}
