package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/user-service/internal/domain"
)

type UserUserRepo interface {
	GetFollowing(context.Context, string) ([]string, error)
	CreateFollowing(context.Context, string, string) error
	Unfollowing(context.Context, string, string) error
	CheckFollowing(context.Context, string, string) (bool, error)
}
type userUserRepo struct {
	db *goqu.Database
}

func NewUserUserRepository(
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
	// fmt.Println(query.ToSQL())
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

func (u userUserRepo) CreateFollowing(ctx context.Context, currentId string, userId string) error {
	query := u.db.
		Insert(domain.TabNameUserUser).
		Cols(domain.TabColUserId, domain.TabColFollower, domain.TabColCreatedAt, domain.TabColUpdatedAt).
		Vals(
			goqu.Vals{currentId, userId, time.Now(), time.Now()},
		)
	result, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (u userUserRepo) CheckFollowing(
	ctx context.Context,
	currentId string,
	userId string,
) (bool, error) {
	query := u.db.
		From(domain.TabNameUserUser).
		Where(
			goqu.C(domain.TabColUserId).Eq(currentId),
			goqu.C(domain.TabColFollower).Eq(userId),
		)

	var users []domain.UserUser
	err := query.Executor().ScanStructsContext(ctx, &users)
	if err != nil || len(users) == 0 {
		return false, err
	}
	if len(users) >= 1 {
		return true, nil
	}
	return false, nil
}

func (u userUserRepo) Unfollowing(
	ctx context.Context,
	currentId string,
	userId string,
) error {
	query := u.db.
		Delete(domain.TabNameUserUser).
		Where(
			goqu.C(domain.TabColUserId).Eq(currentId),
			goqu.C(domain.TabColFollower).Eq(userId),
		)

	fmt.Println(query.ToSQL())
	var users []domain.UserUser
	_, err := query.Executor().ExecContext(ctx)
	if err != nil || len(users) == 0 {
		return err
	}
	// if len(users) >= 1 {
	// 	return true, nil
	// }
	return nil
}
