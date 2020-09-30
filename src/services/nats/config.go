package nats

import "github.com/linshenqi/sptty"

type Config struct {
	sptty.BaseConfig

	Name string   `yaml:"name"`
	Urls []string `yaml:"urls"`
	User string   `yaml:"user"`
	Pwd  string   `yaml:"pwd"`
}

func (s *Config) ConfigName() string {
	return ServiceName
}

func (s *Config) Default() interface{} {
	return &Config{}
}
