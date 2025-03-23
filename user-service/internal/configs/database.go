package configs

type DatabaseType string

const (
	DatabaseTypeMySQL DatabaseType = "mysql"
)

type Database struct {
	Type     DatabaseType `env:"DATABASE_TYPE" envDefault:"mysql"`
	Host     string       `env:"DATABASE_HOST" envDefault:"127.0.0.1"`
	Port     int          `env:"DATABASE_PORT" envDefault:"3309"`
	Username string       `env:"DATABASE_USERNAME" envDefault:"root"`
	Password string       `env:"DATABASE_PASSWORD" envDefault:"root"`
	Database string       `env:"DATABASE_DATABASE" envDefault:"network"`
}
