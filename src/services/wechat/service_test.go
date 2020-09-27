package wechat

import (
	"fmt"
	"testing"

	"github.com/linshenqi/sptty"
)

func getSrv() *Service {
	return &Service{
		http: sptty.CreateHttpClient(sptty.DefaultHttpClientConfig()),
	}
}

func TestContentCheck(t *testing.T) {

	srv := getSrv()
	token, err := srv.GetAccessToken("wx77d83a2aa6c324ab", "2c6ecb6fe8a0394715704149a6afc56b")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	content := fmt.Sprintf("%+v", map[string]string{
		"k1": "裸聊",
		"we": "冰毒,迷奸",
	})

	content = "裸聊"

	if err := srv.SecurityCheckContent(token, content); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("通过检测")
}
