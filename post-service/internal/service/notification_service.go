package service

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"html/template"
	"reflect"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/namnv2496/post-service/internal/domain"
	"github.com/namnv2496/post-service/internal/entity"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/spf13/cast"
)

var (
	defaultFuncs = template.FuncMap{
		"suffix": suffixFunc,
		"prefix": prefixFunc,
	}
)

func suffixFunc(a, b any) bool {
	s1 := cast.ToString(a)
	s2 := cast.ToString(b)
	return strings.HasSuffix(s1, s2)
}

func prefixFunc(a, b any) bool {
	s1 := cast.ToString(a)
	s2 := cast.ToString(b)
	return strings.HasPrefix(s1, s2)
}

type INotificationService interface {
	CreateNotification(ctx context.Context, req *entity.CreateNotificationRequest) (int64, error)
	GetNotifications(ctx context.Context, req *entity.GetNotificationsRequest) ([]*entity.Notification, error)
	UpdateNotification(ctx context.Context, req *entity.UpdateNotificationRequest) error
	Notify(ctx context.Context, req *entity.NotifyRequest) error
}

type notificationService struct {
	notificationRepository repository.INotificationRepository
}

func NewNotificationService(notificationRepository repository.INotificationRepository) INotificationService {
	return &notificationService{notificationRepository: notificationRepository}
}

var _ INotificationService = &notificationService{}

func (s *notificationService) CreateNotification(ctx context.Context, req *entity.CreateNotificationRequest) (int64, error) {
	notification := domain.Notification{
		Title:       req.Notification.Title,
		Description: req.Notification.Description,
		Template:    req.Notification.Template,
		Image:       req.Notification.Image,
		Application: req.Notification.Application,
		Visible:     req.Notification.Visible,
		Link:        req.Notification.Link,
	}
	return s.notificationRepository.AddNotification(ctx, notification)
}

func (s *notificationService) GetNotifications(ctx context.Context, req *entity.GetNotificationsRequest) ([]*entity.Notification, error) {
	result, err := s.notificationRepository.GetNotifications(ctx, req.Application, req.PageNumber, req.PageSize)
	if err != nil {
		return nil, err
	}
	var resp []*entity.Notification
	for _, notification := range result {
		resp = append(resp, &entity.Notification{
			Id:          notification.Id,
			Title:       notification.Title,
			Description: notification.Description,
			Template:    notification.Template,
			Image:       notification.Image,
			Application: notification.Application,
			Visible:     notification.Visible,
			Link:        notification.Link,
		})
	}
	return resp, nil
}

func (s *notificationService) UpdateNotification(ctx context.Context, req *entity.UpdateNotificationRequest) error {
	notification := domain.Notification{
		Id:          req.Id,
		Title:       req.Notification.Title,
		Description: req.Notification.Description,
		Template:    req.Notification.Template,
		Image:       req.Notification.Image,
		Application: req.Notification.Application,
		Visible:     req.Notification.Visible,
		Link:        req.Notification.Link,
	}
	return s.notificationRepository.UpdateNotification(ctx, notification)
}

func (s *notificationService) Notify(ctx context.Context, req *entity.NotifyRequest) error {
	notificattionTemplate, err := s.getNotificationById(ctx, req.Application, req.Id)
	if err != nil {
		return err
	}
	fullfiled := handleTemplate(notificattionTemplate, req.Data)
	fmt.Println(fullfiled)
	fmt.Println("Notification sent successfully")
	return nil
}

func (s *notificationService) getNotificationById(ctx context.Context, application string, id int64) (*entity.Notification, error) {
	result, err := s.notificationRepository.GetNotifications(ctx, []string{application}, 0, 1)
	if err != nil {
		return nil, err
	}
	return &entity.Notification{
		Id:          result[0].Id,
		Title:       result[0].Title,
		Description: result[0].Description,
		Template:    result[0].Template,
		Image:       result[0].Image,
		Application: result[0].Application,
		Visible:     result[0].Visible,
		Link:        result[0].Link,
	}, nil
}

func handleTemplate(input *entity.Notification, params map[string]string) *entity.Notification {
	fields := reflect.ValueOf(input).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fieldName := fields.Type().Field(i).Name
		value := fields.Field(i)
		if value.Kind() == reflect.String {
			// content := sanitizeHtml(value.String())
			// value.SetString(replacePlaceholders(fieldName, content, params))
			value.SetString(replacePlaceholders(fieldName, value.String(), params))
		}
	}
	return input
}

func replacePlaceholders(fieldName string, content string, params map[string]string) string {
	template, err := template.New(fieldName).
		Option("missingkey=zero").
		Funcs(defaultFuncs).
		Parse(content)
	if err != nil {
		return content
	}
	result := new(bytes.Buffer)
	err = template.Execute(result, params)
	if err != nil {
		return content
	}
	return result.String()
}

func sanitizeHtml(input string) string {
	text := bluemonday.StrictPolicy().Sanitize(input)
	return html.UnescapeString(text)
}
