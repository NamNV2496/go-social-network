package service

import (
	"context"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/namnv2496/user-service/internal/repository/cache"
	"github.com/namnv2496/user-service/internal/repository/sms"
)

type OTPTemplate struct {
	OTP       int
	ExpiredAt time.Time
}

type IOTP interface {
	SendOTP(ctx context.Context, phone string, userId string) error
	GenerateOTP(ctx context.Context, userId string) string
	VerifyOTP(ctx context.Context, userId string, otp string) bool
}

type OTP struct {
	redisClient cache.Client
	smsClient   sms.ISms
}

func NewOTPService(
	redisClient cache.Client,
	smsClient sms.ISms,
) IOTP {
	return &OTP{
		redisClient: redisClient,
		smsClient:   smsClient,
	}
}

func (c *OTP) SendOTP(ctx context.Context, phone string, userId string) error {
	otp := c.GenerateOTP(ctx, userId)
	slog.Info("send to phone %s OTP: ", "phone", phone, "otp", otp)
	return c.smsClient.SendOTP(phone, otp)
}

func (c *OTP) GenerateOTP(ctx context.Context, userId string) string {
	otp := rand.Intn(9000) + 1000
	cacheOTP, err := c.redisClient.Get(ctx, userId)
	if err == nil {
		return cacheOTP.(string)
	}
	c.redisClient.Set(ctx, userId, otp, 5*time.Minute)
	return strconv.Itoa(otp)
}

func (c *OTP) VerifyOTP(ctx context.Context, userId string, otp string) bool {
	cacheOTP, err := c.redisClient.Get(ctx, userId)
	if err != nil {
		return false
	}
	return otp == cacheOTP.(string)
}
