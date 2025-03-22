package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNameComment = "comment"
)

type Comment struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64  `gorm:"column:id;primaryKey" json:"id"`
	PostId              int64  `gorm:"column:post_id;not null" json:"post_id"`
	UserId              string `gorm:"column:user_id;not null" json:"user_id"`
	CommentText         string `gorm:"column:comment_text;type:text;not null" json:"comment_text"`
	CommentLevel        int64  `gorm:"column:comment_level;not null" json:"comment_level"`
	CommentParent       int64  `gorm:"column:comment_parent" json:"comment_parent"`
	Images              string `gorm:"column:images;type:text" json:"images"`
	Tags                string `gorm:"column:tags;type:text" json:"tags"`
}

func (_self Comment) TableName() string {
	return TabNameComment
}

func CommentByPostId(postId int64) database.QueryOption {
	return func(tx *gorm.DB) {
		if postId > 0 {
			tx.Where("comment.post_id = ?", postId)
		}
	}
}
