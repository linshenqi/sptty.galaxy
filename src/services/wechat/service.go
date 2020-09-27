package wechat

import (
	"github.com/linshenqi/sptty"
	"gopkg.in/resty.v1"
)

const (
	ServiceName = "wechat"
)

type Service struct {
	http *resty.Client
}

func (s *Service) Init(app sptty.Sptty) error {
	s.http = sptty.CreateHttpClient(sptty.DefaultHttpClientConfig())

	return nil
}

func (s *Service) Release() {

}

func (s *Service) Enable() bool {
	return true
}

func (s *Service) ServiceName() string {
	return ServiceName
}
