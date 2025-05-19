package entity

type Notification struct {
	Id          int64  `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Template    string `json:"template,omitempty"`
	Image       string `json:"image,omitempty"`
	Application string `json:"application,omitempty"`
	Visible     bool   `json:"visible,omitempty"`
	Link        string `json:"link,omitempty"`
}

type CreateNotificationRequest struct {
	Notification *Notification `json:"notification,omitempty"`
}

type CreateNotificationResponse struct {
	NotificationId int64  `json:"notification_id,omitempty"`
	Status         string `json:"status,omitempty"`
}

type GetNotificationsRequest struct {
	Application []string `json:"application,omitempty"`
	PageNumber  uint32   `json:"page_number,omitempty"`
	PageSize    uint32   `json:"page_size,omitempty"`
}

type GetNotificationsResponse struct {
	Notifications []*Notification `json:"notifications,omitempty"`
}

type UpdateNotificationRequest struct {
	Id           int64         `json:"id,omitempty"`
	Notification *Notification `json:"notification,omitempty"`
}

type UpdateNotificationResponse struct {
	Status string `json:"status,omitempty"`
}

type NotifyRequest struct {
	UserId      string            `json:"user_id,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Id          int64             `json:"id,omitempty"`
	Application string            `json:"application,omitempty"`
}

type NotifyResponse struct {
	Status string `json:"status,omitempty"`
}
