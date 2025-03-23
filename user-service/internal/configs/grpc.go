package configs

type GRPC struct {
	UserServiceAddress string `env:"USER_SERVICE_ADDRESS" envDefault:"localhost:5610"`
	// PostServiceAddress     string `env:"post_service_address"`
	// NewfeedsServiceAddress string `env:"newfeeds_service_address"`
}
