package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/namnv2496/user-service/internal/cache"
	"github.com/namnv2496/user-service/internal/domain"
	es "github.com/namnv2496/user-service/internal/elasticsearch"
	userv1 "github.com/namnv2496/user-service/internal/handler/generated/user_core/v1"
	"github.com/namnv2496/user-service/internal/repo"
	"github.com/namnv2496/user-service/internal/security"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAccount(context.Context, string) (domain.User, error)
	FindAccount(context.Context, string) ([]domain.User, error)
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
	esClient     es.ElasticSearchClient
}

func NewUserService(
	userRepo repo.UserRepo,
	userUserRepo repo.UserUserRepo,
	redis cache.Client,
	esClient es.ElasticSearchClient,
) UserService {
	return &userService{
		userRepo:     userRepo,
		userUserRepo: userUserRepo,
		redis:        redis,
		esClient:     esClient,
	}
}

func (u userService) CreateAccount(
	ctx context.Context,
	req *userv1.Account,
) (uint64, error) {

	log.Println("Password: ", req.Password)
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := domain.User{
		Name:      req.Name,
		Email:     req.Email,
		UserId:    req.UserId,
		Password:  string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	log.Println("Called create new account: ", user)
	user, err := u.userRepo.CreateAccount(ctx, user)
	if err != nil {
		return 0, err
	}
	user.Password = ""
	// add to elastic search
	data := map[string]interface{}{
		"id":        user.Id,
		"userId":    user.UserId,
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	}

	err = u.esClient.CreateIndex(ctx, "user")
	if err != nil {
		log.Println("Failed to create index in elastic search: ", err)
		return 0, err
	}
	err = u.esClient.AddDataToIndex(ctx, "user", data)
	if err != nil {
		log.Println("Failed to add to elastic search: ", err)
		return 0, err
	}
	return user.Id, nil
}

func (u userService) GetAccount(
	ctx context.Context,
	userId string,
) (domain.User, error) {

	return u.userRepo.GetAccount(ctx, userId)
}

func (u userService) FindAccount(
	ctx context.Context,
	userId string,
) ([]domain.User, error) {
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"wildcard": map[string]interface{}{
	// 			"name": map[string]interface{}{
	// 				"value": userId + "*",
	// 			},
	// 		},
	// 	},
	// }
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match_phrase_prefix": map[string]interface{}{
							"name": map[string]interface{}{
								"query": userId,
							},
						},
					},
					map[string]interface{}{
						"match_phrase_prefix": map[string]interface{}{
							"userId": map[string]interface{}{
								"query": userId,
							},
						},
					},
				},
			},
		},
		"_source": []string{"id", "userId", "name", "email"},
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Fatalf("Error marshaling query: %s", err)
	}
	log.Println("Search with: ", string(queryJSON))
	result, err := u.esClient.SearchDataFromIndex(ctx, "user", string(queryJSON))
	if err != nil {
		return nil, err
	}
	res := make([]domain.User, 0)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {
		doc := hit.(map[string]interface{})
		source := doc["_source"].(map[string]interface{})

		// Convert the "id" field to uint64
		idStr := fmt.Sprintf("%v", source["id"])
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			log.Fatalf("Error converting id to uint64: %s", err)
		}

		// Map to User struct
		user := domain.User{
			Id: id,
			// Email:  source["email"].(string),
			Name:   source["name"].(string),
			UserId: source["userId"].(string),
		}
		res = append(res, user)
	}
	return res, nil
}

func (u userService) Login(
	ctx context.Context,
	userId string,
	password string,
) (string, error) {

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
