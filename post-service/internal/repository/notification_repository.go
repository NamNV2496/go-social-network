package repository

import (
	"context"

	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/repository/database"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	AddNotification(ctx context.Context, notification domain.Notification) (int64, error)
	GetNotifications(ctx context.Context, applications []string, pageNumber, pageSize uint32) ([]*domain.Notification, error)
	UpdateNotification(ctx context.Context, notification domain.Notification) error
}

type NotificationRepository struct {
	database.ICRUDBase[*domain.Notification]
}

func NewNotificationRepository(
	db *gorm.DB,
) *NotificationRepository {
	return &NotificationRepository{
		ICRUDBase: database.NewCRUDBase[*domain.Notification](db),
	}
}

var _ INotificationRepository = &NotificationRepository{}

func (_self *NotificationRepository) AddNotification(ctx context.Context, notification domain.Notification) (int64, error) {
	if err := _self.Create(ctx, &notification); err != nil {
		return 0, err
	}
	return notification.Id, nil
}

func (_self *NotificationRepository) GetNotifications(ctx context.Context, applications []string, pageNumber, pageSize uint32) ([]*domain.Notification, error) {
	return _self.Find(ctx, domain.NotificationByApplications(applications), domain.WithPagination(int(pageNumber), int(pageSize)))
}

func (_self *NotificationRepository) UpdateNotification(ctx context.Context, notification domain.Notification) error {
	return _self.Update(ctx, &notification)
}
