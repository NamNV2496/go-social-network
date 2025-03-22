package database

import (
	"time"

	"gorm.io/gorm"
)

type Entity interface {
	TableName() string
	Exist() bool
}

type BaseEntity struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (be *BaseEntity) TableName() string {
	return ""
}

func (be *BaseEntity) Exist() bool {
	return false
}

func (be *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	be.CreatedAt = time.Now()
	be.UpdatedAt = time.Now()
	return nil
}

func (be *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	be.UpdatedAt = time.Now()
	return nil
}
