package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type IPostRepository interface {
	AddPost(ctx context.Context, req domain.Post, opts ...database.QueryOption) error
	GetPosts(ctx context.Context, opts ...database.QueryOption) ([]*domain.Post, error)
}

type PostRepository struct {
	database.ICRUDBase[*domain.Post]
}

func NewPostRepository(
	db *gorm.DB,
	config configs.Config,
) (IPostRepository, error) {
	dbConfig := config.Database
	postRepository := database.NewCRUDBase[*domain.Post](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.Post{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &PostRepository{
		ICRUDBase: postRepository,
	}, nil
}

func (_self PostRepository) AddPost(ctx context.Context, req domain.Post, opts ...database.QueryOption) error {
	if err := _self.Create(ctx, &req); err != nil {
		return err
	}
	return nil
}

func (_self PostRepository) GetPosts(ctx context.Context, opts ...database.QueryOption) ([]*domain.Post, error) {
	var posts []*domain.Post
	posts, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
