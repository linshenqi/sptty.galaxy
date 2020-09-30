package mq

import "github.com/linshenqi/sptty"

const (
	ServiceName = "mq"
)

type Service struct {
	sptty.BaseService
}

func (s *Service) Init(app sptty.Sptty) error {
	return nil
}

func (s *Service) Release() {
}

func (s *Service) ServiceName() string {
	return ServiceName
}
