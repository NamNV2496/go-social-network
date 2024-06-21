package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/namnv2496/newsfeed-service/internal/domain"
	"github.com/namnv2496/newsfeed-service/internal/handler/grpc"
	"github.com/namnv2496/newsfeed-service/internal/logic"
	"github.com/namnv2496/newsfeed-service/internal/mq"
	"github.com/namnv2496/newsfeed-service/internal/mq/consumer"
)

type ConsumerHandler interface {
	StartConsumerUp(ctx context.Context) error
}

type consumerHandler struct {
	consumer        consumer.Consumer
	newsfeedService logic.NewsfeedService
	userClient      grpc.ProductGRPCClient
}

func NewKafkaHandler(
	consumer consumer.Consumer,
	newsfeedService logic.NewsfeedService,
	userClient grpc.ProductGRPCClient,
) ConsumerHandler {
	return &consumerHandler{
		consumer:        consumer,
		newsfeedService: newsfeedService,
		userClient:      userClient,
	}
}

func (c consumerHandler) StartConsumerUp(ctx context.Context) error {
	fmt.Println("Listening from topic: ", mq.TOPIC_POST_CONTENT)
	c.consumer.RegisterHandler(
		mq.TOPIC_POST_CONTENT,
		func(ctx context.Context, queueName string, payload []byte) error {
			fmt.Println("listen from queue: " + queueName + ". Data: " + string(payload))
			var newPost domain.Post
			if err := json.Unmarshal([]byte(payload), &newPost); err != nil {
				fmt.Println("error when marshal new post")
				return err
			}
			// update newsfeed for current user
			if err := c.newsfeedService.UpdateNewsfeed(ctx, []string{newPost.User_id}, newPost); err != nil {
				return err
			}
			// update newsfeed for taged users
			var tags []string
			if newPost.Tags != "" {
				tags = strings.Split(newPost.Tags, ",")
				fmt.Println("Update for taged users: ", tags)
				if err := c.newsfeedService.UpdateNewsfeed(ctx, tags, newPost); err != nil {
					return err
				}
			}
			var tagMap = make(map[string]bool, 0)
			for _, tag := range tags {
				tag = strings.Trim(tag, " ")
				tagMap[tag] = true
			}
			// // update newsfeed for followings
			followings, err := c.userClient.GetFollowing(ctx, newPost.User_id)
			var followingUpdate []string
			if err != nil {
				return err
			} else {
				for _, follower := range followings {
					if !tagMap[follower] {
						follower = strings.Trim(follower, " ")
						followingUpdate = append(followingUpdate, follower)
					}
				}
				fmt.Println("Update for followings: ", followingUpdate)
				if err := c.newsfeedService.UpdateNewsfeed(ctx, followingUpdate, newPost); err != nil {
					return err
				}
			}

			return nil
		},
	)
	return nil
}
