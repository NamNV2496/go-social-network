package logic

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/user-service/internal/cache"
	"github.com/namnv2496/user-service/internal/domain"
	userv1 "github.com/namnv2496/user-service/internal/handler/generated/user_core/v1"
	"github.com/namnv2496/user-service/internal/security"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAccount(context.Context, string) (domain.User, error)
	CreateAccount(context.Context, userv1.Account) (uint64, error)
	Login(context.Context, string, string) (string, error)
}
type userService struct {
	db    *goqu.Database
	redis cache.Client
}

func NewUserService(
	db *goqu.Database,
	redis cache.Client,
) UserService {
	return &userService{
		db:    db,
		redis: redis,
	}
}

func (u userService) CreateAccount(ctx context.Context, user userv1.Account) (uint64, error) {

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.UserId), bcrypt.DefaultCost)
	fmt.Println(string(passwordHash))
	return 1, nil
}

func (u userService) GetAccount(ctx context.Context, userId string) (domain.User, error) {

	query := u.db.
		From(domain.TabNameUser).
		Where(
			goqu.C(domain.TabColUserId).Eq(userId),
		)
	fmt.Println(query.ToSQL())
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

func (u userService) Login(ctx context.Context, userId string, password string) (string, error) {

	query := u.db.
		From(domain.TabNameUser).
		Where(
			goqu.C(domain.TabColUserId).Eq(userId),
		)
	var users []domain.User
	err := query.Executor().ScanStructsContext(ctx, &users)
	if err != nil || len(users) == 0 {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(password))
	if err != nil {
		return "", nil
	}
	value, exist := u.redis.Get(ctx, userId)
	if exist == nil {
		return value.(string), nil
	}
	token, err := security.GenerateJWTToken(userId)
	if err != nil {
		return "", err
	}
	u.redis.Set(ctx, userId, token)
	return token, err
}
