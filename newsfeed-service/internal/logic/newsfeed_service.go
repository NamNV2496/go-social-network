package logic

import (
	"context"
	"encoding/json"
	"log"
	"strings"

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
		user = strings.ToLower(user)
		user = strings.Trim(user, " ")
		var jsonData []domain.Post
		data, exist := s.redis.Get(ctx, user)
		if exist != nil {
			log.Println("cache post is not exist. Add new! user_id: ", user)
			jsonData = append(jsonData, newPost)
			redis, err := json.Marshal(jsonData)
			if err != nil {
				return err
			}
			if err := s.redis.Set(ctx, user, redis); err != nil {
				log.Fatalln("Failed to set data to redis")
			}
			continue
		}
		log.Println("cache post is exist! user_id: ", user)
		dataByte, _ := data.(string)
		if err := json.Unmarshal([]byte(dataByte), &jsonData); err != nil {
			log.Println("error when marshal old post")
			return err
		}
		jsonData = append(jsonData, newPost)
		redis, err := json.Marshal(jsonData)
		if err != nil {
			return err
		}
		if err := s.redis.Set(ctx, user, redis); err != nil {
			log.Fatalln("Failed to set data to redis")
		}
	}
	return nil
}

func (s newsfeedService) GetNewsfeed(
	ctx context.Context,
	userId string,
) (*newsfeedv1.GetNewsfeedResponse, error) {
	data, exist := s.redis.Get(ctx, userId)
	if exist == nil {
		var posts []domain.Post
		if err := json.Unmarshal([]byte(data.(string)), &posts); err != nil {
			log.Println("error when marshal old post")
			return nil, err
		}
		postPointers := make([]*newsfeedv1.NewsfeedPost, len(posts))
		reverseSlice(posts)
		for i, post := range posts {
			postPointers[i] = &newsfeedv1.NewsfeedPost{
				UserId:      post.User_id,
				PostId:      post.Id,
				ContentText: post.Content_text,
				Images:      strings.Split(post.Images, ","),
				Tags:        strings.Split(post.Tags, ","),
				Visible:     post.Visible,
				Date:        post.CreatedAt.String(),
			}
		}
		return &newsfeedv1.GetNewsfeedResponse{
			Posts: postPointers,
		}, nil
	}
	return nil, nil
}

func reverseSlice(posts []domain.Post) {
	n := len(posts)
	for i := 0; i < n/2; i++ {
		posts[i], posts[n-1-i] = posts[n-1-i], posts[i]
	}
}
