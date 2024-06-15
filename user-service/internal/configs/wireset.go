package configs

import "github.com/google/wire"

var ConfigWireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Auth"),
	wire.FieldsOf(new(Config), "Redis"),
	wire.FieldsOf(new(Config), "GRPC"),
)
