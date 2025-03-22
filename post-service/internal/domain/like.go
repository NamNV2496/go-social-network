package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNameLike = "like"
)

type Like struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64  `gorm:"column:id;type:bigint" json:"id"`
	PostId              int64  `gorm:"column:post_id;type:bigint" json:"post_id"`
	UserId              string `gorm:"column:user_id;type:text" json:"user_id"`
	Like                bool   `gorm:"column:like;type:bool" json:"like"`
}

func (_self Like) TableName() string {
	return TabNameLike
}

func LikeByPostId(postId int64) database.QueryOption {
	return func(tx *gorm.DB) {
		if postId > 0 {
			tx.Where("like.post_id = ?", postId)
		}
	}
}

func LikeByUserId(userId string) database.QueryOption {
	return func(tx *gorm.DB) {
		if userId != "" {
			tx.Where("like.user_id = ?", userId)
		}
	}
}
