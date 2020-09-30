package nats

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/linshenqi/sptty"
	"github.com/nats-io/nats.go"
)

const (
	ServiceName = "nats"
)

type Service struct {
	sptty.BaseService

	cfg                Config
	queueContext       map[string]*Queue
	mtxQueueContext    sync.Mutex
	nc                 *nats.Conn
	done               chan bool
	pendingHandlers    []PendingQueueHandler
	mtxPendingHandlers sync.Mutex
}

func (s *Service) Init(app sptty.Sptty) error {
	if err := app.GetConfig(s.ServiceName(), &s.cfg); err != nil {
		return err
	}

	if err := s.doInit(); err != nil {
		return err
	}

	return nil
}

func (s *Service) Release() {
	s.nc.Close()
}

func (s *Service) ServiceName() string {
	return ServiceName
}

func (s *Service) PostQueueMsg(queueName string, data []byte) error {
	if err := s.nc.Publish(queueName, data); err != nil {
		return err
	}

	return nil
}

func (s *Service) AddQueueHandler(queueName string, handler QueueHandler) error {
	if !s.connected() {
		s.addPendingHandler(queueName, handler)
		return nil
	}

	q, err := s.getOrCreateQueue(queueName)
	if err != nil {
		return err
	}

	q.AddHandler(handler)
	return nil
}

func (s *Service) connected() bool {
	if s.nc == nil {
		return false
	}

	return s.nc.IsConnected()
}

func (s *Service) addPendingHandler(queueName string, handler QueueHandler) {
	s.mtxPendingHandlers.Lock()
	defer s.mtxPendingHandlers.Unlock()

	s.pendingHandlers = append(s.pendingHandlers, PendingQueueHandler{
		QueueName: queueName,
		Handler:   handler,
	})
}

func (s *Service) handlePendingHandlers() {
	if !s.connected() {
		return
	}

	s.mtxPendingHandlers.Lock()
	defer s.mtxPendingHandlers.Unlock()

	for _, v := range s.pendingHandlers {
		s.AddQueueHandler(v.QueueName, v.Handler)
	}
}

func (s *Service) doInit() error {
	s.nc = nil
	s.mtxPendingHandlers = sync.Mutex{}
	s.queueContext = map[string]*Queue{}
	s.mtxQueueContext = sync.Mutex{}

	go s.asyncInitNats()

	return nil
}

func (s *Service) asyncInitNats() {
	opts := s.setupOpts()
	urls := strings.Join(s.cfg.Urls, ",")

	for {
		select {
		case <-time.After(1 * time.Second):
			nc, err := nats.Connect(urls, opts...)
			if err == nil {
				s.nc = nc
				s.handlePendingHandlers()
				return
			}

		case <-s.done:
			return

		}
	}
}

func (s *Service) setupOpts() []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := 1 * time.Second

	opts := []nats.Option{nats.Name(s.cfg.Name)}
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		sptty.Log(sptty.InfoLevel, fmt.Sprintf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes()), s.ServiceName())
	}))

	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		sptty.Log(sptty.InfoLevel, fmt.Sprintf("Reconnected [%s]", nc.ConnectedUrl()), s.ServiceName())
	}))

	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		sptty.Log(sptty.InfoLevel, fmt.Sprintf("Exiting: %v", nc.LastError()), s.ServiceName())
	}))

	opts = append(opts, nats.UserInfo(s.cfg.User, s.cfg.Pwd))

	return opts
}

func (s *Service) getOrCreateQueue(queueName string) (*Queue, error) {
	s.mtxQueueContext.Lock()
	defer s.mtxQueueContext.Unlock()

	var err error
	q, exist := s.queueContext[queueName]
	if !exist {
		q, err = s.doCreateQueue(queueName)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (s *Service) doCreateQueue(queueName string) (*Queue, error) {
	q := Queue{}
	if err := q.Init(queueName, s.nc); err != nil {
		return nil, err
	}

	s.queueContext[queueName] = &q

	return &q, nil
}
