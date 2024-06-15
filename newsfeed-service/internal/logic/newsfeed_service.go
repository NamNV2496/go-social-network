package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/newsfeed-service/internal/cache"
	"github.com/namnv2496/newsfeed-service/internal/domain"
	"github.com/namnv2496/newsfeed-service/internal/mq/consumer"
)

type NewsfeedService interface {
	UpdateNewsfeed(context.Context, []byte) error
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

func (s newsfeedService) UpdateNewsfeed(ctx context.Context, payload []byte) error {

	var newPost domain.Post
	if err := json.Unmarshal([]byte(payload), &newPost); err != nil {
		fmt.Println("error when marshal")
		return err
	}
	var jsonData []domain.Post
	data, exist := s.redis.Get(ctx, newPost.User_id)
	if exist != nil {
		jsonData = append(jsonData, newPost)
		redis, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}
		s.redis.Set(ctx, newPost.User_id, redis)
		return nil
	}
	dataByte, _ := data.(string)
	if err := json.Unmarshal([]byte(dataByte), &jsonData); err != nil {
		fmt.Println("error when marshal")
		return err
	}
	jsonData = append(jsonData, newPost)
	redis, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}
	s.redis.Set(ctx, newPost.User_id, redis)
	return nil
}
