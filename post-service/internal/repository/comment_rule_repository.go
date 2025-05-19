package repository

import (
	"context"
	"errors"

	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type ICommentRuleRepository interface {
	AddCommentRule(ctx context.Context, commentRule domain.CommentRule) (int64, error)
	GetCommentRuleById(ctx context.Context, id int64, applications []string) (*domain.CommentRule, error)
	GetCommentRules(ctx context.Context, applications []string, page, limit int32) ([]*domain.CommentRule, error)
	UpdateCommentRule(ctx context.Context, comment domain.CommentRule) error
	CountCommentRules(ctx context.Context, applications []string) (int64, error)
}

type CommentRuleRepository struct {
	db database.ICRUDBase[*domain.CommentRule]
}

func NewCommentRuleRepository(
	db *gorm.DB,
	config configs.Config,
) (*CommentRuleRepository, error) {
	dbConfig := config.Database
	commentRuleRepository := database.NewCRUDBase[*domain.CommentRule](db)
	if dbConfig.AutoMigrate {
		if err := db.Debug().AutoMigrate(&domain.CommentRule{}); err != nil {
			return nil, errors.New("cannot create like repository")
		}
	}
	return &CommentRuleRepository{
		db: commentRuleRepository,
	}, nil
}

var _ ICommentRuleRepository = &CommentRuleRepository{}

func (_self *CommentRuleRepository) AddCommentRule(ctx context.Context, commentRule domain.CommentRule) (int64, error) {
	if err := _self.db.Create(ctx, &commentRule); err != nil {
		return 0, err
	}
	return commentRule.Id, nil
}

func (_self *CommentRuleRepository) GetCommentRuleById(ctx context.Context, id int64, applications []string) (*domain.CommentRule, error) {
	result, err := _self.db.Find(ctx, domain.CommentRuleByIdAndApplications(id, applications))
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result[0], nil
}

func (_self *CommentRuleRepository) GetCommentRules(ctx context.Context, applications []string, page, limit int32) ([]*domain.CommentRule, error) {
	return _self.db.Find(ctx,
		domain.CommentRuleByApplications(applications),
		domain.PaginationCommentRule(int(page), int(limit)),
	)
}

func (_self *CommentRuleRepository) UpdateCommentRule(ctx context.Context, commentRule domain.CommentRule) error {
	if commentRule.Id == 0 {
		return errors.New("invalid comment rule id")
	}
	if err := _self.db.Update(ctx,
		&commentRule,
		domain.CommentRuleByIdAndApplication(commentRule.Id, commentRule.Application),
	); err != nil {
		return err
	}
	return nil
}

func (_self *CommentRuleRepository) CountCommentRules(ctx context.Context, applications []string) (int64, error) {
	return _self.db.Count(ctx, domain.CommentRuleByApplications(applications))
}
