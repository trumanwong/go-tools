package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type WechatUser struct {
	OpenId     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        uint32   `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgUrl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionId    string   `json:"unionid"`
}

func GetAccessToken(appId, secret, code string) (*string, *string, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
			appId,
			secret,
			code,
		))
	if err != nil {
		return nil, nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	m := make(map[string]any)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, nil, err
	}
	if _, ok := m["access_token"]; ok {
		accessToken := m["access_token"].(string)
		openid := m["openid"].(string)
		return &accessToken, &openid, nil
	}
	return nil, nil, errors.New("解析失败")
}

func GetUserInfo(accessToken, openid string) (*WechatUser, error) {
	resp, err := http.Get(
		fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
			accessToken,
			openid,
		))
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result WechatUser
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
