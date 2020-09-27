package wechat

import (
	"encoding/json"
	"fmt"
)

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

func (s *Service) GetAccessToken(appID string, appSecret string) (string, error) {

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appID,
		appSecret)

	resp, err := s.http.R().Get(url)

	if err != nil {
		return "", err
	}

	respBody := TokenResp{}
	if err := json.Unmarshal(resp.Body(), &respBody); err != nil {
		return "", err
	}

	if respBody.Errcode != 0 {
		return "", fmt.Errorf("%+v", respBody)
	}

	return respBody.AccessToken, nil
}
