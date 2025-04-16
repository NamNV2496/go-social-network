package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/domain"
	"github.com/namnv2496/user-service/internal/repository/repo"
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
	SendEmail(ctx context.Context, request *SendEmailRequest) (*domain.SendRawQueryResponse, error)
	SendEmailByTemplate(ctx context.Context, request *domain.SendEmailByTemplate) (*domain.SendEmailByTemplateResponse, error)

	GetEmailTemplateById(context.Context, *domain.GetEmailTemplateRequest) (*domain.GetEmailTemplateResponse, error)
	GetEmailTemplateByTemplateId(context.Context, *domain.GetEmailTemplateByTemplateIdRequest) (*domain.GetEmailTemplateResponse, error)
	AddEmailTemplate(ctx context.Context, req *domain.AddEmailTemplateRequest) (*domain.AddEmailTemplateResponse, error)
	UpdateEmailTemplate(ctx context.Context, req *domain.UpdateEmailTemplateRequest) (*domain.UpdateEmailTemplateResponse, error)
}

type Email struct {
	emailClient *email.Email
	emailRepo   repo.IEmailTemplateRepo
	cfg         configs.Email
}

func NewEmailClient(
	conf *configs.Config,
	emailRepo repo.IEmailTemplateRepo,
) IEmail {
	return &Email{
		cfg:         conf.Email,
		emailClient: email.NewEmail(),
		emailRepo:   emailRepo,
	}
}

func (e *Email) SendEmail(ctx context.Context, request *SendEmailRequest) (*domain.SendRawQueryResponse, error) {
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
		return &domain.SendRawQueryResponse{
			Success: false,
			Payload: "fail",
		}, err
	}
	return &domain.SendRawQueryResponse{
		Success: true,
		Payload: "success",
	}, nil
}

func (e *Email) SendEmailByTemplate(ctx context.Context, request *domain.SendEmailByTemplate) (*domain.SendEmailByTemplateResponse, error) {
	templateBody := e.getTemplate(ctx, request.TemplateId)
	body := parseTemplate(templateBody, request.Params)
	slog.Info("send email with body", "body", body)

	e.emailClient.From = fmt.Sprintf("%s <%s>", "Registation account", request.FromEmail)
	e.emailClient.Subject = "Registation account confirmation"
	e.emailClient.HTML = []byte(body)
	e.emailClient.To = []string{request.ToEmail}
	e.emailClient.Cc = []string{request.Cc}
	e.emailClient.Bcc = []string{}

	// smtpAuth := smtp.PlainAuth("", e.cfg.Email, e.cfg.Password, smtpAuthAddress)
	// if err := e.emailClient.Send(smtpServerAddress, smtpAuth); err != nil {
	// 	return &domain.SendEmailByTemplateResponse{
	// 		Success: false,
	// 		Payload: "fail",
	// 	}, err
	// }
	return &domain.SendEmailByTemplateResponse{
		Success: true,
		Payload: "success",
	}, nil
}

func (e *Email) GetEmailTemplateById(ctx context.Context, req *domain.GetEmailTemplateRequest) (*domain.GetEmailTemplateResponse, error) {
	result, err := e.emailRepo.GetEmailTemplateById(ctx, req)
	if err != nil {
		return nil, err
	}
	emailTemplates := make([]*domain.GetEmailTemplate, 0, len(result))
	for _, template := range result {
		emailTemplates = append(emailTemplates, &domain.GetEmailTemplate{
			Id:         uint64(template.Id),
			TemplateId: template.TemplateId,
			Template:   template.Template,
		})
	}
	return &domain.GetEmailTemplateResponse{
		Response: emailTemplates,
	}, nil
}

func (e *Email) GetEmailTemplateByTemplateId(ctx context.Context, req *domain.GetEmailTemplateByTemplateIdRequest) (*domain.GetEmailTemplateResponse, error) {
	result, err := e.emailRepo.GetEmailTemplateByTemplateId(ctx, req)
	if err != nil {
		return nil, err
	}
	emailTemplates := make([]*domain.GetEmailTemplate, 0, len(result))
	for _, template := range result {
		emailTemplates = append(emailTemplates, &domain.GetEmailTemplate{
			Id:         uint64(template.Id),
			TemplateId: template.TemplateId,
			Template:   template.Template,
		})
	}
	return &domain.GetEmailTemplateResponse{
		Response: emailTemplates,
	}, nil
}

func (e *Email) AddEmailTemplate(ctx context.Context, req *domain.AddEmailTemplateRequest) (*domain.AddEmailTemplateResponse, error) {
	result, err := e.emailRepo.AddEmailTemplate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.AddEmailTemplateResponse{
		Status: result,
	}, nil
}

func (e *Email) UpdateEmailTemplate(ctx context.Context, req *domain.UpdateEmailTemplateRequest) (*domain.UpdateEmailTemplateResponse, error) {
	result, err := e.emailRepo.UpdateEmailTemplate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &domain.UpdateEmailTemplateResponse{
		Status: result,
	}, nil
}

func (e *Email) getTemplate(ctx context.Context, templateId string) string {
	result, err := e.emailRepo.GetEmailTemplateByTemplateId(ctx, &domain.GetEmailTemplateByTemplateIdRequest{
		TemplateId: templateId,
	})
	if err != nil {
		return ""
	}
	return result[0].Template
}

func parseTemplate(templateformat string, input map[string]string) string {
	body := strings.Replace(templateformat, "\n", "<br />", -1)
	body = strings.Replace(body, "<b>", "", -1)
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
