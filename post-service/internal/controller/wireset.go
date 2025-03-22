package controller

import (
	"github.com/google/wire"
)

var ControllerWireSet = wire.NewSet(
	NewController,
)
