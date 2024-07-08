package repo

import (
	"context"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/user-service/internal/domain"
)

type UserRepo interface {
	GetAccount(ctx context.Context, userId string) (domain.User, error)
	CreateAccount(ctx context.Context, user domain.User) (domain.User, error)
}
type userRepo struct {
	db *goqu.Database
}

func NewUserService(
	db *goqu.Database,
) UserRepo {
	return &userRepo{
		db: db,
	}
}

func (u userRepo) GetAccount(ctx context.Context, userId string) (domain.User, error) {
	query := u.db.
		From(domain.TabNameUser).
		Where(
			goqu.C(domain.TabColUserId).Eq(userId),
		)
	// fmt.Println(query.ToSQL())
	var users []domain.User
	err := query.Executor().ScanStructsContext(ctx, &users)
	if err != nil {
		return domain.User{}, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return domain.User{}, err
}
func (u userRepo) CreateAccount(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	query := u.db.
		Insert(domain.TabNameUser).
		Rows(user)
	_, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return domain.User{}, err
	}

	var newUser domain.User
	getUserQuery := u.db.From(domain.TabNameUser).Where(goqu.C(domain.TabColUserId).Eq(user.UserId))
	_, err = getUserQuery.Executor().ScanStructContext(ctx, &newUser)
	if err != nil {
		return domain.User{}, err
	}

	log.Println("User: ", newUser)
	return newUser, nil
}
