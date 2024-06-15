package consumers

import (
	"context"
	"fmt"

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
}

func NewKafkaHandler(
	consumer consumer.Consumer,
	newsfeedService logic.NewsfeedService,
) ConsumerHandler {
	return &consumerHandler{
		consumer:        consumer,
		newsfeedService: newsfeedService,
	}
}

func (c consumerHandler) StartConsumerUp(ctx context.Context) error {
	fmt.Println("Add consumer for topic: ", mq.TOPIC_POST_CONTENT)
	c.consumer.RegisterHandler(
		mq.TOPIC_POST_CONTENT,
		func(ctx context.Context, queueName string, payload []byte) error {
			fmt.Println("listen from queue: " + queueName + ". Data: " + string(payload))
			c.newsfeedService.UpdateNewsfeed(ctx, payload)
			return nil
		},
	)
	return nil
}
