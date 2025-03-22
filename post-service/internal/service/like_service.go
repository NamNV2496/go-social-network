package service

import (
	"context"

	"github.com/namnv2496/post-service/internal/domain"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/repository/mq/producer"
)

type ILikeService interface {
	LikeAction(context.Context, *postv1.LikeRequest) (*postv1.LikeResponse, error)
	GetLikeStatsByPostIdAndUserId(context.Context, *postv1.GetLikeRequest) (*postv1.GetLikeResponse, error)
}

type likeService struct {
	likeRepository        repository.ILikeRepository
	likeCountRepository   repository.ILikeCountRepository
	transactionRepository repository.ITransaction
}

func NewLikeService(
	likeRepository repository.ILikeRepository,
	likeCountRepository repository.ILikeCountRepository,
	transactionRepository repository.ITransaction,
	kafkaClient producer.Client,
) ILikeService {
	return &likeService{
		likeRepository:        likeRepository,
		likeCountRepository:   likeCountRepository,
		transactionRepository: transactionRepository,
	}
}

func (l *likeService) LikeAction(
	ctx context.Context,
	req *postv1.LikeRequest,
) (*postv1.LikeResponse, error) {
	if req.Like.UserId == "" || req.Like.PostId == 0 {
		return nil, nil
	}
	var totalLike int64
	var likeAction bool
	if err := l.transactionRepository.WithTransaction(ctx, func(ctx context.Context) error {
		// get Like status of user
		likes, err := l.likeRepository.GetlikeByPostIdAndUserId(ctx, domain.LikeByPostId(int64(req.Like.PostId)), domain.LikeByUserId(req.Like.UserId))
		if err != nil {
			return err
		}
		entity := &domain.Like{
			PostId: int64(req.Like.PostId),
			UserId: req.Like.UserId,
			Like:   req.Like.Action == postv1.Like_Enum_LIKE,
		}
		if len(likes) > 0 && likes[0] != nil {
			entity = likes[0]
			entity.Like = !entity.Like
		}
		likeAction = entity.Like
		if err := l.likeRepository.UpsertLikeByPostIdAndUserId(ctx, entity, domain.LikeByPostId(int64(req.Like.PostId)), domain.LikeByUserId(req.Like.UserId)); err != nil {
			return err
		}
		// update total like count
		likeCounts, err := l.likeCountRepository.UpsertTotalLikeByPostIds(ctx, domain.LikeCount{
			PostId:    int64(req.Like.PostId),
			TotalLike: 1,
		}, likeAction, domain.LikeCountByPostId(int64(req.Like.PostId)))
		if err != nil {
			return err
		}
		totalLike = likeCounts[0].TotalLike
		return nil
	}); err != nil {
		return nil, err
	}
	return &postv1.LikeResponse{
		Response: &postv1.LikePostResponse{
			PostId:    req.Like.PostId,
			Like:      likeAction,
			TotalLike: uint64(totalLike),
		},
	}, nil
}

func (l *likeService) GetLikeStatsByPostIdAndUserId(
	ctx context.Context,
	req *postv1.GetLikeRequest,
) (*postv1.GetLikeResponse, error) {
	if len(req.PostId) == 0 {
		return nil, nil
	}
	var resp []*postv1.LikePostResponse
	if err := l.transactionRepository.WithTransaction(ctx, func(ctx context.Context) error {
		// Get total like by postId
		likeCounts, err := l.likeCountRepository.GetTotalLikeByPostIds(ctx, domain.LikeCountByPostIds(pkg.ConvertListUintToInt64(req.PostId)))
		if err != nil {
			return err
		}
		// Get Like status of user
		likes, err := l.likeRepository.GetlikeByPostIdAndUserId(ctx, domain.LikeByPostId(int64(req.PostId[0])), domain.LikeByUserId(req.UserId))
		if err != nil {
			return err
		}
		userLikeStatus := make(map[int64]bool)
		for _, like := range likes {
			userLikeStatus[like.PostId] = like.Like
		}
		for _, like := range likeCounts {
			resp = append(resp, &postv1.LikePostResponse{
				PostId:    uint64(like.PostId),
				Like:      userLikeStatus[like.PostId],
				TotalLike: uint64(like.TotalLike),
			})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &postv1.GetLikeResponse{
		Response: resp,
	}, nil
}
