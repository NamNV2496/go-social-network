package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/newsfeed-service/internal/cache"
	"github.com/namnv2496/newsfeed-service/internal/domain"
	newsfeedv1 "github.com/namnv2496/newsfeed-service/internal/handler/generated/newsfeed_core/v1"
	"github.com/namnv2496/newsfeed-service/internal/mq/consumer"
)

type NewsfeedService interface {
	UpdateNewsfeed(context.Context, []string, domain.Post) error
	GetNewsfeed(context.Context, string) (*newsfeedv1.GetNewsfeedResponse, error)
}

type newsfeedService struct {
	db          *goqu.Database
	kafkaClient consumer.Consumer
	redis       cache.Client
}

func NewUserService(
	db *goqu.Database,
	kafkaClient consumer.Consumer,
	redis cache.Client,
) NewsfeedService {
	return &newsfeedService{
		db:          db,
		kafkaClient: kafkaClient,
		redis:       redis,
	}
}

func (s newsfeedService) UpdateNewsfeed(
	ctx context.Context,
	users []string,
	newPost domain.Post,
) error {

	if len(users) == 0 {
		return nil
	}
	for _, user := range users {
		var jsonData []domain.Post
		data, exist := s.redis.Get(ctx, user)
		if exist != nil {
			fmt.Println("cache post is not exist. Add new! user_id: ", user)
			jsonData = append(jsonData, newPost)
			redis, err := json.Marshal(jsonData)
			if err != nil {
				return err
			}
			s.redis.Set(ctx, user, redis)
			continue
		}
		fmt.Println("cache post is exist! user_id: ", user)
		dataByte, _ := data.(string)
		if err := json.Unmarshal([]byte(dataByte), &jsonData); err != nil {
			fmt.Println("error when marshal old post")
			return err
		}
		jsonData = append(jsonData, newPost)
		redis, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}
		s.redis.Set(ctx, user, redis)
	}
	return nil
}

func (s newsfeedService) GetNewsfeed(ctx context.Context, userId string) (*newsfeedv1.GetNewsfeedResponse, error) {
	data, exist := s.redis.Get(ctx, userId)
	if exist == nil {
		var posts []newsfeedv1.NewsfeedPost
		if err := json.Unmarshal([]byte(data.(string)), &posts); err != nil {
			fmt.Println("error when marshal old post")
			return nil, err
		}
		postPointers := make([]*newsfeedv1.NewsfeedPost, len(posts))
		for i := range posts {
			postPointers[i] = &posts[i]
		}
		return &newsfeedv1.GetNewsfeedResponse{
			Posts: postPointers,
		}, nil
	}
	return nil, nil
}
