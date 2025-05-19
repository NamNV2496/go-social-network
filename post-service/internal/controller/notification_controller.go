package controller

import (
	"context"

	"github.com/namnv2496/post-service/internal/entity"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg"
	"github.com/namnv2496/post-service/internal/service"
)

type NotificationController struct {
	postv1.UnimplementedNotificationServiceServer
	notificationService service.INotificationService
}

func NewNotificationController(
	notificationService service.INotificationService,
) postv1.NotificationServiceServer {
	return &NotificationController{
		notificationService: notificationService,
	}
}

var _ postv1.NotificationServiceServer = &NotificationController{}

func (c *NotificationController) CreateNotification(ctx context.Context, req *postv1.CreateNotificationRequest) (*postv1.CreateNotificationResponse, error) {
	if req.Notification == nil {
		return nil, nil
	}
	var request entity.CreateNotificationRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	_, err := c.notificationService.CreateNotification(ctx, &request)
	if err != nil {
		return nil, err
	}
	return &postv1.CreateNotificationResponse{
		Status: "success",
	}, nil
}

func (c *NotificationController) GetNotifications(ctx context.Context, req *postv1.GetNotificationsRequest) (*postv1.GetNotificationsResponse, error) {
	var request entity.GetNotificationsRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	notifications, err := c.notificationService.GetNotifications(ctx, &request)
	if err != nil {
		return nil, err
	}
	var resp []*postv1.Notification
	for _, notification := range notifications {
		resp = append(resp, &postv1.Notification{
			Id:          notification.Id,
			Title:       notification.Title,
			Description: notification.Description,
			Template:    notification.Template,
			Image:       notification.Image,
			Application: notification.Application,
			Visible:     notification.Visible,
		})
	}
	return &postv1.GetNotificationsResponse{
		Notifications: resp,
	}, nil
}

func (c *NotificationController) UpdateNotifications(ctx context.Context, req *postv1.UpdateNotificationsRequest) (*postv1.UpdateNotificationsResponse, error) {
	var request entity.UpdateNotificationRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	err := c.notificationService.UpdateNotification(ctx, &request)
	if err != nil {
		return nil, err
	}
	return &postv1.UpdateNotificationsResponse{
		Status: "success",
	}, nil
}

func (c *NotificationController) Notify(ctx context.Context, req *postv1.NotifyRequest) (*postv1.NotifyResponse, error) {
	var request entity.NotifyRequest
	if err := pkg.Copy(&request, req); err != nil {
		return nil, err
	}
	err := c.notificationService.Notify(ctx, &request)
	if err != nil {
		return nil, err
	}
	return &postv1.NotifyResponse{
		Status: "success",
	}, nil
}
