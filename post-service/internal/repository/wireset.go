package repository

import (
	"github.com/google/wire"
	"github.com/namnv2496/post-service/internal/repository/database"
	"github.com/namnv2496/post-service/internal/repository/mq"
)

var RepositoryWireSet = wire.NewSet(
	database.DatabaseWireSet,
	mq.MQWireSet,
	NewLikeRepository,
	NewLikeCountRepository,
	NewCommentRepository,
	NewPostRepository,
	NewTransaction,
)
