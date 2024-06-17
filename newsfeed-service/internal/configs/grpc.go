package configs

type GRPC struct {
	UserServiceAddress string `yaml:"user_service_address"`
	PostServiceAddress string `yaml:"post_service_address"`
}
