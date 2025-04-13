package configs

type DatabaseType string

const (
	DatabaseTypeMySQL DatabaseType = "mysql"
)

type Database struct {
	Type        DatabaseType `yaml:"type"`
	Host        string       `yaml:"host"`
	Port        int          `yaml:"port"`
	Username    string       `yaml:"username"`
	Password    string       `yaml:"password"`
	Database    string       `yaml:"database"`
	AutoMigrate bool         `yaml:"auto_migrate"`
	Logger      Logger       `yaml:"logger"`
}
