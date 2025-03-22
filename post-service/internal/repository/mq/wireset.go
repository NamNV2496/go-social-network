package mq

import (
	"github.com/google/wire"

	"github.com/namnv2496/post-service/internal/repository/mq/producer"
)

var MQWireSet = wire.NewSet(
	producer.NewClient,
)
