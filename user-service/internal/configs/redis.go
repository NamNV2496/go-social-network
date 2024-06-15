package configs

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
	TTL      int    `yaml:"ttl"`
}
