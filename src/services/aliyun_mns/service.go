package aliyun_mns

import (
	"fmt"
	"sync"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/linshenqi/sptty"
)

const (
	ServiceName = "aliyun_mns"
)

type Service struct {
	cfg       Config
	mnsClient ali_mns.MNSClient

	queueContext    map[string]*Queue
	mtxQueueContext sync.Mutex
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
}

func (s *Service) Enable() bool {
	return true
}

func (s *Service) ServiceName() string {
	return ServiceName
}

func (s *Service) doInit() error {
	// init client
	s.mnsClient = ali_mns.NewAliMNSClient(s.cfg.Url, s.cfg.AccessKey, s.cfg.AccessSecret)

	// init queues
	s.mtxQueueContext = sync.Mutex{}
	s.queueContext = map[string]*Queue{}
	for _, v := range s.cfg.Queues {
		q := Queue{
			Queue:   ali_mns.NewMNSQueue(v, s.mnsClient),
			RecvBuf: make(chan ali_mns.MessageReceiveResponse, BufferSize),
			ErrBuf:  make(chan error),
		}

		s.queueContext[v] = &q
		go s.asyncQueueHanlder(&q)
	}

	return nil
}

func (s *Service) asyncQueueHanlder(queue *Queue) {
	name := queue.Queue.Name()

	for {
		select {
		case recv := <-queue.RecvBuf:
			s.notifyQueueHandlers(name, recv.MessageBody, nil)
		case err := <-queue.ErrBuf:
			s.notifyQueueHandlers(name, "", err)
			return
		}
	}
}

func (s *Service) getQueue(name string) (*Queue, error) {
	q, exist := s.queueContext[name]
	if !exist {
		return nil, fmt.Errorf("Queue %s Not Found ", name)
	}

	return q, nil
}

func (s *Service) notifyQueueHandlers(queueName string, msg string, err error) {
	q, err := s.getQueue(queueName)
	if err != nil {
		return
	}

	s.mtxQueueContext.Lock()
	defer s.mtxQueueContext.Unlock()

	for _, handler := range q.Handlers {
		handler(msg, err)
	}
}

func (s *Service) PostQueueMsg(queueName string, msg string) error {

	q, err := s.getQueue(queueName)
	if err != nil {
		return err
	}

	s.mtxQueueContext.Lock()
	defer s.mtxQueueContext.Unlock()

	_, err = q.Queue.SendMessage(ali_mns.MessageSendRequest{
		MessageBody:  msg,
		DelaySeconds: 0,
		Priority:     8,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddQueueHandler(queueName string, handler QueueHandler) error {

	q, err := s.getQueue(queueName)
	if err != nil {
		return err
	}

	s.mtxQueueContext.Lock()
	defer s.mtxQueueContext.Unlock()

	q.Handlers = append(q.Handlers, handler)

	return nil
}
