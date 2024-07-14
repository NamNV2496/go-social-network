package logic

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/namnv2496/post-service/internal/domain"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/mq/producer"
)

type LikeService interface {
	LikeAction(context.Context, *postv1.LikeRequest) (*postv1.LikeResponse, error)
	Getlike(context.Context, *postv1.GetLikeRequest) (*postv1.GetLikeResponse, error)
}

type likeService struct {
	db *goqu.Database
}

func NewLikeService(
	db *goqu.Database,
	kafkaClient producer.Client,
) LikeService {
	return &likeService{
		db: db,
	}
}

func (l likeService) Getlike(
	ctx context.Context,
	req *postv1.GetLikeRequest,
) (*postv1.GetLikeResponse, error) {

	// get Like status of user
	postUserLikes := make(map[uint64]bool)
	if req.UserId != "" {
		query := l.db.
			From(domain.TabNameLike).
			Where(
				goqu.C(domain.TabColPostId).In(req.PostId),
				goqu.C(domain.TabColUserId).Eq(req.UserId),
			)
		var userLikes []domain.Like
		err := query.Executor().ScanStructsContext(ctx, &userLikes)
		if err != nil {
			return nil, err
		}
		for _, userLike := range userLikes {
			postUserLikes[userLike.PostId] = userLike.Like
		}
	}

	// get total like of post
	query := l.db.
		From(domain.TabNameLikeCount).
		Where(
			goqu.C(domain.TabColPostId).In(req.PostId),
		)
	var likes []domain.LikeCount
	err := query.Executor().ScanStructsContext(ctx, &likes)
	if err != nil {
		return nil, err
	}
	var response []*postv1.LikePostResponse

	for _, like := range likes {
		var element postv1.LikePostResponse
		element.PostId = like.PostId
		element.TotalLike = uint64(like.TotalLike)
		val, ok := postUserLikes[like.PostId]
		if ok {
			element.Like = val
		}
		response = append(response, &element)
	}

	return &postv1.GetLikeResponse{
		Response: response,
	}, nil
}

func (l likeService) LikeAction(
	ctx context.Context,
	req *postv1.LikeRequest,
) (*postv1.LikeResponse, error) {

	exist, status, err := l.checkExist(ctx, req)
	if err != nil {
		return nil, err
	}
	if (status == false && req.Like.Action == 0) || (status == true && req.Like.Action == 1) {
		return nil, errors.New("invalid input")
	}
	action := 1
	if req.Like.Action == 0 {
		// unlike
		action = -1
	}
	if exist {
		log.Println("Update like")
		if err := l.updateLike(ctx, req); err != nil {
			return nil, err
		}
	} else {
		log.Println("Add new like")
		if err := l.newLike(ctx, req); err != nil {
			return nil, err
		}
	}

	// update totalLike
	query := l.db.
		From(domain.TabNameLikeCount).
		Where(
			goqu.C(domain.TabColPostId).Eq(req.Like.PostId),
		)
	var likeCount domain.LikeCount
	exist, err = query.Executor().ScanStructContext(ctx, &likeCount)
	if err != nil {
		return nil, err
	}
	if !exist {
		l.newLikeCount(ctx, req.Like.PostId)
	} else {
		l.updateLikeCount(ctx, req.Like.PostId, likeCount.TotalLike, action)
	}

	return &postv1.LikeResponse{
		Response: &postv1.LikePostResponse{
			PostId:    req.Like.PostId,
			Like:      req.Like.Action == 0,
			TotalLike: uint64(likeCount.TotalLike + action),
		},
	}, nil
}

func (l likeService) newLikeCount(
	ctx context.Context,
	postId uint64,
) error {
	like := domain.LikeCount{
		PostId:    postId,
		TotalLike: 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	query := l.db.
		Insert(domain.TabNameLikeCount).
		Rows(like)
	_, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l likeService) updateLikeCount(
	ctx context.Context,
	postId uint64,
	totalLike int,
	action int,
) error {
	updateQuery := l.db.
		Update(domain.TabNameLikeCount).
		Set(goqu.Record{
			domain.TabColTotalLike: totalLike + action,
			domain.TabColUpdatedAt: time.Now(),
		}).
		Where(
			goqu.C(domain.TabColPostId).Eq(postId),
		)

	_, err := updateQuery.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l likeService) newLike(
	ctx context.Context,
	req *postv1.LikeRequest,
) error {
	like := domain.Like{
		PostId:    req.Like.PostId,
		UserId:    req.Like.UserId,
		Like:      req.Like.Action == 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	query := l.db.
		Insert(domain.TabNameLike).
		Rows(like)
	_, err := query.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l likeService) checkExist(
	ctx context.Context,
	req *postv1.LikeRequest,
) (bool, bool, error) {
	query := l.db.
		From(domain.TabNameLike).
		Where(
			goqu.C(domain.TabColPostId).Eq(req.Like.PostId),
			goqu.C(domain.TabColUserId).Eq(req.Like.UserId),
		)
	var like domain.Like
	exist, err := query.Executor().ScanStructContext(ctx, &like)
	if err != nil {
		return false, false, err
	}
	return exist, like.Like, nil
}

func (l likeService) updateLike(
	ctx context.Context,
	req *postv1.LikeRequest,
) error {

	updateQuery := l.db.
		Update(domain.TabNameLike).
		Set(goqu.Record{
			domain.TabColLike:      req.Like.Action == 1,
			domain.TabColUpdatedAt: time.Now(),
		}).
		Where(
			goqu.C(domain.TabColPostId).Eq(req.Like.PostId),
			goqu.C(domain.TabColUserId).Eq(req.Like.UserId),
		)

	// log.Println(updateQuery.ToSQL())
	_, err := updateQuery.Executor().ExecContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
