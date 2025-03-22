package controller

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
)

func (c Controller) LikeAction(ctx context.Context, req *postv1.LikeRequest) (*postv1.LikeResponse, error) {
	return c.likeService.LikeAction(ctx, req)
}

func (c Controller) Getlike(ctx context.Context, req *postv1.GetLikeRequest) (*postv1.GetLikeResponse, error) {
	return c.likeService.GetLikeStatsByPostIdAndUserId(ctx, req)
}
