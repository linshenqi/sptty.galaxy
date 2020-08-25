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
	Handlers []QueueHandler
}

func (s *Queue) doRecv() {
	s.Queue.ReceiveMessage(s.RecvBuf, s.ErrBuf, 8)
}
