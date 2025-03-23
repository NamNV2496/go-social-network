package configs

type Redis struct {
	Address  string `env:"REDIS_URL" envDefault:"127.0.0.1:6379"`
	Username string `env:"REDIS_USERNAME" envDefault:""`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	Database int    `env:"REDIS_DATABASW" envDefault:"0"`
	TTL      int    `env:"REDIS_TTL" envDefault:"3000"`
}
