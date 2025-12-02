package controller

import (
	"github.com/namnv2496/post-service/configs"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg/logger"
	"github.com/namnv2496/post-service/internal/service"
	"github.com/namnv2496/post-service/internal/validator"
)

type Controller struct {
	postv1.UnimplementedPostServiceServer
	logger         *logger.Logger
	likeService    service.ILikeService
	postService    service.IPostService
	commentService service.ICommentService
	validatorSvc   validator.IPostValidator
}

func NewController(
	cfg configs.Config,
	likeService service.ILikeService,
	postService service.IPostService,
	commentService service.ICommentService,
	validatorSvc validator.IPostValidator,
) postv1.PostServiceServer {
	return &Controller{
		logger:         logger.NewLogger("post-controller", "post"),
		likeService:    likeService,
		postService:    postService,
		commentService: commentService,
		validatorSvc:   validatorSvc,
	}
}
