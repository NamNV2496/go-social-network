package controller

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
)

func (c Controller) CreateComment(
	ctx context.Context,
	req *postv1.CreateCommentRequest,
) (*postv1.CreateCommentResponse, error) {
	return c.commentService.AddComment(ctx, req)
}

func (c Controller) GetComment(
	ctx context.Context,
	req *postv1.GetCommentRequest,
) (*postv1.GetCommentResponse, error) {
	return c.commentService.GetComment(ctx, req)
}
