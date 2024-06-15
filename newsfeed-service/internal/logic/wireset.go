package logic

import "github.com/google/wire"

var LogicWireSet = wire.NewSet(
	NewUserService,
)
