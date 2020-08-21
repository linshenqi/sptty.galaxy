package aliyun_mns

import (
	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
)

const (
	BufferSize = 65535
)

type QueueHandler func(string, error)

type Queue struct {
	RecvBuf  chan ali_mns.MessageReceiveResponse
	Queue    ali_mns.AliMNSQueue
	ErrBuf   chan error
	Handlers []QueueHandler
}
