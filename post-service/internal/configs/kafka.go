package configs

type Kafka struct {
	Addresses string `yaml:"addresses"`
	ClientID  string `yaml:"client_id"`
}
