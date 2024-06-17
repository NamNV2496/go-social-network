package repo

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/user-service/internal/domain"
)

type UserUserRepo interface {
	GetFollowing(context.Context, string) ([]string, error)
}
type userUserRepo struct {
	db *goqu.Database
}

func NewUserUserService(
	db *goqu.Database,
) UserUserRepo {
	return &userUserRepo{
		db: db,
	}
}

func (u userUserRepo) GetFollowing(ctx context.Context, userId string) ([]string, error) {
	query := u.db.
		From(domain.TabNameUserUser).
		Where(
			goqu.C(domain.TabColFollower).Eq(userId),
		)
	fmt.Println(query.ToSQL())
	var users []domain.UserUser
	err := query.Executor().ScanStructsContext(ctx, &users)
	if err != nil || len(users) == 0 {
		return []string{}, err
	}
	var res []string
	for _, user := range users {
		res = append(res, user.UserId)
	}
	return res, nil
}
