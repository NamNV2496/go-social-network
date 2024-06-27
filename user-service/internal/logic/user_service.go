package logic

import (
	"context"
	"log"

	"github.com/namnv2496/user-service/internal/cache"
	"github.com/namnv2496/user-service/internal/domain"
	userv1 "github.com/namnv2496/user-service/internal/handler/generated/user_core/v1"
	"github.com/namnv2496/user-service/internal/repo"
	"github.com/namnv2496/user-service/internal/security"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAccount(context.Context, string) (domain.User, error)
	CreateAccount(context.Context, *userv1.Account) (uint64, error)
	Login(context.Context, string, string) (string, error)
	GetFollowing(context.Context, string) ([]string, error)
	CreateFollowing(context.Context, string, string) (bool, error)
	Unfollowing(context.Context, string, string) error
	CheckFollowing(context.Context, string, string) (bool, error)
}
type userService struct {
	userRepo     repo.UserRepo
	userUserRepo repo.UserUserRepo
	redis        cache.Client
}

func NewUserService(
	userRepo repo.UserRepo,
	userUserRepo repo.UserUserRepo,
	redis cache.Client,
) UserService {
	return &userService{
		userRepo:     userRepo,
		userUserRepo: userUserRepo,
		redis:        redis,
	}
}

func (u userService) CreateAccount(ctx context.Context, user *userv1.Account) (uint64, error) {

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	log.Println(string(passwordHash))
	return 1, nil
}

func (u userService) GetAccount(ctx context.Context, userId string) (domain.User, error) {

	return u.userRepo.GetAccount(ctx, userId)
}

func (u userService) Login(ctx context.Context, userId string, password string) (string, error) {

	user, err := u.userRepo.GetAccount(ctx, userId)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil
	}
	value, exist := u.redis.Get(ctx, userId+"_token")
	if exist == nil {
		return value.(string), nil
	}
	token, err := security.GenerateJWTToken(userId)
	if err != nil {
		return "", err
	}
	if err := u.redis.Set(ctx, userId+"_token", token); err != nil {
		log.Fatalln("Failed to save token to redis")
	}
	return token, err
}

func (u userService) GetFollowing(
	ctx context.Context,
	userId string,
) ([]string, error) {

	return u.userUserRepo.GetFollowing(ctx, userId)
}

func (u userService) CreateFollowing(
	ctx context.Context,
	currentId string,
	userId string,
) (bool, error) {
	log.Println("Add new following: ", currentId, " - ", userId)
	err := u.userUserRepo.CreateFollowing(ctx, currentId, userId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u userService) CheckFollowing(
	ctx context.Context,
	currentId string,
	userId string,
) (bool, error) {

	exist, err := u.userUserRepo.CheckFollowing(ctx, currentId, userId)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (u userService) Unfollowing(
	ctx context.Context,
	currentId string,
	userId string,
) error {

	err := u.userUserRepo.Unfollowing(ctx, currentId, userId)
	if err != nil {
		return err
	}
	return nil
}
