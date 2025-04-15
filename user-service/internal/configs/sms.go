package configs

type SMS struct {
	Enable          bool   `env:"SMS_ENABLE" envDefault:"false"`
	Sid             string `env:"SMS_SID" envDefault:"sid"`
	Password        string `env:"SMS_PASSWORD" envDefault:"password"`
	FromPhoneNumber string `env:"SMS_FROM_PHONE_NUMBER" envDefault:"0999999999"`
}
