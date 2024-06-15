package mq

import (
	"time"
)

var (
	REDIS_NEWSFEED_KEY = "newsfeed"
	RedisTTL           = 86400 * time.Second
)
