package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/internal/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment, opts ...database.QueryOption) error
	GetComment(ctx context.Context, opts ...database.QueryOption) ([]*domain.Comment, error)
}

type CommentRepository struct {
	database.ICRUDBase[*domain.Comment]
}

func NewCommentRepository(
	db *gorm.DB,
	dbConfig configs.Database,
) (ICommentRepository, error) {
	commentRepository := database.NewCRUDBase[*domain.Comment](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.Comment{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &CommentRepository{
		ICRUDBase: commentRepository,
	}, nil
}

func (_self *CommentRepository) AddComment(ctx context.Context, comment domain.Comment, opts ...database.QueryOption) error {
	if err := _self.Create(ctx, &comment); err != nil {
		return err
	}
	return nil
}

func (_self *CommentRepository) GetComment(ctx context.Context, opts ...database.QueryOption) ([]*domain.Comment, error) {
	comments, err := _self.Find(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
