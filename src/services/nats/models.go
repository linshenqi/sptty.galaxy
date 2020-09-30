package nats

import (
	"sync"

	"github.com/nats-io/nats.go"
)

const (
	BufferSize = 65535
)

type PendingQueueHandler struct {
	QueueName string
	Handler   QueueHandler
}

// topic, msg
type QueueHandler func(string, []byte)

type Queue struct {
	name     string
	nc       *nats.Conn
	mtx      sync.Mutex
	handlers []QueueHandler
	recvBuf  chan *nats.Msg
	done     chan bool
}

// func (s *Queue) rawRecvHandler(msg *nats.Msg) {
// 	if msg == nil {
// 		return
// 	}

// 	s.recvBuf <- msg
// }

func (s *Queue) asyncHandleRecv() {
	for {
		select {
		case <-s.done:
			return

		case msg := <-s.recvBuf:
			s.doNotify(msg)
		}
	}
}

func (s *Queue) doNotify(msg *nats.Msg) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	for _, v := range s.handlers {
		if v != nil {
			v(msg.Subject, msg.Data)
		}
	}
}

func (s *Queue) Init(name string, nc *nats.Conn) error {
	s.nc = nc
	s.mtx = sync.Mutex{}
	s.recvBuf = make(chan *nats.Msg, BufferSize)
	s.done = make(chan bool)

	_, err := s.nc.QueueSubscribeSyncWithChan(name, name, s.recvBuf)
	if err != nil {
		return err
	}

	go s.asyncHandleRecv()

	return nil
}

func (s *Queue) AddHandler(handler QueueHandler) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.handlers = append(s.handlers, handler)
}
