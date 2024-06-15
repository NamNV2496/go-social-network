package grpc

import (
	postv1 "github.com/namnv2496/newsfeed-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/newsfeed-service/internal/logic"
)

type GrpcHandler struct {
	postv1.UnimplementedPostServiceServer
	newsfeedService logic.NewsfeedService
}

func NewGrpcHander(
	newsfeedService logic.NewsfeedService,
) postv1.PostServiceServer {
	return &GrpcHandler{
		newsfeedService: newsfeedService,
	}
}
