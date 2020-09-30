package nats

import (
	"fmt"
	"testing"
	"time"
)

const (
	q1 = "q1"
)

func getSrv() *Service {
	srv := Service{cfg: Config{
		Name: "n1",
		Urls: []string{
			"nats://118.190.133.189:4222",
			"nats://118.190.133.189:4223",
			"nats://118.190.133.189:4224",
		},
		User: "admin",
		Pwd:  "T0pS3cr3t",
	}}

	srv.doInit()

	return &srv
}

func onRecv1(topic string, data []byte) {
	fmt.Printf("onRecv1 topic: %s data: %s\n", topic, string(data))
}

func onRecv2(topic string, data []byte) {
	fmt.Printf("onRecv2 topic: %s data: %s\n", topic, string(data))
}

func TestNats(t *testing.T) {
	srv1 := getSrv()
	srv2 := getSrv()

	if err := srv1.AddQueueHandler(q1, onRecv1); err != nil {
		fmt.Printf("conn err: %s", err.Error())
	}

	if err := srv2.AddQueueHandler(q1, onRecv2); err != nil {
		fmt.Printf("conn err: %s", err.Error())
	}

	i := 0
	for {
		i++
		srv1.PostQueueMsg(q1, []byte(fmt.Sprintf("%d", i)))
		time.Sleep(100 * time.Millisecond)
	}
}
