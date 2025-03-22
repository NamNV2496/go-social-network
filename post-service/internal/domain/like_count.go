package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNameLikeCount = "like_count"
)

type LikeCount struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64 `gorm:"id" json:"id"`
	PostId              int64 `gorm:"post_id" json:"post_id"`
	TotalLike           int64 `gorm:"total_like" json:"total_like"`
}

func (_self LikeCount) TableName() string {
	return TabNameLikeCount
}

func LikeCountByPostId(postId int64) database.QueryOption {
	return func(tx *gorm.DB) {
		if postId > 0 {
			tx.Where("like_count.post_id = ?", postId)
		}
	}
}

func LikeCountByPostIds(postIds []int64) database.QueryOption {
	return func(tx *gorm.DB) {
		if len(postIds) > 0 {
			tx.Where("like_count.post_id in ?", postIds)
		}
	}
}
