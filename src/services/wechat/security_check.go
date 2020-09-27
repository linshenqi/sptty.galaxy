package wechat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RespSec struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errMsg"`
}

func (s *Service) SecurityCheckContent(token string, content string) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", token)

	body := map[string]string{
		"content": content,
	}

	r := s.http.R().SetBody(&body).SetHeader("content-type", "application/json")

	respBody := RespSec{}
	resp, err := r.Post(url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%+v", resp)
	}

	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		return err
	}

	if respBody.ErrCode != 0 {
		return fmt.Errorf("%+v", respBody.ErrMsg)
	}

	return nil
}
