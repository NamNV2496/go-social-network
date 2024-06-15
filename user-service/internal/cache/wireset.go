package cache

import "github.com/google/wire"

var CacheWireSet = wire.NewSet(
	NewRedisClient,
)
