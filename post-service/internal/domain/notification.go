package domain

import (
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

var (
	TabNameNotification = "notification"
)

type Notification struct {
	database.BaseEntity `gorm:"embedded"`
	Id                  int64  `gorm:"column:id;primaryKey" json:"id"`
	Title               string `gorm:"column:title;not null" json:"title"`
	Description         string `gorm:"column:description" json:"description"`
	Template            string `gorm:"column:template" json:"template"`
	Image               string `gorm:"column:image" json:"image"`
	Application         string `gorm:"column:application" json:"application"`
	Visible             bool   `gorm:"column:visible" json:"visible"`
	Link                string `gorm:"column:link" json:"link"`
}

func (_self *Notification) TableName() string {
	return TabNameNotification
}

func NotificationByIdAndApplications(id int64, applications []string) database.QueryOption {
	return func(tx *gorm.DB) {
		if id > 0 && len(applications) > 0 {
			tx.Where("notification.id = ? AND notification.application in ? AND notification.visible = true", id, applications)
		}
	}
}

func NotificationByIdAndApplication(id int64, application string) database.QueryOption {
	return func(tx *gorm.DB) {
		if id > 0 && len(application) > 0 {
			tx.Where("notification.id = ? AND notification.application = ? AND notification.visible = true", id, application)
		}
	}
}

func NotificationByApplications(applications []string) database.QueryOption {
	return func(tx *gorm.DB) {
		if len(applications) > 0 {
			tx.Where("notification.application in ? AND notification.visible = true", applications)
		}
	}
}

func WithPagination(page, limit int) database.QueryOption {
	return func(tx *gorm.DB) {
		if limit > 0 {
			tx.Offset(page * limit).Limit(limit)
		}
	}
}
