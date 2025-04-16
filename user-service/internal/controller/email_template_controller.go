package controller

import (
	"context"

	"github.com/namnv2496/user-service/internal/domain"
	"github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
)

type EmailTemplateHanlder struct {
	userv1.UnimplementedEmailTemplateServiceServer
	emailService service.IEmail
}

func NewEmailTemplateHander(
	emailService service.IEmail,
) userv1.EmailTemplateServiceServer {
	return &EmailTemplateHanlder{
		emailService: emailService,
	}
}

func (c *EmailTemplateHanlder) GetEmailTemplateById(ctx context.Context, req *userv1.GetEmailTemplateRequest) (*userv1.GetEmailTemplateResponse, error) {
	resp, err := c.emailService.GetEmailTemplateById(ctx, &domain.GetEmailTemplateRequest{
		Id: req.Id,
	})
	if err != nil || resp == nil {
		return nil, err
	}
	response := make([]*userv1.GetEmailTemplate, 0, len(resp.Response))
	for _, temp := range resp.Response {
		response = append(response, &userv1.GetEmailTemplate{
			Id:         temp.Id,
			Template:   temp.Template,
			TemplateId: temp.TemplateId,
		})
	}
	return &userv1.GetEmailTemplateResponse{
		Response: response,
	}, nil
}
func (c *EmailTemplateHanlder) GetEmailTemplateByTemplateId(ctx context.Context, req *userv1.GetEmailTemplateByTemplateIdRequest) (*userv1.GetEmailTemplateResponse, error) {
	resp, err := c.emailService.GetEmailTemplateByTemplateId(ctx, &domain.GetEmailTemplateByTemplateIdRequest{
		TemplateId: req.TemplateId,
	})
	if err != nil || resp == nil {
		return nil, err
	}
	response := make([]*userv1.GetEmailTemplate, 0, len(resp.Response))
	for _, temp := range resp.Response {
		response = append(response, &userv1.GetEmailTemplate{
			Id:         temp.Id,
			Template:   temp.Template,
			TemplateId: temp.TemplateId,
		})
	}
	return &userv1.GetEmailTemplateResponse{
		Response: response,
	}, nil
}

func (c *EmailTemplateHanlder) AddEmailTemplate(ctx context.Context, req *userv1.AddEmailTemplateRequest) (*userv1.AddEmailTemplateResponse, error) {
	resp, err := c.emailService.AddEmailTemplate(ctx, &domain.AddEmailTemplateRequest{
		Template: &domain.GetEmailTemplate{
			Id:         req.Template.Id,
			Template:   req.Template.Template,
			TemplateId: req.Template.TemplateId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &userv1.AddEmailTemplateResponse{
		Status: resp.Status,
	}, nil
}

func (c *EmailTemplateHanlder) UpdateEmailTemplate(ctx context.Context, req *userv1.UpdateEmailTemplateRequest) (*userv1.UpdateEmailTemplateResponse, error) {
	resp, err := c.emailService.UpdateEmailTemplate(ctx, &domain.UpdateEmailTemplateRequest{
		Template: &domain.GetEmailTemplate{
			Id:         req.Template.Id,
			Template:   req.Template.Template,
			TemplateId: req.Template.TemplateId,
		},
	})
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateEmailTemplateResponse{
		Status: resp.Status,
	}, nil
}
