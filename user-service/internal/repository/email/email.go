package email

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"github.com/namnv2496/user-service/internal/configs"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type SendEmailRequest struct {
	FromEmail    string    `json:"from_email,omitempty"`
	ToEmail      string    `json:"to_email,omitempty"`
	Cc           string    `json:"cc,omitempty"`
	ReplyToName  string    `json:"reply_to_name,omitempty"`
	ReplyToEmail string    `json:"reply_to_email,omitempty"`
	SchedulerAt  time.Time `json:"scheduler_at,omitempty"`
	// option 1
	Subject string `json:"subject,omitempty" validate:"required_without=Template"`
	Body    string `json:"body,omitempty" validate:"required_with=Subject"`
	// option 2
	TemplateId string            `json:"template_id,omitempty" validate:"required_without=Subject"`
	Params     map[string]string `json:"params,omitempty" validate:"required_with=Template"`
}

type SendEmailResponse struct {
	Id     string
	Status string
}

type IEmail interface {
	SendEmail(ctx context.Context, request *SendEmailRequest) (*SendRawQueryResponse, error)
	SendEmailByTemplate(ctx context.Context, request SendEmailByTemplate) (*SendEmailByTemplateResponse, error)
}

type Email struct {
	emailClient *email.Email
	cfg         configs.Email
}

func NewEmailClient(
	conf *configs.Config,
) IEmail {
	return &Email{
		cfg:         conf.Email,
		emailClient: email.NewEmail(),
	}
}

func (e *Email) SendEmail(ctx context.Context, request *SendEmailRequest) (*SendRawQueryResponse, error) {
	// requestBody := &SendRawQuery{
	// 	FromEmail:    request.FromEmail,
	// 	ToEmail:      request.ToEmail,
	// 	Cc:           request.Cc,
	// 	ReplyToName:  request.ReplyToName,
	// 	ReplyToEmail: request.ReplyToEmail,
	// 	SchedulerAt:  request.SchedulerAt,
	// }
	e.emailClient.From = fmt.Sprintf("%s <%s>", "Registation account", request.FromEmail)
	e.emailClient.Subject = request.Subject
	e.emailClient.HTML = []byte(request.Body)
	e.emailClient.To = []string{request.ToEmail}
	e.emailClient.Cc = []string{request.Cc}
	e.emailClient.Bcc = []string{}

	smtpAuth := smtp.PlainAuth("", e.cfg.Email, e.cfg.Password, smtpAuthAddress)
	if err := e.emailClient.Send(smtpServerAddress, smtpAuth); err != nil {
		return &SendRawQueryResponse{
			Success: false,
			Payload: "fail",
		}, err
	}
	return &SendRawQueryResponse{
		Success: true,
		Payload: "success",
	}, nil
}

func (e *Email) SendEmailByTemplate(ctx context.Context, request SendEmailByTemplate) (*SendEmailByTemplateResponse, error) {
	templateBody := getTemplate(request.TemplateId)
	body := parseTemplate(templateBody, request.Params)

	e.emailClient.From = fmt.Sprintf("%s <%s>", "Registation account", request.FromEmail)
	e.emailClient.Subject = "Registation account confirmation"
	e.emailClient.HTML = []byte(body)
	e.emailClient.To = []string{request.ToEmail}
	e.emailClient.Cc = []string{request.Cc}
	e.emailClient.Bcc = []string{}

	smtpAuth := smtp.PlainAuth("", e.cfg.Email, e.cfg.Password, smtpAuthAddress)
	if err := e.emailClient.Send(smtpServerAddress, smtpAuth); err != nil {
		return &SendEmailByTemplateResponse{
			Success: false,
			Payload: "fail",
		}, err
	}
	return &SendEmailByTemplateResponse{
		Success: true,
		Payload: "success",
	}, nil
}

func getTemplate(templateId string) string {
	// get template by template id in DB
	return "<b>Hi {{.full_name}} this is your OTP: {{.otp}}, please don't publish this OTP to another people"
}

func parseTemplate(templateformat string, input map[string]string) string {
	body := strings.Replace(templateformat, "\n", "<br />", -1)
	body = strings.Replace(body, "[", "{{.", -1)
	body = strings.Replace(body, "]", "}}.", -1)
	body = strings.Replace(body, `\"`, `"`, -1)
	t, err := template.New("parse template").Parse(body)
	if err != nil {
		return ""
	}
	var b bytes.Buffer
	err = t.Execute(&b, &input)
	if err != nil {
		return ""
	}
	return b.String()
}
