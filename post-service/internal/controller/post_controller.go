package controller

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg/metric"
	// "github.com/opentracing/opentracing-go"
)

func (c *Controller) CreatePost(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	// defer span.Finish()
	metric.MetricNewPostCnt("CreatePost")
	return c.postService.AddPost(ctx, req)
}

func (c *Controller) GetPost(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {
	// span, ctx := opentracing.StartSpanFromContext(ctx, "user.Create")
	// defer span.Finish()
	metric.MetricGetPostCnt("GetPost")
	return c.postService.GetPosts(ctx, req)
}
