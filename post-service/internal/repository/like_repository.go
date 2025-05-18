package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type ILikeRepository interface {
	GetlikeByPostIdAndUserId(ctx context.Context, opts ...database.QueryOption) ([]*domain.Like, error)
	UpsertLikeByPostIdAndUserId(ctx context.Context, req *domain.Like, opts ...database.QueryOption) error
}

type LikeRepository struct {
	database.ICRUDBase[*domain.Like]
}

func NewLikeRepository(
	db *gorm.DB,
	dbConfig configs.Database,
) (ILikeRepository, error) {
	likeRepository := database.NewCRUDBase[*domain.Like](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.Like{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &LikeRepository{
		ICRUDBase: likeRepository,
	}, nil
}

func (_self *LikeRepository) GetlikeByPostIdAndUserId(ctx context.Context, opts ...database.QueryOption) ([]*domain.Like, error) {
	var likes []*domain.Like
	likes, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func (_self *LikeRepository) UpsertLikeByPostIdAndUserId(ctx context.Context, req *domain.Like, opts ...database.QueryOption) error {
	if req.Id > 0 {
		if err := _self.Update(ctx, req, opts...); err != nil {
			return err
		}
	} else {
		if err := _self.Create(ctx, req); err != nil {
			return err
		}
	}
	return nil
}
