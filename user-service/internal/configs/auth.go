package configs

import "time"

type Auth struct {
	Duration string `yaml:"duration"`
}

func (t *Auth) GetTokenDuration() time.Duration {
	duration, _ := time.ParseDuration(t.Duration)
	return duration
}
