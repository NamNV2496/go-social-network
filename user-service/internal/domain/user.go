package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameUser = goqu.T("user")
)

type User struct {
	Id        uint64    `db:"id" goqu:"omitnil"`
	Email     string    `db:"email" goqu:"omitnil"`
	Name      string    `db:"name" goqu:"omitnil"`
	Picture   string    `db:"picture" goqu:"omitnil"`
	UserId    string    `db:"user_id" goqu:"omitnil"`
	Password  string    `db:"password" goqu:"omitnil"`
	CreatedAt time.Time `db:"created_at" goqu:"omitnil"`
	UpdatedAt time.Time `db:"updated_at" goqu:"omitnil"`
}
