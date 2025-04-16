package domain

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

type GetEmailTemplate struct {
	Id         uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Template   string `protobuf:"bytes,2,opt,name=template,proto3" json:"template,omitempty"`
	TemplateId string `protobuf:"bytes,2,opt,name=template_id,proto3" json:"template_id,omitempty"`
}

type GetEmailTemplateRequest struct {
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

type GetEmailTemplateResponse struct {
	Response []*GetEmailTemplate `protobuf:"bytes,1,rep,name=response,proto3" json:"response,omitempty"`
}

type AddEmailTemplateRequest struct {
	Template *GetEmailTemplate `protobuf:"bytes,1,opt,name=template,proto3" json:"template,omitempty"`
}

type AddEmailTemplateResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

type UpdateEmailTemplateRequest struct {
	Template *GetEmailTemplate `protobuf:"bytes,1,opt,name=template,proto3" json:"template,omitempty"`
}

type UpdateEmailTemplateResponse struct {
	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

type GetEmailTemplateByTemplateIdRequest struct {
	TemplateId string `protobuf:"bytes,1,opt,name=template_id,json=templateId,proto3" json:"template_id,omitempty"`
}
