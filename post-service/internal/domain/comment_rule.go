package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNameCommentRule = "comment_rule"
)

type CommentRule struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64  `gorm:"column:id;primaryKey" json:"id"`
	CommentText         string `gorm:"column:comment_text;not null" json:"comment_text"`
	Application         string `gorm:"column:application;primaryKey" json:"application"`
	Visible             bool   `gorm:"column:visible;primaryKey" json:"visible"`
}

func (_self *CommentRule) TableName() string {
	return TabNameCommentRule
}

func CommentRuleByIdAndApplications(id int64, applications []string) database.QueryOption {
	return func(tx *gorm.DB) {
		if id > 0 && len(applications) > 0 {
			tx.Where("comment_rule.id = ? AND comment_rule.application in ? AND comment_rule.visible = true", id, applications)
		}
	}
}

func CommentRuleByIdAndApplication(id int64, application string) database.QueryOption {
	return func(tx *gorm.DB) {
		if id > 0 && len(application) > 0 {
			tx.Where("comment_rule.id = ? AND comment_rule.application = ? AND comment_rule.visible = true", id, application)
		}
	}
}

func CommentRuleByApplications(applications []string) database.QueryOption {
	return func(tx *gorm.DB) {
		if len(applications) > 0 {
			tx.Where("comment_rule.application in ? AND comment_rule.visible = true", applications)
		}
	}
}

func PaginationCommentRule(page, limit int) database.QueryOption {
	return func(tx *gorm.DB) {
		if limit > 0 {
			tx.Offset(page * limit).Limit(limit)
		}
	}
}
