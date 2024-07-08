package es

import "github.com/google/wire"

var ESWireSet = wire.NewSet(
	NewElasticSearch,
)
