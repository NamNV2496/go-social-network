package configs

type Logger struct {
	Development       bool   `yaml:"development"`
	DisableCaller     bool   `yaml:"disable_caller"`
	DisableStacktrace bool   `yaml:"disable_stacktrace"`
	Encoding          string `yaml:"encoding"`
	Level             string `yaml:"level"`
}
