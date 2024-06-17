package grpc

import (
	"context"

	newsfeedv1 "github.com/namnv2496/newsfeed-service/internal/handler/generated/newsfeed_core/v1"
	"github.com/namnv2496/newsfeed-service/internal/logic"
)

type GrpcHandler struct {
	newsfeedv1.UnimplementedNewsfeedServiceServer
	newsfeedService logic.NewsfeedService
}

func NewGrpcHander(
	newsfeedService logic.NewsfeedService,
) newsfeedv1.NewsfeedServiceServer {
	return &GrpcHandler{
		newsfeedService: newsfeedService,
	}
}

func (h GrpcHandler) GetNewsfeed(
	ctx context.Context,
	req *newsfeedv1.GetNewsfeedRequest,
) (*newsfeedv1.GetNewsfeedResponse, error) {

	return h.newsfeedService.GetNewsfeed(ctx, req.UserId)
}
