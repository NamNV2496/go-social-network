package configs

import "github.com/google/wire"

var ConfigWireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "Kafka"),
	wire.FieldsOf(new(Config), "GRPC"),
	wire.FieldsOf(new(Config), "Database"),
	wire.FieldsOf(new(Config), "Redis"),
)
