package aliyun_mns

type Config struct {
	Url          string   `yaml:"url"`
	AccessKey    string   `yaml:"access_key"`
	AccessSecret string   `yaml:"access_secret"`
	Queues       []string `yaml:"queues"`
}

func (s *Config) ConfigName() string {
	return ServiceName
}

func (s *Config) Validate() error {
	return nil
}

func (s *Config) Default() interface{} {
	return &Config{}
}
