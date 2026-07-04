package configs

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

//go:embed local.yaml
var DefaultConfigBytes []byte

type Config struct {
	GRPC        GRPC        `yaml:"grpc"`
	Kafka       Kafka       `yaml:"kafka"`
	Database    Database    `yaml:"database"`
	Logger      Logger      `yaml:"logger"`
	CommentRule CommentRule `yaml:"comment_rule"`
}

type CommentRule struct {
	Mode          string        `yaml:"mode"`
	AI            AIConfig      `yaml:"ai"`
	Langfuse      Langfuse      `yaml:"langfuse"`
	LangfuseTrace LangfuseTrace `yaml:"langfuse_trace"`
}

type AIConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Endpoint    string `yaml:"endpoint"`
	Model       string `yaml:"model"`
	APIKey      string `yaml:"api_key"`
	ProjectName string `yaml:"project_name"`
	PromptName  string `yaml:"prompt_name"`
	// PromptLabelName string `yaml:"prompt_label_name"`
}

type Langfuse struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	PublicKey string `yaml:"public_key"`
	SecretKey string `yaml:"secret_key"`
}

type LangfuseTrace struct {
	Endpoint  string `yaml:"endpoint"`
	PublicKey string `yaml:"public_key"`
	SecretKey string `yaml:"secret_key"`
	TraceName string `yaml:"trace_name"`
	// VerboseDump enables base64 dumping of OTLP request/response bodies for debugging.
	VerboseDump bool `yaml:"verbose_dump"`
}

func NewConfig() (Config, error) {
	var (
		configBytes []byte = DefaultConfigBytes
		config      Config
		err         error
	)
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshal configuration file: %w", err)
	}

	return config, nil
}
