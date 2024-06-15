package mq

import (
	"github.com/google/wire"

	"github.com/namnv2496/newsfeed-service/internal/mq/consumer"
)

var MQWireSet = wire.NewSet(
	consumer.NewConsumer,
)
