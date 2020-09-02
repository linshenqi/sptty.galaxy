package aliyun_mns

import (
	"fmt"
	"testing"
	"time"

	alimns "github.com/aliyun/aliyun-mns-go-sdk"
)

const (
	q1 = "ashibro-message"
)

var srv1 *Service
var srv2 *Service

func getSrv() *Service {
	return &Service{cfg: Config{
		Url:          "http://1693478565140903.mns.cn-hangzhou.aliyuncs.com/",
		AccessKey:    "",
		AccessSecret: "",
		Queues: []string{
			q1,
		},
	}}
}

func onQueueMsg1(queue string, msg *alimns.MessageReceiveResponse, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if msg == nil {
		return
	}

	if err := srv1.DeleteQueueMsg(queue, msg.ReceiptHandle); err != nil {
		fmt.Printf("1 delete failed: %s\n", err.Error())
	}

	fmt.Printf("1 recv: %s\n", msg.MessageBody)
}

func onQueueMsg2(queue string, msg *alimns.MessageReceiveResponse, err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if msg == nil {
		return
	}

	if err := srv1.DeleteQueueMsg(queue, msg.ReceiptHandle); err != nil {
		fmt.Printf("2 delete failed: %s\n", err.Error())
	}

	fmt.Printf("2 recv: %s\n", msg.MessageBody)
}

func TestNMS(t *testing.T) {
	// content := "{\"author_id\":\"broau47bd4mu467s8tgg\",\"content\":\"用户申请取消行程 用户ID: broau47bd4mu467s8tgg 行程ID: bt2tfuvbd4mq3ml1hhjg\",\"session_id\":\"internal_customer_service\",\"target_user_id\":\"internal_customer_service\",\"type\":\"text\"}"
	srv1 = getSrv()
	srv2 = getSrv()

	if err := srv1.doInit(); err != nil {
		fmt.Println(err.Error())
	}

	if err := srv2.doInit(); err != nil {
		fmt.Println(err.Error())
	}

	if err := srv1.AddQueueHandler(q1, onQueueMsg1); err != nil {
		fmt.Println(err.Error())
	}

	if err := srv2.AddQueueHandler(q1, onQueueMsg2); err != nil {
		fmt.Println(err.Error())
	}

	i := 0
	for {
		i++
		// if err := srv1.PostQueueMsg(q1, fmt.Sprintf("%d: %s", i, content)); err != nil {
		// 	fmt.Println(err.Error())
		// }

		// i++
		// if err := srv1.PostQueueMsg(q1, fmt.Sprintf("%d: %s", i, content)); err != nil {
		// 	fmt.Println(err.Error())
		// }

		time.Sleep(1 * time.Second)
	}
}
