package gpt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BaiDuAuthorization struct {
	// 访问凭证
	RefreshToken string `json:"refresh_token"`
	// 有效期，Access Token的有效期。
	// 说明：单位是秒，有效期30天
	ExpiresIn int `json:"expires_in"`
	// 错误码
	// 说明：响应失败时返回该字段，成功时不返回
	Error *string `json:"error,omitempty"`
	// 错误描述信息，帮助理解和解决发生的错误
	// 说明：响应失败时返回该字段，成功时不返回
	ErrorDescription *string `json:"error_description,omitempty"`
	// 暂时未使用，可忽略
	SessionKey string `json:"session_key"`
	// 暂时未使用，可忽略
	AccessToken string `json:"access_token"`
	// 暂时未使用，可忽略
	Scope string `json:"scope"`
	// 暂时未使用，可忽略
	SessionSecret string `json:"session_secret"`
}

func GetBaidubceAccessToken(clientId, secret string) (*BaiDuAuthorization, error) {
	url := "https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials"
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf(url, clientId, secret), payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var baiDuAuthorization BaiDuAuthorization
	err = json.Unmarshal(body, &baiDuAuthorization)
	if err != nil {
		return nil, err
	}
	if baiDuAuthorization.Error != nil {
		if baiDuAuthorization.ErrorDescription != nil {
			return nil, fmt.Errorf("[%s]获取百度AI访问凭证失败: %s", *baiDuAuthorization.Error, *baiDuAuthorization.ErrorDescription)
		}
		return nil, fmt.Errorf("[%s]获取百度AI访问凭证失败", *baiDuAuthorization.Error)
	}
	return &baiDuAuthorization, nil
}
