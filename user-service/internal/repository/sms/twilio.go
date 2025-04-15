package sms

import (
	"fmt"
	"log/slog"

	"github.com/namnv2496/user-service/internal/configs"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type ISms interface {
	SendOTP(phone string, otp string) error
}

type Sms struct {
	client          *twilio.RestClient
	enable          bool
	fromPhoneNumber string
}

func NewSms(
	conf *configs.Config,
) ISms {
	return &Sms{
		enable:          conf.SMS.Enable,
		fromPhoneNumber: conf.SMS.FromPhoneNumber,
		client: twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: conf.SMS.Sid,
			Password: conf.SMS.Password,
		}),
	}
}

func (c *Sms) SendOTP(phone string, otp string) error {
	if !c.enable {
		return nil
	}
	params := &openapi.CreateMessageParams{}

	params.SetTo(phone)
	params.SetFrom(c.fromPhoneNumber)

	msg := fmt.Sprintf("Your OTP is %s", otp)
	params.SetBody(msg)

	_, err := c.client.Api.CreateMessage(params)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("SMS sent successfully!")
	return nil
}
