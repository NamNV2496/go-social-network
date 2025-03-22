package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNamePost = "post"
)

type Post struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64  `gorm:"column:id;primaryKey;autoIncrement" json:"post_id"`
	Uuid                string `gorm:"column:uuid;type:varchar(36);not null" json:"uuid"`
	User_id             string `gorm:"column:user_id;type:varchar(255);not null" json:"user_id"`
	Content_text        string `gorm:"column:content_text;type:text" json:"content_text"`
	Images              string `gorm:"column:images;type:text" json:"images"`
	Tags                string `gorm:"column:tags;type:text" json:"tags"`
	Visible             bool   `gorm:"column:visible;type:boolean;default:true" json:"visible"`
}

func PostByUserId(userId string) database.QueryOption {
	return func(tx *gorm.DB) {
		if userId != "" {
			tx.Where("post.user_id = ?", userId)
		}
	}
}

func PostByUuid(uuid string) database.QueryOption {
	return func(tx *gorm.DB) {
		if uuid != "" {
			tx.Where("post.uuid = ?", uuid)
		}
	}
}

func PostOrderById() database.QueryOption {
	return func(tx *gorm.DB) {
		tx.Order("id DESC")
	}
}

func (_self Post) TableName() string {
	return TabNamePost
}
