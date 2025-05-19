package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type ILikeCountRepository interface {
	UpsertTotalLikeByPostIds(ctx context.Context, req domain.LikeCount, action bool, opts ...database.QueryOption) ([]*domain.LikeCount, error)
	GetTotalLikeByPostIds(ctx context.Context, opts ...database.QueryOption) ([]*domain.LikeCount, error)
}

type LikeCountRepository struct {
	database.ICRUDBase[*domain.LikeCount]
}

func NewLikeCountRepository(
	db *gorm.DB,
	config configs.Config,
) (ILikeCountRepository, error) {
	dbConfig := config.Database
	likeCountRepository := database.NewCRUDBase[*domain.LikeCount](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.LikeCount{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &LikeCountRepository{
		ICRUDBase: likeCountRepository,
	}, nil
}

func (_self *LikeCountRepository) UpsertTotalLikeByPostIds(ctx context.Context, req domain.LikeCount, action bool, opts ...database.QueryOption) ([]*domain.LikeCount, error) {
	likeCounts, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if len(likeCounts) == 0 {
		if err := _self.Create(ctx, &req); err != nil {
			return nil, err
		}
		return []*domain.LikeCount{&req}, nil
	}
	if action {
		likeCounts[0].TotalLike += 1
	} else {
		likeCounts[0].TotalLike -= 1
	}
	if err := _self.Update(ctx, likeCounts[0], opts...); err != nil {
		return nil, err
	}
	return likeCounts, nil
}

func (_self *LikeCountRepository) GetTotalLikeByPostIds(ctx context.Context, opts ...database.QueryOption) ([]*domain.LikeCount, error) {
	likeCounts, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return likeCounts, nil
}
