package configs

import "time"

type Auth struct {
	Duration string `env:"duration" envDefault:"24h"`
}

func (t *Auth) GetTokenDuration() time.Duration {
	duration, _ := time.ParseDuration(t.Duration)
	return duration
}
