package grpc

import (
	"context"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/logic"
)

type GrpcHandler struct {
	postv1.UnimplementedPostServiceServer
	postService    logic.PostService
	commentService logic.CommentService
	likeService    logic.LikeService
}

func NewGrpcHander(
	postService logic.PostService,
	commentService logic.CommentService,
	likeService logic.LikeService,
) postv1.PostServiceServer {
	return &GrpcHandler{
		postService:    postService,
		commentService: commentService,
		likeService:    likeService,
	}
}

func (h GrpcHandler) CreatePost(
	ctx context.Context,
	req *postv1.CreatePostRequest,
) (*postv1.CreatePostResponse, error) {

	return h.postService.Post(ctx, req)
}

func (h GrpcHandler) GetPost(
	ctx context.Context,
	req *postv1.GetPostRequest,
) (*postv1.GetPostResponse, error) {
	return h.postService.GetPosts(ctx, req)
}

func (h GrpcHandler) CreateComment(
	ctx context.Context,
	req *postv1.CreateCommentRequest,
) (*postv1.CreateCommentResponse, error) {

	return h.commentService.Comment(ctx, req)
}

func (h GrpcHandler) GetComment(
	ctx context.Context,
	req *postv1.GetCommentRequest,
) (*postv1.GetCommentResponse, error) {

	return h.commentService.GetComment(ctx, req)
}

func (h GrpcHandler) LikeAction(
	ctx context.Context,
	req *postv1.LikeRequest,
) (*postv1.LikeResponse, error) {

	return h.likeService.LikeAction(ctx, req)
}

func (h GrpcHandler) Getlike(
	ctx context.Context,
	req *postv1.GetLikeRequest,
) (*postv1.GetLikeResponse, error) {

	return h.likeService.Getlike(ctx, req)
}
