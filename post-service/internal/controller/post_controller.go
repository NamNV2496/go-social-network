package controller

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
)

func (c Controller) CreatePost(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {
	return c.postService.AddPost(ctx, req)
}

func (c Controller) GetPost(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {
	return c.postService.GetPosts(ctx, req)
}
