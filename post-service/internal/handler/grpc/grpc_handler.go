package grpc

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/logic"
)

type GrpcHandler struct {
	postv1.UnimplementedPostServiceServer
	userService logic.UserService
}

func NewGrpcHander(
	userService logic.UserService,
) postv1.PostServiceServer {
	return &GrpcHandler{
		userService: userService,
	}
}

func (h GrpcHandler) CreatePost(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {

	return h.userService.Post(ctx, req)
}

func (h GrpcHandler) GetPost(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {
	return nil, nil
}
