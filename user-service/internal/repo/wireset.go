package repo

import "github.com/google/wire"

var RepoWireSet = wire.NewSet(
	NewUserService,
	NewUserUserService,
)
