package repository

import (
	"context"

	"gorm.io/gorm"
)

type ITransaction interface {
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

type tracsaction struct {
	db *gorm.DB
}

func NewTransaction(
	db *gorm.DB,
) ITransaction {
	return &tracsaction{
		db: db,
	}
}

type transactionKey struct {
}

func (t *tracsaction) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	newCtx := context.WithValue(ctx, transactionKey{}, t.db)
	db := t.db.Begin()
	if err := fn(newCtx); err != nil {
		db.Rollback()
		return err
	}
	if err := db.Commit(); err != nil {
		return err.Error
	}
	return nil
}
