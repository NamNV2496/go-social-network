package database

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICRUDBase[T Entity] interface {
	Database(ctx context.Context) *gorm.DB

	Create(context.Context, T) error
	Find(context.Context, ...QueryOption) ([]T, error)
	Update(ctx context.Context, entity T, queryOpts ...QueryOption) error
	Updates(ctx context.Context, entities []T, queryOpts ...QueryOption) error
	Upsert(context.Context, []T, ...QueryOption) error
	DeleteById(context.Context, int64) error

	EnableDebug(debug bool)
	Debug()
}

type CRUDBase[T Entity] struct {
	database *gorm.DB
	debug    bool
}

type QueryOption func(tx *gorm.DB)

type transactionKey struct {
}

func NewCRUDBase[T Entity](db *gorm.DB) ICRUDBase[T] {
	return &CRUDBase[T]{
		database: db,
		debug:    false,
	}
}

func GetDatabase(ctx context.Context) *gorm.DB {
	txValue := ctx.Value(transactionKey{})
	switch txSession := txValue.(type) {
	case *gorm.DB:
		return txSession
	}
	return nil
}

func SetDatabase(ctx context.Context, database *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey{}, database)
}

func (_self *CRUDBase[T]) Database(ctx context.Context) *gorm.DB {
	tx := GetDatabase(ctx)
	if tx != nil {
		return tx
	}
	SetDatabase(ctx, _self.database)
	return _self.database
}

// func (_self *CRUDBase[T]) beginTransaction(ctx context.Context) (context.Context, *gorm.DB) {
// 	tx := _self.database.Begin()
// 	newCtx := SetDatabase(ctx, tx)
// 	return newCtx, tx
// }

// func (_self *CRUDBase[T]) commitTransaction(ctx context.Context) error {
// 	tx := _self.Database(ctx)
// 	if tx != nil {
// 		return tx.Commit().Error
// 	}
// 	return errors.New("operation is not in any tracsaction")
// }

// func (_self *CRUDBase[T]) rollbackTransaction(ctx context.Context) error {
// 	tx := _self.Database(ctx)
// 	if tx != nil {
// 		return tx.Rollback().Error
// 	}
// 	return errors.New("operation is not in any tracsaction")
// }

func (_self *CRUDBase[T]) Create(ctx context.Context, _entity T) error {
	tx := _self.Database(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	return tx.WithContext(ctx).Create(_entity).Error
}

func (_self *CRUDBase[T]) Find(ctx context.Context, queryOpts ...QueryOption) ([]T, error) {
	tx := _self.Database(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	var _entity T
	entities := make([]T, 0)
	tx = tx.WithContext(ctx).Model(&_entity)
	setOptToQuery(tx, queryOpts...)
	err := tx.Find(&entities).Error
	return entities, err
}

// upsert
func (_self *CRUDBase[T]) Update(ctx context.Context, entity T, queryOpts ...QueryOption) error {
	tx := _self.Database(ctx).WithContext(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	setOptToQuery(tx, queryOpts...)
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&entity).Error
}

func (_self *CRUDBase[T]) Updates(ctx context.Context, entities []T, queryOpts ...QueryOption) error {
	if len(entities) == 0 {
		return nil
	}
	tx := _self.Database(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	var _entities T
	tx = tx.WithContext(ctx).Model(&_entities)
	setOptToQuery(tx, queryOpts...)
	return tx.Save(&entities).Error
}

func (_self *CRUDBase[T]) Upsert(ctx context.Context, entities []T, queryOpts ...QueryOption) error {
	if len(entities) == 0 {
		return nil
	}
	tx := _self.Database(ctx).WithContext(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	setOptToQuery(tx, queryOpts...)
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(entities).Error
}

func (_self *CRUDBase[T]) DeleteById(ctx context.Context, id int64) error {
	tx := _self.Database(ctx)
	tx = tx.WithContext(ctx)
	if _self.debug {
		tx = tx.Debug()
	}
	var _entity T
	return tx.Model(&_entity).Where("id=?", id).Delete(&_entity).Error
}

func (_self *CRUDBase[T]) EnableDebug(debug bool) {
	_self.debug = debug
}

func (_self *CRUDBase[T]) Debug() {
	if _self.debug {
		_self.database.Debug()
	}
}

func setOptToQuery(db *gorm.DB, queryOpts ...QueryOption) {
	if len(queryOpts) == 0 {
		return
	}
	var queryOpt QueryOption
	for _, queryOption := range queryOpts {
		queryOpt = queryOption
		queryOpt(db)
	}
}
