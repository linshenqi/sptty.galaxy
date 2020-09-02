package aliyun_mns

import (
	alimns "github.com/aliyun/aliyun-mns-go-sdk"
)

const (
	BufferSize = 65535
)

type QueueHandler func(string, *alimns.MessageReceiveResponse, error)

type Queue struct {
	RecvBuf  chan alimns.MessageReceiveResponse
	Queue    alimns.AliMNSQueue
	ErrBuf   chan error
	Done     chan bool
	Handlers []QueueHandler
}

func (s *Queue) doRecv() {
	s.Queue.ReceiveMessage(s.RecvBuf, s.ErrBuf, 8)
}

func (s *Queue) release() {
	s.Done <- true
}
