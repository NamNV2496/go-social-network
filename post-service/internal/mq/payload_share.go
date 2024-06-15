package mq

import (
	"time"
)

var (
	REDIS_KEY_ORDER = "kitchen"
	RedisTTL        = 86400 * time.Second
)
