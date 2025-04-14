package email

import "time"

type SendRawQuery struct {
	FromEmail    string    `json:"from_email,omitempty"`
	ToEmail      string    `json:"to_email,omitempty"`
	Cc           string    `json:"cc,omitempty"`
	ReplyToName  string    `json:"reply_to_name,omitempty"`
	ReplyToEmail string    `json:"reply_to_email,omitempty"`
	SchedulerAt  time.Time `json:"scheduler_at,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Body         string    `json:"body,omitempty"`
}

type SendRawQueryResponse struct {
	Success bool `json:"success,omitempty"`
	Payload any  `json:"payload,omitempty"`
}

type SendEmailByTemplate struct {
	FromEmail    string            `json:"from_email,omitempty"`
	ToEmail      string            `json:"to_email,omitempty"`
	Cc           string            `json:"cc,omitempty"`
	ReplyToName  string            `json:"reply_to_name,omitempty"`
	ReplyToEmail string            `json:"reply_to_email,omitempty"`
	SchedulerAt  time.Time         `json:"scheduler_at,omitempty"`
	TemplateId   string            `json:"template_id,omitempty"`
	Params       map[string]string `json:"params,omitempty"`
}

type SendEmailByTemplateResponse struct {
	Success bool `json:"success,omitempty"`
	Payload any  `json:"payload,omitempty"`
}
