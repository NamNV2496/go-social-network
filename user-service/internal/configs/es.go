package configs

type ElasticSearch struct {
	Addr          []string `env:"ES_ADDR" envDefault:"http://51.79.250.32:9200"`
	Username      string   `env:"ES_USERNAME" envDefault:"elastic"`
	Password      string   `env:"ES_PASSWORD" envDefault:"admin"`
	RetryOnStatus []int    `env:"ES_RETRY_ON_STATUS" envDefault:"429,502,503,504"`
	DisableRetry  bool     `env:"ES_DISABLE_RETRY" envDefault:"false"`
	MaxRetries    int      `env:"ES_MAX_RETRIES" envDefault:"3"`
	RetryBackoff  int      `env:"ES_RETRY_BACKOFF" envDefault:"100"`
}
